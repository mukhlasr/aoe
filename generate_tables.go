package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"

	"golang.org/x/tools/imports"
)

type TemplateTableFieldModel struct {
	Name string
	Type string
}

type TemplateTableModel struct {
	Name   string
	Fields []TemplateTableFieldModel
}

func generateTablesAsModels(packageName string, out io.Writer) error {
	tmpl, err := template.ParseFS(templateFS, "templates/model.go.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	model := struct {
		PackageName      string
		ImportedPackages []string
		Types            []TemplateTableModel
	}{
		PackageName: packageName,
	}

	tables, err := GetTables(context.TODO())
	if err != nil {
		return err
	}

	mapPackages := map[string]struct{}{}
	for _, table := range tables {
		var tableModel TemplateTableModel

		tableModel.Name = table.GoTypeName()
		for _, column := range table.Columns {
			var fieldModel TemplateTableFieldModel

			fieldModel.Name = snakeToPascalCase(column.Name)
			fieldModel.Type = column.GoType().Name

			if column.GoType().PackageName != "" {
				mapPackages[column.GoType().PackageName] = struct{}{}
			}

			tableModel.Fields = append(tableModel.Fields, fieldModel)
		}

		model.Types = append(model.Types, tableModel)
	}

	for pack := range mapPackages {
		model.ImportedPackages = append(model.ImportedPackages, pack)
	}

	err = tmpl.Execute(&buf, model)
	if err != nil {
		return err
	}

	b, err := imports.Process("", buf.Bytes(), &imports.Options{
		FormatOnly: true,
		Comments:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to format source code: %w", err)
	}

	_, err = io.Copy(out, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return err
}
