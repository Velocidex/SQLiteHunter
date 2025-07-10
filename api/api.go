package api

type ConfigDefinitions struct {
	Globs map[string][]string `yaml:"Globs"`
}

type Definition struct {
	Name                string      `yaml:"Name" json:"Name,omitempty"`
	Author              string      `yaml:"Author" json:"Author,omitempty"`
	Description         string      `yaml:"Description" json:"Description,omitempty"`
	Email               string      `yaml:"Email" json:"Email,omitempty"`
	Reference           string      `yaml:"Reference" json:"Reference,omitempty"`
	Categories          []string    `yaml:"Categories" json:"Categories,omitempty"`
	SQLiteIdentifyQuery string      `yaml:"SQLiteIdentifyQuery" json:"SQLiteIdentifyQuery,omitempty"`
	SQLiteIdentifyValue interface{} `yaml:"SQLiteIdentifyValue" json:"SQLiteIdentifyValue,omitempty"`
	Globs               []string    `yaml:"Globs" json:"Globs,omitempty"`
	FilenameRegex       string      `yaml:"FilenameRegex" json:"FilenameRegex,omitempty"`
	Sources             []Source    `yaml:"Sources" json:"Sources,omitempty"`

	Filename_ string `yaml:"Filename" json:"Filename,omitempty"`
	RawData_  string `yaml:"RawData" json:"RawData,omitempty"`
}

type Source struct {
	Name string `yaml:"name"`

	// VQL to include prior to the VQL query - for example contains
	// custom VQL functions
	Preamble string `yaml:"Preamble" json:"Preamble,omitempty"`

	// Specialized VQL to post process the rows. Default is a
	// passthrough `SELECT * FROM Rows`
	VQL                 string      `yaml:"VQL" json:"VQL,omitempty"`
	SQL                 string      `yaml:"SQL" json:"SQL,omitempty"`
	SQLiteIdentifyQuery string      `json:"id_query,omitempty"`
	SQLiteIdentifyValue interface{} `json:"id_value,omitempty"`
	Filename            string      `json:"filename,omitempty"`
}
