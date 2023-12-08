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
		Name: "Generic.Forensic.SQLiteHunter",
		Description: `Hunt for SQLite files.

SQLite has become the de-facto standard for storing application data,
in many types of applications:

- Web Browsers
- Operating Systems
- Various applications, such as iMessage, TCC etc

This artifact can hunt for these artifacts in a mostly automated way.
More info at https://github.com/Velocidex/SQLiteHunter

NOTE: If you want to use this artifact on just a bunch of files already
collected (for example the files collected using the
Windows.KapeFiles.Targets artifact) you can use the CustomGlob parameter
(for example set it to "/tmp/unpacked/**" to consider all files in the
unpacked directory).

`,
		Category: ordereddict.NewDict().Set("All", true),
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
		categories := d.Categories
		if len(categories) == 0 {
			categories = []string{"Misc"}
		}

		// All artifacts include the All category as well.
		categories = append(categories, "All")

		for _, c := range categories {
			res.Category.Update(c, true)
		}

		globs := definitions.ExpandGlobs(d, config_obj)
		for _, g := range globs {
			res.Spec.Globs = append(res.Spec.Globs, api.GlobSpec{
				Glob:     g,
				Tags:     categories,
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
