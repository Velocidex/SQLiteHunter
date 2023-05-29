package definitions

import (
	"regexp"
	"strings"

	"github.com/Velocidex/SQLiteHunter/api"
)

var (
	templateRegex = regexp.MustCompile("\\{\\{.+\\}\\}")
)

func ExpandGlobs(definition api.Definition, config_obj *api.ConfigDefinitions) []string {
	var res []string

	for _, glob := range definition.Globs {
		has_template := false

		// Expand each template into the glob
		templateRegex.ReplaceAllStringFunc(glob,
			func(in string) string {
				has_template = true
				subs, pres := config_obj.Globs[strings.Trim(in, "{}")]
				if !pres {
					res = append(res, glob)
					return ""
				}
				// Add a glob for each substitution
				for _, s := range subs {
					res = append(res,
						strings.ReplaceAll(glob, in, s))
				}
				return ""
			})

		// The glob has no template expansion in it just include it
		// literally.
		if !has_template {
			res = append(res, glob)
		}
	}

	return res
}
