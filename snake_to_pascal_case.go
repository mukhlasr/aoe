package main

import (
	"strings"
)

func snakeToPascalCase(str string) string {
	// make sure to lowercase it all
	var res string

	checkpoint := 0
	for i := 0; i <= len(str); i++ {
		if i == len(str) || str[i] == '_' {
			res += CapitalizeWord(str[checkpoint:i])
			checkpoint = i + 1
			continue
		}
	}

	return res
}

func CapitalizeWord(str string) string {
	str = strings.ToLower(str)
	if str == "" {
		return ""
	}

	if str == "id" {
		return "ID"
	}

	b := []byte(str)
	b[0] = capitalize(b[0])

	return string(b)
}

func capitalize(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b - 'a' + 'A'
	}

	return b
}
