VERSION=`go run main.go -version`
default:
	env GOOS=windows GOARCH=amd64 go build -o build/get-ssm-params.$(VERSION).windows-amd64.exe
	env GOOS=linux GOARCH=amd64 go build -o build/get-ssm-params.$(VERSION).linux-amd64
	env GOOS=darwin GOARCH=amd64 go build -o build/get-ssm-params.$(VERSION).darwin-amd64

build: default

tag:
	git tag $(VERSION)
