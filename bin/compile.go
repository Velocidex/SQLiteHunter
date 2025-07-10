package main

import (
	"os"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/SQLiteHunter/compile"
	"github.com/Velocidex/SQLiteHunter/definitions"
	"github.com/alecthomas/kingpin"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed config.yaml
var config_string string

var (
	compile_cmd = app.Command("compile",
		"Build an artifact from a set of SQLiteHunter yaml files.")

	compile_yaml = compile_cmd.Arg("input",
		"Path to the SQLiteHunter yaml directory to compile").
		Required().String()

	output_artifact = compile_cmd.Arg("output",
		"Where to write the final artifact").
		Required().String()

	output_zip = compile_cmd.Flag("output_zip",
		"Produce a ZIP file we can use to hunt").
		String()

	output_index = compile_cmd.Flag("index", "Where to write the rules index").
			String()
)

func loadConfig() (*api.ConfigDefinitions, error) {
	config_obj := &api.ConfigDefinitions{}
	err := yaml.Unmarshal([]byte(config_string), config_obj)
	return config_obj, err
}

func doCompile() error {
	config_obj, err := loadConfig()
	if err != nil {
		return err
	}

	defs, err := definitions.LoadDefinitions(*compile_yaml)
	if err != nil {
		return err
	}

	spec, err := compile.Compile(defs, config_obj)
	if err != nil {
		return err
	}

	// Serialize the artifact to YAML
	res, err := spec.Yaml()
	if err != nil {
		return err
	}

	fd, err := os.OpenFile(*output_artifact,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write([]byte(res))
	if err != nil {
		return err
	}

	if *output_zip != "" {
		err := spec.WriteZip(*output_zip)
		if err != nil {
			return err
		}
	}

	if *output_index != "" {
		fd, err := os.OpenFile(*output_index,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		defer fd.Close()

		_, err = fd.Write(spec.BuildIndex())
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case compile_cmd.FullCommand():
			err := doCompile()
			kingpin.FatalIfError(err, "Compiling artifact")

		default:
			return false
		}
		return true
	})
}
