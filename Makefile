all: build windows

build:
	go build -o sqlitehunter_compiler ./bin/*.go

windows:
	GOOS=windows GOARCH=amd64 \
	go build -o sqlitehunter_compiler.exe ./bin/*.go

compile: FORCE
	./sqlitehunter_compiler -config ./config.yaml -definition_directory ./definitions > output/SQLiteHunter.yaml

FORCE:
