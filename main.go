package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/SQLiteHunter/compile"
	"github.com/Velocidex/SQLiteHunter/definitions"
	"gopkg.in/yaml.v3"
)

const (
	CONFIG = "./config.yaml"
)

func loadConfig() (*api.ConfigDefinitions, error) {
	fd, err := os.Open(CONFIG)
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
	config_obj, err := loadConfig()
	if err != nil {
		panic(err)
	}

	defs, err := definitions.LoadDefinitions("./definitions")
	if err != nil {
		panic(err)
	}

	spec, err := compile.Compile(defs, config_obj)
	if err != nil {
		panic(err)
	}

	fmt.Println(spec.Yaml())
}
