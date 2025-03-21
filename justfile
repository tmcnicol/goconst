build:
	go build -o goconst

generate:
	go generate -v ./testpackage/...
