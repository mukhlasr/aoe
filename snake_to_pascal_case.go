package main

import "strings"

func snakeToPascalCase(str string) string {
	// make sure to lowercase it all
	str = strings.ToLower(str)

	builder := strings.Builder{}
	_ = builder.WriteByte(capitalize(str[0]))

	for i := 1; i < len(str); i++ {
		if str[i] == '_' && i < len(str)-1 {
			_ = builder.WriteByte(capitalize(str[i+1]))
			i++
			continue
		}

		_ = builder.WriteByte(str[i])
	}

	return builder.String()
}

func capitalize(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}

	return b
}
