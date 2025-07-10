package compile

import (
	"fmt"
	"regexp"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/SQLiteHunter/definitions"
	"github.com/Velocidex/ordereddict"
)

var (
	selectRegex = regexp.MustCompile(`(?i)^\s*SELECT`)
)

type Artifact struct {
	Spec     api.Spec
	Name     string
	Category *ordereddict.Dict
	Sources  []api.Definition
}

func newArtifact() *Artifact {
	return &Artifact{
		Name:     "Generic.Forensic.SQLiteHunter",
		Category: ordereddict.NewDict(),
		Spec: api.Spec{
			Sources: ordereddict.NewDict(),
		},
	}
}

// Build the spec and artifact
func Compile(defs []api.Definition,
	config_obj *api.ConfigDefinitions) (*Artifact, error) {

	res := newArtifact()
	res.Sources = defs

	for _, d := range defs {
		categories := d.Categories
		if len(categories) == 0 {
			categories = []string{"Misc"}
		}

		for _, c := range categories {
			res.Category.Update(c, true)
		}

		globs := definitions.ExpandGlobs(d, config_obj)
		for _, g := range globs {
			res.Spec.Globs = append(res.Spec.Globs, api.GlobSpec{
				Glob:     g,
				Tags:     categories,
				Filename: d.FilenameRegex,
				Rule:     d.Name,
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

			if s.VQL != "" && !selectRegex.MatchString(s.VQL) {
				return nil, fmt.Errorf("Source %v: %v/%v: Invalid VQL clause: Must begin with SELECT. Any VQL definitions should go in a Preamble",
					d.Filename_, d.Name, name)
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
