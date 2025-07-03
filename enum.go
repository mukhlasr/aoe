package main

import "context"

type Enum struct {
	Name   string
	Values []string
}

const qGetEnum = `
SELECT
	t.typname AS enum_name,
	e.enumlabel AS enum_value
FROM
	pg_type t
	JOIN pg_enum e ON t.oid = e.enumtypid
	JOIN pg_namespace n ON n.oid = t.typnamespace
WHERE
	n.nspname = 'public'
ORDER BY
	t.typname, e.enumsortorder`

func GetEnums(ctx context.Context) ([]Enum, error) {
	rows, err := DB.Query(ctx, qGetEnum)
	if err != nil {
		return nil, err
	}

	m := map[string][]string{}
	for rows.Next() {
		var name, val string
		if err := rows.Scan(&name, &val); err != nil {
			break
		}

		vals := m[name]

		m[name] = append(vals, val)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var res []Enum
	for name, vals := range m {
		res = append(res, Enum{
			Name:   name,
			Values: vals,
		})
	}

	return res, nil

}

func (e Enum) GoTypeName() string {
	return snakeToPascalCase(e.Name)
}

// GoConstants returns map[constant_name] = value
func (e Enum) GoConstants() map[string]string {
	if len(e.Values) == 0 {
		return nil
	}

	res := map[string]string{}

	for _, v := range e.Values {
		key := snakeToPascalCase(e.Name + "_" + v)
		res[key] = v
	}

	return res
}
