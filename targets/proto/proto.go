// Package proto contains proto related mage targets.
package proto

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const (
	skimCMD = "skim"
)

var (
	Owner             = "taxibeat"
	Registry          = "proto-schemas"
	SchemasLocation   = "proto/schemas"
	GeneratedLocation = "proto/generated"
	Service           = "" // Should be set in magefile init func.
)

// Proto groups together proto related tasks.
type Proto mg.Namespace

// SchemaGenerate generates a single proto schema.
func (Proto) SchemaGenerate(schema, version string) error {
	if schema == "" {
		return errors.New("schema is mandatory")
	}
	if version == "" {
		return errors.New("version is mandatory")
	}

	pathToSchema := fmt.Sprintf("%s/%s/%s.proto", schema, version, schema)
	fmt.Printf("proto schema: generate schema %s\n", pathToSchema)

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %s", err)
	}

	args := append(
		getDefaultSkimArgs(Service),
		"generate",
		"-s",
		SchemasLocation,
		"-out",
		tmpDir,
		"--schema",
		pathToSchema,
	)

	err = sh.RunV(skimCMD, args...)
	if err != nil {
		return err
	}

	generatedFiles, err := getGeneratedFiles(tmpDir)
	if err != nil {
		return err
	}
	return moveGeneratedFiles(generatedFiles)
}

// SchemaGenerateAll generates all the schemas found.
func (Proto) SchemaGenerateAll() error {
	fmt.Printf("proto schema: generate all schemas for service: %q\n", Service)

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("failed to create tmp dir: %s", err)
	}

	args := append(
		getDefaultSkimArgs(Service),
		"generate-all",
		"-s",
		SchemasLocation,
		"-out",
		tmpDir,
	)

	err = sh.RunV(skimCMD, args...)
	if err != nil {
		return err
	}

	generatedFiles, err := getGeneratedFiles(tmpDir)
	if err != nil {
		return err
	}
	return moveGeneratedFiles(generatedFiles)
}

// SchemaValidateAll lints the schemas in the repository against the GitHub schemas.
func (p Proto) SchemaValidateAll() error {
	fmt.Printf("proto schema: validate all schemas for service: %q\n", Service)

	args := append(
		getDefaultSkimArgs(Service),
		"-t",
		os.Getenv("GITHUB_TOKEN"),
		"validate-all",
		"-s",
		SchemasLocation,
	)

	return sh.RunV(skimCMD, args...)
}

func getDefaultSkimArgs(service string) []string {
	return []string{
		"-r",
		Registry,
		"-o",
		Owner,
		"-n",
		service,
	}
}

func getGeneratedFiles(tmpDir string) ([]string, error) {
	var generatedFiles []string
	err := filepath.Walk(tmpDir, func(path string, fInfo os.FileInfo, err error) error {
		if fInfo.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		generatedFiles = append(generatedFiles, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list generated files: %v", err)
	}
	return generatedFiles, nil
}

func moveGeneratedFiles(generatedFiles []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %v", wd)
	}

	for _, generatedFile := range generatedFiles {
		fileName := filepath.Base(generatedFile)
		schemaName := strings.Split(fileName, ".")[0]
		outDir := fmt.Sprintf("%s/%s/%s", wd, GeneratedLocation, schemaName)
		err := os.MkdirAll(outDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create out dir: %v", err)
		}

		outFile := fmt.Sprintf("%s/%s", outDir, fileName)
		err = os.Rename(generatedFile, outFile)
		if err != nil {
			return fmt.Errorf("failed to move generated file: %v", err)
		}
		fmt.Printf("schema generated successfully under %s\n", outFile)
	}
	return nil
}
