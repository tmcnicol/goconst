run:
	go run ./... --type eventType,eventType2 ./...

build:
	go build -o goconst

generate:
	go generate -v ./testpackage/...
