package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/SQLiteHunter/compile"
	"github.com/Velocidex/SQLiteHunter/definitions"
	"gopkg.in/yaml.v3"
)

func loadConfig(config_path string) (*api.ConfigDefinitions, error) {
	fd, err := os.Open(config_path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	config_obj := &api.ConfigDefinitions{}
	err = yaml.Unmarshal(data, config_obj)

	return config_obj, err
}

func main() {
	config_path := flag.String("config", "./config.yaml",
		"The path to the config file")

	definition_directory := flag.String("definition_directory", "./definitions",
		"A directory containing all definitiions")

	flag.Parse()

	config_obj, err := loadConfig(*config_path)
	if err != nil {
		panic(err)
	}

	defs, err := definitions.LoadDefinitions(*definition_directory)
	if err != nil {
		panic(err)
	}

	spec, err := compile.Compile(defs, config_obj)
	if err != nil {
		panic(err)
	}

	// Serialize the artifact to YAML
	fmt.Println(spec.Yaml())
}
