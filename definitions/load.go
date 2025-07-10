package definitions

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Velocidex/SQLiteHunter/api"
	"gopkg.in/yaml.v3"
)

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

				data, err := ioutil.ReadAll(fd)
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

				definitions = append(definitions, definition)
			}
			return nil
		})

	return definitions, err
}
