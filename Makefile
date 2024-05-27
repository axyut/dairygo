templ:
	@templ generate

build_linux: templ
	@GOOS=linux GOARCH=amd64 go build -o bin/app.out .
build_windows: templ
	@GOOS=windows GOARCH=amd64 go build -o bin/app.exe .
build_darwin: templ
	@GOOS=darwin GOARCH=amd64 go build -o bin/app .

test:
	@go test -v ./...
	
dev:
	air

run: templ
	@go run .

install:
	@go mod tidy




