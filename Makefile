templ:
	@templ generate

build_linux: templ
	@GOOS=linux GOARCH=amd64 go build -o bin/app.out .
build_windows: templ
	@GOOS=windows GOARCH=amd64 go build -o bin/app.exe .

test:
	@go test -v ./...
	
dev:
	air

run: templ
	@go run .




