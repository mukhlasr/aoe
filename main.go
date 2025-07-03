package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
)

func main() {
	var (
		dsn    string
		schema string
	)
	flag.StringVar(&dsn, "dsn", "", "dsn to postgres")
	flag.StringVar(&schema, "schema", "public", "schema name")
	flag.Parse()

	if dsn == "" {
		flag.PrintDefaults()
		return
	}

	if err := initDB(context.TODO(), dsn); err != nil {
		fmt.Println("aoe: failed to init db:", err)
		os.Exit(-1)
	}

	log.Println(generateEnum())
}

func generateEnum() error {
	type tmplModel struct {
		PackageName string
		Types       []Enum
	}
	tmpl, err := template.ParseFS(templateFS, "templates/enum.go.tmpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	enums, err := GetEnums(context.TODO())
	if err != nil {
		return err
	}

	err = tmpl.Execute(&buf, tmplModel{PackageName: "db", Types: enums})
	if err != nil {
		return err
	}

	b, err := io.ReadAll(&buf)
	if err != nil {
		return err
	}

	b, err = format.Source(b)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
