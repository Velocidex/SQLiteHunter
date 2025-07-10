package compile

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/ordereddict"

	_ "embed"
)

//go:embed template.yaml
var artifact_template string

type Expansion struct {
	ArtifactName   string
	CompressedSpec string
	Categories     []string
	Spec           api.Spec
}

type dictValue struct {
	Key   string
	Value interface{}
}

// Produce the YAML of the artifact definition.
func (self *Artifact) Yaml() (string, error) {
	tmpl, err := template.New("").Funcs(
		template.FuncMap{
			"Indent": func(in string) interface{} {
				return strings.TrimSpace(indent(in, 4))
			},
			"Quote": func(in interface{}) interface{} {
				return fmt.Sprintf("%q", in)
			},
			"DictRange": func(in interface{}) interface{} {
				var res []dictValue

				dict, ok := in.(*ordereddict.Dict)
				if ok {

					for _, k := range dict.Keys() {
						v, _ := dict.Get(k)
						res = append(res, dictValue{Key: k, Value: v})
					}
				}
				return res
			},
		}).Parse(artifact_template)
	if err != nil {
		return "", err
	}

	exp := &Expansion{
		ArtifactName:   self.Name,
		CompressedSpec: self.encodeSpec(),
		Categories:     self.Category.Keys(),
		Spec:           self.Spec,
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, exp)
	return string(b.Bytes()), err
}

func (self *Artifact) encodeSpec() string {
	serialized, _ := json.Marshal(self.Spec)

	// Compress the string
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(serialized)
	gz.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func (self *Artifact) BuildIndex() []byte {
	var defs []api.Definition

	for _, d := range self.Sources {
		var sources []api.Source
		for _, s := range d.Sources {
			sources = append(sources, api.Source{
				Name: s.Name,
			})
		}

		defs = append(defs, api.Definition{
			RawData_:    d.RawData_,
			Name:        d.Name,
			Author:      d.Author,
			Description: d.Description,
			Categories:  d.Categories,
			Sources:     sources,
		})
	}

	serialized, _ := json.Marshal(defs)
	return serialized
}
