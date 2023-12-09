package testing

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Velocidex/ordereddict"
	"github.com/alecthomas/assert"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	VelociraptorUrl        = "https://github.com/Velocidex/velociraptor/releases/download/v0.7.0/velociraptor-v0.7.0-4-linux-amd64-musl"
	VelociraptorBinaryPath = "./velociraptor.bin"
)

type SQLiteHunterTestSuite struct {
	suite.Suite
}

func (self *SQLiteHunterTestSuite) SetupSuite() {
	self.findAndPrepareBinary()
}

func (self *SQLiteHunterTestSuite) findAndPrepareBinary() {
	t := self.T()

	_, err := os.Lstat(VelociraptorBinaryPath)
	if err != nil {
		fmt.Printf("Downloading %v from %v\n", VelociraptorBinaryPath,
			VelociraptorUrl)
		resp, err := http.Get(VelociraptorUrl)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// Create the file
		out, err := os.OpenFile(VelociraptorBinaryPath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
		assert.NoError(t, err)
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		assert.NoError(t, err)
	}
}

func (self *SQLiteHunterTestSuite) TestFirefoxHistory() {
	t := self.T()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	golden := ordereddict.NewDict()
	argv := []string{
		"--definitions", "../output",
		"query", `LET S = scope()

SELECT LastVisitDate, URL, Description, _Source
FROM Artifact.Generic.Forensic.SQLiteHunter(
   DateBefore=S.DateBefore || "2200-10-10",
   DateAfter=S.DateAfter || "1900-01-01",
   FilterRegex=S.FilterRegex || ".",
   MatchFilename=FALSE, All=FALSE, Chrome=TRUE, CustomGlob=CustomGlob)
WHERE VisitID
`, "--env", "CustomGlob=" + filepath.Join(cwd, "../test_files/Firefox/*")}

	// This file has these dates:
	// 2020-06-27T09:29:54.51375Z
	// 2020-06-27T09:30:05.721357Z
	// 2020-06-30T05:53:37.171Z
	// 2021-02-21T08:55:10.488Z
	out, err := runWithArgs(argv)
	require.NoError(t, err, out)

	golden.Set("All Records", filterOut(out))

	out, err = runWithArgs(argv,
		"--env", "DateAfter=2021-02-20T08:55:10.488Z")
	assert.NoError(t, err, out)

	golden.Set("After 2021-02-20T08:55:10.488Z should be only 2021-02-21T08:55:10.488Z",
		filterOut(out))

	out, err = runWithArgs(argv, "--env", "DateBefore=2020-06-27T09:30:00Z")
	assert.NoError(t, err, out)

	golden.Set("DateBefore=2020-06-27T09:30:00Z should be only 2020-06-27T09:29:54.51375Z",
		filterOut(out))

	out, err = runWithArgs(argv,
		"--env", "FilterRegex=Firefox Developer Edition")
	golden.Set("FilterRegex=Firefox Developer Edition",
		filterOut(out))

	g := goldie.New(t,
		goldie.WithFixtureDir("fixtures"),
		goldie.WithNameSuffix(".golden"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)
	g.AssertJson(t, "TestArtifact", golden)
}

func TestSQLiteHunter(t *testing.T) {
	suite.Run(t, &SQLiteHunterTestSuite{})
}

func filterOut(out string) []string {
	res := []string{}

	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "OSPath") {
			res = append(res, line)
		}
	}
	return res
}

func runWithArgs(argv []string, args ...string) (string, error) {
	full_argv := append(argv, args...)

	fmt.Printf("Running %v %v\n", VelociraptorBinaryPath,
		strings.Join(full_argv, " "))
	cmd := exec.Command(VelociraptorBinaryPath, full_argv...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
