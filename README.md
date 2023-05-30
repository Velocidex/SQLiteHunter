# SQLite Hunter

This repository maintains the source for the
`Generic.Forensic.SQLiteHunter` VQL artifact. This artifact is desiged
to be an efficient and mostly automated artifact to analyse and
collect SQLite based artifacts from various applications on the
endpoint.

The produced artifact is self contained and can be loaded into
Velociraptor (https://docs.velociraptor.app) to hunt quickly and
efficiently across a large number of endpoints.

SQLite has become the de-facto standard for storing application data,
in many types of applications:

- Web Browsers
- Operating Systems
- Various applications, such as iMessage, TCC etc

## How do we hunt for SQLite files?

Compiling this repository with produce a single artifact called
`Generic.Forensic.SQLiteHunter` with multiple sources. Each artifact
source targets a single aspect of a single application and is applied
to a single SQLite file.

Since SQLite files can be used for many different applications we use
three phases; Collection of SQLite, Identification of the SQLite and
finally analysis of the file:

1. In the first phase we collect prospective SQLite files for the
   desired targets based on glob expressions to quickly locate the
   usual places these are stored. For example, looking for Chrome
   Browser History files typically these are stored in
   `C:\Users\*\AppData\{Roaming,Local}/Google/Chrome/User Data`.

   By employing targetted glob expressions we can quickly locate
   relevat files. However the user can also provide a generic glob
   expression for us to use other files (e.g. files collected by some
   other means off a different system).

2. Since different applications use SQLite in different ways, we want
   to have specialized treatment for each application type -
   extracting relevant data and potentially enriching it for enhanced
   analysis.

   Looking at the prospective files found in stage 1 we need to
   classify each file to a specific type. Each artifact source targets
   a specific application and sqlite file. In order to identify the
   file the source runs the `SQLiteIdentifyQuery` on the sqlite file
   (as described below).

   In the common mode we can use the filename itself to quickly
   classify the file this is a shortcut to speed things up. If the
   files could have been renamed, you can specify `MatchFilename` to
   be false in which case only the SQLiteIdentifyQuery method will be
   used (this will be slower).

3. Once a file is identified as belonging to a particular application,
   the artifact source can run the specified SQL on the file. Since
   pure SQL is very limited in the type of data it can use and it is
   also harder to use, the output is enriched further via a VQL query.

   Being able to apply VQL to the output of the SQL query makes
   crafting the SQL much easier (for example timestamp conversions are
   much easier in VQL than SQL). Additionally the VQL can be used to
   enrich the data from other sources (e.g. geoip etc).

## How is this repository organized?

The main logic is stored in YAML definitions stored in the `definitions` directory:

1. Name: This is the first part of the artifact source name that will be produced.

2. Author,Email, Reference: Self explanatory.

3. SQLiteIdentifyQuery and SQLiteIdentifyValue: To test if the SQLite
   file is one that should be targeted by this definition,
   Velociraptor will run the SQLiteIdentifyQuery which should produce
   one row and one columns called `Check`. The value in this column
   will be checked against SQLiteIdentifyValue to determine if the
   file qualifies for this map.

4. Categories: A list of keywords that can be used to limit the
   collection to only certain categories. Note that some categories
   may overlap (e.g. Chrome and Browser)

5. FilenameRegex: A regex that can be used to the filename to shortcut
   identification of the file when "MatchFilename" is enabled. NOTE
   that we do this in addition to the `SQLiteIdentifyQuery` so it is
   only an optimization to speed up processing.

6. Globs: A list of glob expression. This list can be interpolated
   with the globs in `config.yaml`

7. Sources: This is a list of source definitions that will be
   converted to an artifact source. Each of these may contain:

   * Name: If more than one source is specified in a definition, they
     can have a name. This name will be used together with the main
     definition source to build the Artifact source name in the final
     artifact.
   * VQL: This is a VQL query that will be used to build the artifact
     source. The query must end with `SELECT .... FROM Rows`
   * SQL: This is the SQL query that will be applied to the SQLite
     file. Generally it is easier to apply enrichment, processing etc
     in the VQL so the SQL query can be much simpler.
   * SQLiteIdentifyQuery and SQLiteIdentifyValue - if these appear
     within the source they will override the definition. This allows
     for different sources to be written for different versions of the
     SQLite tables.
