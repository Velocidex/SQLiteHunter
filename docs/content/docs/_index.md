---
title: SQLiteHunter Site
date: 2023-10-15T00:14:44+10:00
bookToc: false
---

# The Velociraptor SQLiteHunter Site

This repository maintains the source for the
`Generic.Forensic.SQLiteHunter` VQL artifact. This artifact is
designed to be an efficient and mostly automated artifact to analyze
and collect SQLite based artifacts from various applications on the
endpoint.

The produced artifact is self contained and can be loaded into
Velociraptor (https://docs.velociraptor.app) to hunt quickly and
efficiently across a large number of endpoints.

You can download the latest artifact pack [as a zip
file](/SQLiteHunter.zip), or [as a YAML
file](../artifact/SQLiteHunter.yaml) and add it manually to
Velociraptor.

## Parameters

1. **RuleFilter**: If you dont want to run all the rules, you can
   filter the ones you need using this regular expression.

2. **Rules**: Alteratively, the rules may be specified one at the time
   using a multi-choice selector.

3. **MatchFilename**: Rules generally look for SQLite files using
   known filenames. If this option is unset, we relay on automatic
   detection to identify the filenames (For example, enumerate the
   tables in the SQLite file). This makes scanning much slower so by
   default this setting is enabled.

4. **CustomGlob**: Rules default to search for SQLites using known
   globs. However, if you have a bunch of SQLite files in a different
   location, you may specify the custom glob to search for files.

5. **DateAfter** and **DateBefore**: These setting allow you to time
   box the returned rows to only return items that occurred between
   the specified dates.

6. **FilterRegex**: A filter that applies on the entire row (encoded
   as JSON). This is very useful to find all relevant rows relating to
   a specific item. For example, if you want to know any rows
   accessing www.example.com you can specify this filter which will
   return records like `Visited links`, `bookmarks`, `favicons` etc.

7. **SQLITE_ALWAYS_MAKE_TEMPFILE**: By default Velociraptor will make
   a temporary copy of the SQLite file before parsing it. This ensure
   the file is not locked and can be freely accessed. If this setting
   is set to off parsing might be a lot slower as Velociraptor will
   have to contend with application locks. There is probably no reason
   to disable this.

8. **AlsoUpload**: This option also uploads the raw SQLite files.


## Artifact

<div style="max-height: 500px; overflow-y: auto; ">
<pre >
<code style="margin-top: -40px;font-size: medium;" class="language-yaml">
{{< insert "../static/artifact/SQLiteHunter.yaml" >}}
</code>
</pre>
</div>
