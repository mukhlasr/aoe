package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"

	"golang.org/x/tools/imports"
)

type TemplateConstModel struct {
	Name  string
	Type  string
	Value string
}

type TemplateEnumModel struct {
	Name   string
	Consts []TemplateConstModel
}

func GenerateEnumsAsConstants(packageName string, out io.Writer) error {
	tmpl, err := template.ParseFS(templateFS, "templates/const.go.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	model := struct {
		PackageName string
		Types       []TemplateEnumModel
	}{
		PackageName: packageName,
	}

	enums, err := GetEnums(context.TODO())
	if err != nil {
		return err
	}

	for _, enum := range enums {
		var enumModel TemplateEnumModel

		typeName := enum.GoTypeName()

		enumModel.Name = typeName
		for _, val := range enum.Values {
			var constModel TemplateConstModel

			constModel.Name = enum.GoTypeName() + snakeToPascalCase(val)
			constModel.Type = typeName
			constModel.Value = val

			enumModel.Consts = append(enumModel.Consts, constModel)
		}

		model.Types = append(model.Types, enumModel)
	}

	err = tmpl.Execute(&buf, model)
	if err != nil {
		return err
	}

	b, err := imports.Process("", buf.Bytes(), &imports.Options{
		FormatOnly: true,
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
