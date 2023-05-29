package api

type ConfigDefinitions struct {
	Globs map[string][]string `yaml:"Globs"`
}

type Definition struct {
	Name                string      `yaml:"Name"`
	Author              string      `yaml:"Author"`
	Email               string      `yaml:"Email"`
	Category            string      `yaml:"Category"`
	SQLiteIdentifyQuery string      `yaml:"SQLiteIdentifyQuery"`
	SQLiteIdentifyValue interface{} `yaml:"SQLiteIdentifyValue"`
	Globs               []string    `yaml:"Globs"`
	FilenameRegex       string      `yaml:"FilenameRegex"`
	Sources             []Source    `yaml:"Sources"`
}

type Source struct {
	Name string `yaml:"name"`
	// Specialized VQL to post process the rows. Default is a
	// passthrough `SELECT * FROM Rows`
	VQL                 string      `yaml:"VQL"`
	SQL                 string      `yaml:"SQL"`
	SQLiteIdentifyQuery string      `json:"id_query"`
	SQLiteIdentifyValue interface{} `json:"id_value"`
	Filename            string      `json:"filename"`
}
