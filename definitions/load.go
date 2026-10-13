package definitions

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Velocidex/SQLiteHunter/api"
	"gopkg.in/yaml.v3"
)

// Artifact names are more strict now so we need to clean them up.
var artifactNameRegex = regexp.MustCompile("[^a-zA-Z0-9_.]")

func SanitizeName(name string) string {
	return artifactNameRegex.ReplaceAllLiteralString(name, "_")
}

func LoadDefinitions(root string) ([]api.Definition, error) {
	var definitions []api.Definition

	err := filepath.Walk(root,
		func(path string, d fs.FileInfo, err error) error {
			if err == nil && strings.HasSuffix(path, ".yaml") {
				var definition api.Definition
				fd, err := os.Open(path)
				if err != nil {
					fmt.Printf("Error processing %v: %v\n", path, err)
					return nil
				}

				data, err := io.ReadAll(fd)
				if err != nil {
					fmt.Printf("Error processing %v: %v\n", path, err)
					return nil
				}

				err = yaml.Unmarshal(data, &definition)
				if err != nil {
					fmt.Printf("Error processing %v: %v\n", path, err)
					return err
				}
				definition.Filename_ = path
				definition.RawData_ = string(data)

				definitions = append(definitions, sanitizeDefinition(definition))
			}
			return nil
		})

	return definitions, err
}

func sanitizeDefinition(definition api.Definition) api.Definition {
	definition.Name = SanitizeName(definition.Name)
	for _, s := range definition.Sources {
		s.Name = SanitizeName(s.Name)
	}

	return definition
}
