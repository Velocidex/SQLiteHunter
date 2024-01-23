package testing

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/assert"
)

func (self *SQLiteHunterTestSuite) TestGolden() {
	t := self.T()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	test_files, err := filepath.Abs(filepath.Join(cwd, "../test_files/"))
	assert.NoError(t, err)

	config_file, err := filepath.Abs(filepath.Join(cwd, "test.config.yaml"))
	assert.NoError(t, err)

	test_cases, err := filepath.Abs(filepath.Join(cwd, "./testcases/"))
	assert.NoError(t, err)

	argv := []string{
		"--definitions", "../output",
		"--config", config_file,
		"golden", "--env", "testFiles=" + test_files, test_cases,
	}

	out, err := runWithArgs(argv)
	assert.NoError(t, err, string(out))

	fmt.Println(string(out))
}
