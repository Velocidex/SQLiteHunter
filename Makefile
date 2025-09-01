all: build windows

build:
	go build -o sqlitehunter_compiler ./bin/*.go

windows:
	GOOS=windows GOARCH=amd64 \
	go build -o sqlitehunter_compiler.exe ./bin/*.go

compile: FORCE
	go run ./bin/ compile ./definitions/ ./output/SQLiteHunter.yaml --output_zip ./output/SQLiteHunter.zip --index ./docs/content/docs/rules/index.json

golden: compile
	./testing/velociraptor.bin --definitions ./output --config ./testing/test.config.yaml golden --env testFiles=`pwd`/test_files ./testing/testcases -v --filter=${GOLDEN}


FORCE:
