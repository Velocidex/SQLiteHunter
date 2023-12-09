package testing

import (
	"os"
	"path/filepath"

	"github.com/Velocidex/ordereddict"
	"github.com/alecthomas/assert"
	"github.com/sebdah/goldie/v2"
)

func (self *SQLiteHunterTestSuite) TestEdgeWebAssistDatabase() {
	t := self.T()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	golden := ordereddict.NewDict()
	argv := []string{
		"--definitions", "../output",
		"query", `LET S = scope()

SELECT *, "" AS OSPath
FROM Artifact.Generic.Forensic.SQLiteHunter(
   FilterRegex=S.FilterRegex,
   MatchFilename=FALSE, All=FALSE, Edge=TRUE, CustomGlob=CustomGlob)
`, "--env", "CustomGlob=" + filepath.Join(cwd, "../test_files/Edge/WebAssistDatabase")}

	out, err := runWithArgs(argv, "--env", "FilterRegex=Audio")
	golden.Set("FilterRegex=Audio", filterOut(out))

	g := goldie.New(t,
		goldie.WithFixtureDir("fixtures"),
		goldie.WithNameSuffix(".golden"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)
	g.AssertJson(t, "TestEdgeWebAssistDatabase", golden)
}
