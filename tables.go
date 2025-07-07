package main

import (
	"context"

	"github.com/jinzhu/inflection"
)

type Table struct {
	Name    string
	Columns []Column
}

const qGetTables = `
SELECT table_name
	FROM information_schema.tables
	WHERE table_schema = 'public' AND table_type = 'BASE TABLE'`

func GetTables(ctx context.Context) ([]Table, error) {
	var tabNames []string

	rows, err := DB.Query(ctx, qGetTables)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tabName string
		if err := rows.Scan(&tabName); err != nil {
			return nil, err
		}

		tabNames = append(tabNames, tabName)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var res []Table

	for _, n := range tabNames {
		var tab Table

		if n == "migrations" {
			continue
		}

		tab.Name = n
		cols, err := GetTableColumns(ctx, n)
		if err != nil {
			return nil, err
		}

		tab.Columns = cols

		res = append(res, tab)
	}

	return res, nil
}

func (t Table) GoTypeName() string {
	name := snakeToPascalCase(t.Name)
	return inflection.Singular(name)
}
