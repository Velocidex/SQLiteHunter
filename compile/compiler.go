package compile

import (
	"fmt"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/SQLiteHunter/definitions"
	"github.com/Velocidex/ordereddict"
)

type Artifact struct {
	Spec              api.Spec
	Name, Description string
	Category          *ordereddict.Dict
	Sources           []api.Definition
}

func newArtifact() *Artifact {
	return &Artifact{
		Name:        "Generic.Forensic.SQLiteHunter",
		Description: "Hunt for SQLite files",
		Category:    ordereddict.NewDict(),
		Spec: api.Spec{
			Sources: ordereddict.NewDict(),
		},
	}
}

// Build the spec and artifact
func Compile(defs []api.Definition,
	config_obj *api.ConfigDefinitions) (*Artifact, error) {

	res := newArtifact()
	for _, d := range defs {
		category := d.Category
		if category == "" {
			category = "Misc"
		}
		res.Category.Set(category, true)
		globs := definitions.ExpandGlobs(d, config_obj)
		for _, g := range globs {
			res.Spec.Globs = append(res.Spec.Globs, api.GlobSpec{
				Glob:     g,
				Tag:      category,
				Filename: d.FilenameRegex,
			})
		}

		// Each definition can contain multiple queries. Each such
		// query ends up in a separate source.
		for idx, s := range d.Sources {
			// Calculate a unique name for the source
			name := s.Name
			if name == "" {
				name = fmt.Sprintf("%v", idx)
			}

			if name == "0" {
				s.Name = d.Name
			} else {
				s.Name = d.Name + "_" + name
			}

			if s.SQLiteIdentifyQuery == "" {
				s.SQLiteIdentifyQuery = d.SQLiteIdentifyQuery
			}

			if s.SQLiteIdentifyValue == nil {
				s.SQLiteIdentifyValue = d.SQLiteIdentifyValue
			}

			s.Filename = d.FilenameRegex
			res.Spec.Sources.Update(s.Name, s)
		}
	}

	return res, nil
}
