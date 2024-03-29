package compile

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Velocidex/SQLiteHunter/api"
	"github.com/Velocidex/SQLiteHunter/utils"
)

// Produce the YAML of the artifact definition.
func (self *Artifact) Yaml() string {
	return fmt.Sprintf(`
name: %v
description: |
%v

column_types:
- name: Image
  type: preview_upload

export: |
  LET SPEC <= %q
  LET Specs <= parse_json(data=gunzip(string=base64decode(string=SPEC)))
  LET CheckHeader(OSPath) = read_file(filename=OSPath, length=12) = "SQLite forma"
  LET Bool(Value) = if(condition=Value, then="Yes", else="No")

  -- In fast mode we check the filename, then the header then run the sqlite precondition
  LET matchFilename(SourceName, OSPath) = OSPath =~ get(item=Specs.sources, field=SourceName).filename
    AND CheckHeader(OSPath=OSPath)
    AND Identify(SourceName= SourceName, OSPath= OSPath)
    AND log(message=format(format="%%v matched by filename %%v",
            args=[OSPath, get(item=Specs.sources, field=SourceName).filename]))

  -- If the user wanted to also upload the file, do so now
  LET MaybeUpload(OSPath) = if(condition=AlsoUpload, then=upload(file=OSPath)) OR TRUE

  LET Identify(SourceName, OSPath) = SELECT if(
    condition=CheckHeader(OSPath=OSPath),
    then={
      SELECT *
      FROM sqlite(file=OSPath, query=get(item=Specs.sources, field=SourceName).id_query)
    }) AS Hits
  FROM scope()
  WHERE if(condition=Hits[0].Check = get(item=Specs.sources, field=SourceName).id_value,
    then= log(message="%%v was identified as %%v",
            args=[OSPath, get(item=Specs.sources, field=SourceName).Name]),
    else=log(message="%%v was not identified as %%v (got %%v, wanted %%v)",
             args=[OSPath, get(item=Specs.sources, field=SourceName).Name, str(str=Hits),
                   get(item=Specs.sources, field=SourceName).id_value]) AND FALSE)

  LET ApplyFile(SourceName) = SELECT * FROM foreach(row={
     SELECT OSPath FROM AllFiles
     WHERE if(condition=MatchFilename,  then=matchFilename(SourceName=SourceName, OSPath=OSPath),
      else=Identify(SourceName= SourceName, OSPath= OSPath))

  }, query={
     SELECT *, OSPath FROM sqlite(
        file=OSPath, query=get(item=Specs.sources, field=SourceName).SQL)
  })

  -- Filter for matching files without sqlite checks.
  LET FilterFile(SourceName) =
     SELECT OSPath FROM AllFiles
     WHERE if(condition=MatchFilename,
              then=OSPath =~ get(item=Specs.sources, field=SourceName).filename)

  -- Build a regex for all enabled categories.
  LET all_categories = SELECT _value FROM foreach(row=%v) WHERE get(field=_value)
  LET category_regex <= join(sep="|", array=all_categories._value)
  LET AllGlobs <= filter(list=Specs.globs, condition="x=> x.tags =~ category_regex")
  LET _ <= log(message="Globs for category %%v is %%v", args=[category_regex, CustomGlob || AllGlobs.glob])
  LET AllFiles <= SELECT OSPath FROM glob(globs=CustomGlob || AllGlobs.glob)
    WHERE NOT IsDir AND MaybeUpload(OSPath=OSPath)

parameters:
- name: MatchFilename
  description: |
    If set we use the filename to detect the type of sqlite file.
    When unset we use heristics (slower)
  type: bool
  default: Y

- name: CustomGlob
  description: Specify this glob to select other files

- name: DateAfter
  description: Timebox output to rows after this time.
  type: timestamp
  default: "1970-01-01T00:00:00Z"

- name: DateBefore
  description: Timebox output to rows after this time.
  type: timestamp
  default: "2100-01-01T00:00:00Z"

- name: FilterRegex
  description: Filter critical rows by this regex
  type: regex
  default: .

%v

- name: SQLITE_ALWAYS_MAKE_TEMPFILE
  type: bool
  default: Y

- name: AlsoUpload
  description: If specified we also upload the identified file.
  type: bool

sources:
- name: AllFiles
  query: |
    SELECT * FROM AllFiles

%v

`, self.Name, indent(self.Description, 4),
		self.encodeSpec(),
		utils.MustMarshalString(self.Category.Keys()),
		self.getParameters(),
		self.getSources())
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

func (self *Artifact) getParameters() string {
	res := []string{}
	for _, k := range self.Category.Keys() {
		default_value := "N"
		if k == "All" {
			default_value = "Y"
		}
		res = append(res, fmt.Sprintf(`
- name: %v
  description: Select targets with category %v
  type: bool
  default: %v
`, k, k, default_value))
	}
	return strings.Join(res, "\n")
}

func (self *Artifact) getSources() string {
	res := []string{}
	for _, k := range self.Spec.Sources.Keys() {
		v_any, _ := self.Spec.Sources.Get(k)
		v, ok := v_any.(api.Source)
		if !ok {
			continue
		}
		// If it is not an SQLite query at all, just pass the files
		// directly.
		if v.SQL == "" {
			res = append(res, fmt.Sprintf(`
- name: %v
  query: |
    LET Rows = SELECT * FROM FilterFile(SourceName=%q)
%v
`, k, k, indent(v.VQL, 4)))

		} else {
			res = append(res, fmt.Sprintf(`
- name: %v
  query: |
    LET Rows = SELECT * FROM ApplyFile(SourceName=%q)
%v
`, k, k, indent(v.VQL, 4)))
		}
	}
	return strings.Join(res, "\n")
}
