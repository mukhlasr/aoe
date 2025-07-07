package main

import "testing"

func TestSnakeToPascalCase(t *testing.T) {
	for _, testCase := range []struct {
		input  string
		output string
	}{
		{"foo_bar_baz", "FooBarBaz"},
		{"FOO_BAR_BAZ", "FooBarBaz"},
		{"foo_bar_123", "FooBar123"},
		{"foobar123", "Foobar123"},
		{"foo_bar_b", "FooBarB"},
		{"foo_bar_", "FooBar"},
		{"_bar", "Bar"},
		{"user_id", "UserID"},
		{"id", "ID"},
	} {
		output := snakeToPascalCase(testCase.input)
		if output != testCase.output {
			t.Errorf("expecting %s, but got: %s", testCase.output, output)
		}
	}
}
