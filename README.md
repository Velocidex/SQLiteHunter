# SQLite Hunter

This repository maintains the source for the
`Generic.Forensic.SQLiteHunter` VQL artifact. This artifact is
designed to be an efficient and mostly automated artifact to analyze
and collect SQLite based artifacts from various applications on the
endpoint.

The produced artifact is self contained and can be loaded into
Velociraptor (https://docs.velociraptor.app) to hunt quickly and
efficiently across a large number of endpoints.

SQLite has become the de-facto standard for storing application data,
in many types of applications:

- Web Browsers, e.g. Chrome, Firefox, Opera, Edge
- Operating Systems
- Various applications, such as iMessage, TCC etc

## How do we hunt for SQLite files?

Compiling this repository will produce a single artifact called
`Generic.Forensic.SQLiteHunter` with multiple sources. Each artifact
source targets a single aspect of a single application and is applied
to a single SQLite file.

Since SQLite files can be used for many different applications we use
three phases; Collection of SQLite files, Identification of the SQLite
application based on the file, and finally analysis of the file:

1. In the first phase we collect prospective SQLite files for the
   desired targets based on glob expressions to quickly locate the
   usual places these are stored. For example, looking for Chrome
   Browser History files typically these are stored in
   `C:\Users\*\AppData\{Roaming,Local}\Google\Chrome\User Data`.

   By employing targeted glob expressions we can quickly locate
   relevant files. However the user can also provide a generic glob
   expression for us to use other files (e.g. files collected by some
   other means off a different system).

2. Since different applications use SQLite in different ways, we want
   to have specialized treatment for each application type -
   extracting relevant data and potentially enriching it for enhanced
   analysis.

   Looking at the prospective files found in stage 1 we need to
   classify each file to a specific type. Each artifact source targets
   a specific application and SQLite file. In order to identify the
   file the source runs the `SQLiteIdentifyQuery` on the SQLite file
   (as described below).

   In the common mode we can use the filename itself to quickly
   classify the file this is a shortcut to speed things up. If the
   files could have been renamed, you can specify `MatchFilename` to
   be false in which case only the `SQLiteIdentifyQuery` method will be
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

The main logic is stored in YAML definitions stored in the
`definitions` directory:

1. `Name`: This is the first part of the artifact source name that
   will be produced.

2. `Author`,`Email`, `Reference`: Self explanatory.

3. `SQLiteIdentifyQuery` and `SQLiteIdentifyValue`: To test if the SQLite
   file is one that should be targeted by this definition,
   Velociraptor will run the SQLiteIdentifyQuery which should produce
   one row and one columns called `Check`. The value in this column
   will be checked against SQLiteIdentifyValue to determine if the
   file qualifies for this map.

4. `Categories`: A list of keywords that can be used to limit the
   collection to only certain categories. Note that some categories
   may overlap (e.g. Chrome and Browser)

5. `FilenameRegex`: A regex that can be used to the filename to shortcut
   identification of the file when `MatchFilename` is enabled. NOTE
   that we do this in addition to the `SQLiteIdentifyQuery` so it is
   only an optimization to speed up processing.

6. `Globs`: A list of glob expression. This list can be interpolated
   with the globs in `config.yaml`

7. `Sources`: This is a list of source definitions that will be
   converted to an artifact source. Each of these may contain:

   * `Name`: If more than one source is specified in a definition, they
     can have a name. This name will be used together with the main
     definition source to build the Artifact source name in the final
     artifact.
   * `VQL`: This is a VQL query that will be used to build the artifact
     source. The query must end with `SELECT .... FROM Rows`
   * `SQL`: This is the SQL query that will be applied to the SQLite
     file. Generally it is easier to apply enrichment, processing etc
     in the VQL so the SQL query can be much simpler.
   * `SQLiteIdentifyQuery` and SQLiteIdentifyValue - if these appear
     within the source they will override the definition. This allows
     for different sources to be written for different versions of the
     SQLite tables.

## Example Development Walk Through

In the following section I will describe how to add new definitions to
the SQLiteHunter artifact with this repository. Because SQLiteHunter
is a combined artifact that operates on all targets we need to compile
the artifact each time we want to use it.

The general process for development is:

1. Obtain a test file (e.g. a SQLite file from a target system or
   similar). Store this file in the test_files directory somewhere.
2. Write a definitions file in the definitions directory (more on that later).
3. Compile the artifact using `make compile` or just from the top level

```
./sqlitehunter_compiler > output/SQLiteHunter.yaml
```

4. Now simply collect the new target using Velociraptor directly.

Lets work through an example. For this example I will write the `Edge
Browser Navigation History` target from the SQLECmd project.

### Step 1: Get a sample file.

This targets the file `C:\Users\User\AppData\Local\Microsoft\Edge\User Data\Default\WebAssistDatabase` so I copy this file into the
test_files directory. It is highly recommended to share a sample file
in your PR in order for automated tests to be built.

### Writing the definition file.

I will start off creating a new definition file in the definitions
directory: `EdgeBrowser_NavigationHistory.yaml`

The file starts off with the common fields:
```yaml
Name: Edge Browser Navigation History
Author: Suyash Tripathi
Email: suyash.tripathi@cybercx.com.au
Reference: https://github.com/EricZimmerman/SQLECmd
```

Next I will add the `SQLiteIdentifyQuery` that Velociraptor will run
to determine if this is in fact a `WebAssistDatabase`. A good check
(which is used in the original SQLECmd map is to check if the file
contains a `navigation_history` table.

```yaml
SQLiteIdentifyQuery: |
  SELECT count(*) AS `Check`
  FROM sqlite_master
  WHERE type='table'
    AND name='navigation_history';
SQLiteIdentifyValue: 1
```

The query is expected to return 1 row.

Next I will add a new category for this definition. I will give it the
Test category for now so I can isolate just this definition during
development. Normally SQLiteHunter is designed to operate on many
targets automatically which makes it a bit harder to use in
development. This way we can just run a single target using the
`--args All=N --args Test=Y` args.

```yaml
Categories:
  - Edge
  - Test
  - Browser
```

Next we set the file matching filters. These allow Velociraptor to
identify potential files by filename which is a lot faster than having
to read and test every file. Usually the filename is expected to be
`WebAssistDatabase` and it lives in the Edge profile directories.

The Edge browser is also available on MacOS so we need to add globs
for that.

```yaml
FilenameRegex: "WebAssistDatabase"
Globs:
  - "{{WindowsChromeProfiles}}/*/WebAssistDatabase"
  - "{{MacOSChromeProfiles}}/*/WebAssistDatabase"
```

Now come the interesting part - we need to add the actual Source for
extracting the data. The SQLiteHunter artifact is structured in a two
pass form - first the SQL is run on the sqlite file, then the
resulting rows are passed through a VQL query which is able to
enrich/post process the data.

For the moment we just want to add an SQL query that will run on the
SQLite file and simply pass the VQL through unchanged.

A good start is the SQL from the SQLECmd repository:

```yaml
Sources:
- name: Navigation History
  VQL: |
    SELECT * FROM Rows
  SQL: |
    SELECT
      navigation_history.id AS ID,
      datetime(navigation_history.last_visited_time, 'unixepoch') AS 'Last Visited Time',
      navigation_history.title AS Title,
      navigation_history.url AS URL,
      navigation_history.num_visits AS VisitCount
    FROM
      navigation_history
    ORDER BY
      navigation_history.last_visited_time ASC;
```

The VQL part is a simple passthrough query while the SQL part is take
directly from the SQLECmd project.

### Testing the definition

We are now ready to test the definition. First compile it with `make
compile`, next test with Velociraptor (from the top level directory):

```
make compile  && ./velociraptor-v0.7.1-linux-amd64 --definitions ./output/ -v artifacts collect Generic.Forensic.SQLiteHunter --args CustomGlob=`pwd`/test_files/Edge/* --args All=N --args Test=Y
```

I you do not want to build the `sqlitehunter_compiler` you can just
download it from the Releases page of this repository and place it at
the top level of the repository - otherwise you can build it from
source using just `make` at the top level.

This command:
1. Uses the Velociraptor binary appropriate for the platform you are
   running on
2. Adds the `--definitions` to get Velociraptor to automatically load
   the new artifact (overriding the built in version).
3. Uses the `-v` flag to have detailed logging - you should look for
   helpful messages or errors during development.
4. Adds the `CustomGlob` parameter to force the `SQLiteHunter`
   artifact to search the `test_files` directory instead of the
   system. Leaving this out will force it to search the current system
   which may be useful as well.
5. Finally we turn all `All` processing and focus on collecting only
   the `Test` category. This can be omitted if the `CustomGlob` is
   very specific so other targets are not triggered anyway. The
   purpose of this is to just speed up the development cycle.

Let's look at some of the output on my system:

```text
[INFO] 2023-12-08T23:10:29Z Globs for category Test is /home/mic/projects/SQLiteHunter/test_files/Edge/*
[INFO] 2023-12-08T23:10:29Z Starting collection of Generic.Forensic.SQLiteHunter/AllFiles
...
[INFO] 2023-12-08T23:10:29Z sqlite: Will try to copy /home/mic/projects/SQLiteHunter/test_files/Edge/WebAssistDatabase to temp file
[INFO] 2023-12-08T23:10:29Z sqlite: Using local copy /tmp/tmp210798471.sqlite
[INFO] 2023-12-08T23:10:29Z /home/mic/projects/SQLiteHunter/test_files/Edge/WebAssistDatabase was identified as Edge Browser Navigation History_Navigation History
[INFO] 2023-12-08T23:10:29Z sqlite: removing tempfile /tmp/tmp210798471.sqlite
[INFO] 2023-12-08T23:10:29Z /home/mic/projects/SQLiteHunter/test_files/Edge/WebAssistDatabase matched by filename WebAssistDatabase
[
 {
  "OSPath": "/home/mic/projects/SQLiteHunter/test_files/Edge/WebAssistDatabase"
 },
 {
  "ID": 0,
  "Last Visited Time": "2023-08-21 02:13:07",
  "Title": "Japanese language - Wikipedia",
  "URL": "https://en.wikipedia.org/wiki/Japanese_language",
  "VisitCount": 1,
  "OSPath": "/home/mic/projects/SQLiteHunter/test_files/Edge/WebAssistDatabase"
 },
```

The first logged message shows that selecting the Test category
results in the `CustomGlob` being used (if CustomGlob is not specified
it will resolve to the globs given in the definition file.

Next we see the sqlite file being copied to a temp file (this is done
to protect the integrity of the SQLite files from changes due to
journaling etc and to avoid locked sqlite files from genering an
error).

Next we see the file is identified as an `Edge Browser Navigation
History` file based on the `SQLiteIdentifyQuery` query.

Then we see the rows generated by the SQLite file.

The output looks almost right but there is a problem - the `Last
Visited Time` timestamp is not formatted correctly as an ISO timestamp
(there is a missing timezone specifier). This is because formatting
times was done using the SQL query but this does not generate correct
timestamps.

It is generally better to use VQL to format times correctly. Lets fix
this by moving the timestamp formatting code from SQL to VQL:

```yaml
Sources:
- name: Navigation History
  VQL: |
    SELECT ID,
       timestamp(epoch=`Last Visited Time`) AS `Last Visited Time`,
       Title, URL, VisitCount
    FROM Rows

  SQL: |
    SELECT
      navigation_history.id AS ID,
      navigation_history.last_visited_time AS 'Last Visited Time',
      navigation_history.title AS Title,
      navigation_history.url AS URL,
      navigation_history.num_visits AS VisitCount
    FROM
      navigation_history
    ORDER BY
      navigation_history.last_visited_time ASC;
```

Now the output is more correct and properly formatted:
```
{
  "ID": 0,
  "Last Visited Time": "2023-08-27T08:02:34Z",
  "Title": "Microsoft Edge | What's New",
  "URL": "https://www.microsoft.com/en-us/edge/update/116?form=MT00GR\u0026channel=stable\u0026version=116.0.1938.54",
  "VisitCount": 1
},
```

### Time boxing and filtering

While this works pretty well, we lack the ability to control the
output of the artifact based on filtering or time boxing. The user may
specify the following parameters: `DateAfter`, `DateBefore` and
`FilterRegex` to narrow output in the artifact.

Each source interprets these contraints in the way that makes sense to
them. In this case we should implement time boxing based on the `Last
Visit Time` and allow the user to filter by Title and URL:

```yaml
Sources:
- name: Navigation History
  VQL: |
    SELECT ID,
       timestamp(epoch=`Last Visited Time`) AS `Last Visited Time`,
       Title, URL, VisitCount
    FROM Rows
    WHERE `Last Visited Time` > DateAfter
      AND `Last Visited Time` < DateBefore
      AND (Title, URL) =~ FilterRegex
```

You can verify these filters work by specifying the parameters on the
command line:

```
make compile  && ./velociraptor-v0.7.1-linux-amd64 --definitions ./output/ -v artifacts collect Generic.Forensic.SQLiteHunter --args CustomGlob=`pwd`/test_files/Edge/* --args All=N --args Test=Y --args FilterRegex=Audio
```
