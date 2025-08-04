package main

import "context"

const defaultUUIDPackage = "github.com/google/uuid"

type GoType struct {
	Name        string
	PackageName string
}

type Column struct {
	Name     string
	DataType string
	UDTName  string
	Nullable bool
}

const qGetColumn = `
SELECT
	c.column_name,
	c.udt_name,
	c.is_nullable,
	c.data_type
FROM
	information_schema.columns c
WHERE
	c.table_schema = 'public'
	AND c.table_name = $1
ORDER BY
	c.ordinal_position`

func GetTableColumns(ctx context.Context, tabName string) ([]Column, error) {
	var res []Column

	rows, err := DB.Query(ctx, qGetColumn, tabName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var colName, udtName, isNullable, dataType string
		if err := rows.Scan(&colName, &udtName, &isNullable, &dataType); err != nil {
			return nil, err
		}

		res = append(res, Column{
			Name:     colName,
			UDTName:  udtName,
			Nullable: isNullable == "YES",
			DataType: dataType,
		})
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (c Column) GoType() GoType {
	var typ string
	var packageName string
	switch c.DataType {
	case "integer":
		typ = "int"
		if c.Nullable {
			packageName = "database/sql"
			typ = "sql.NullInt32"
		}
	case "bigint":
		typ = "int64"
		if c.Nullable {
			packageName = "database/sql"
			typ = "sql.NullInt64"
		}
	case "text", "character varying":
		typ = "string"

		if c.Nullable {
			packageName = "database/sql"
			typ = "sql.NullString"
		}
	case "boolean":
		typ = "bool"

		if c.Nullable {
			packageName = "database/sql"
			typ = "sql.NullBool"
		}
	case "timestamp without time zone", "timestamp with time zone":
		packageName = "time"
		typ = "time.Time"

		if c.Nullable {
			packageName = "database/sql"
			typ = "sql.NullTime"
		}
	case "uuid":
		packageName = defaultUUIDPackage
		typ = "uuid.UUID"

		if c.Nullable {
			typ = "uuid.NullUUID"
		}
	case "USER-DEFINED":
		typ = snakeToPascalCase(c.UDTName)

		if c.Nullable {
			typ = "Null" + typ
		}
	case "numeric":
		typ = "float64"

		if c.Nullable {
			packageName = "database/sql"
			typ = "sql.NullFloat64"
		}
	default:
		// Fallback to any for unsupported types
		typ = "any"
	}

	return GoType{
		Name:        typ,
		PackageName: packageName,
	}
}
