package api

import "github.com/Velocidex/ordereddict"

// This is the main data structure that will be used by the
// application.
type Spec struct {
	Globs []GlobSpec `json:"globs"`

	// map[string]Source
	Sources *ordereddict.Dict `json:"sources"`
}

type GlobSpec struct {
	Glob     string `json:"glob"`
	Tag      string `json:"tag"`
	Filename string `json:"name"`
}
