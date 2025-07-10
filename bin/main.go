package main

import (
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("sqlitehunter_compiler",
		"A tool for packaging the SQLiteHunt artifact.")

	command_handlers []CommandHandler
)

type CommandHandler func(command string) bool

func main() {
	app.HelpFlag.Short('h')
	app.UsageTemplate(kingpin.CompactUsageTemplate)
	command := kingpin.MustParse(app.Parse(os.Args[1:]))

	for _, handler := range command_handlers {
		if handler(command) {
			break
		}
	}
}
