package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func main() {
	var (
		dsn         string
		schema      string
		dir         string
		packageName string
	)
	flag.StringVar(&dsn, "dsn", "", "dsn to postgres")
	flag.StringVar(&schema, "schema", "public", "schema name")
	flag.StringVar(&dir, "dir", "./", "output directory")
	flag.StringVar(&packageName, "package", "", "package name for the generated code. The default will be the same as output directory")
	flag.Parse()

	if dsn == "" || dir == "" {
		flag.PrintDefaults()
		return
	}

	exitOnError(initDB(context.TODO(), dsn))

	if packageName == "" {
		name, err := packageNameFromDir(dir)
		exitOnError(err)

		packageName = name
	}

	constFile, err := os.OpenFile(path.Join(dir, "const.go"), os.O_CREATE|os.O_WRONLY, 0644)
	exitOnError(err)

	modelFile, err := os.OpenFile(path.Join(dir, "model.go"), os.O_CREATE|os.O_WRONLY, 0644)
	exitOnError(err)

	exitOnError(GenerateEnumsAsConstants(packageName, constFile))
	exitOnError(generateTablesAsModels(packageName, modelFile))
}

func packageNameFromDir(dir string) (string, error) {
	d, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to get package name from directory: %w", err)
	}

	if d == "/" {
		return "", errors.New("unable to generate on the root directory")
	}

	return path.Base(d), nil

}

func exitOnError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, "aoe:", err.Error())
}
