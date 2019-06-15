all:
	GOOS=windows GOARCH=386 go build -o binaries/cfip-win32.exe cfip.go
	GOOS=windows GOARCH=amd64 go build -o binaries/cfip-win64.exe cfip.go
	GOOS=linux GOARCH=386 go build  -o binaries/cfip-linux32 cfip.go
	GOOS=linux GOARCH=amd64 go build -o binaries/cfip-linux64 cfip.go
	GOOS=darwin GOARCH=386 go build -o binaries/cfip-osx32 cfip.go
	GOOS=darwin GOARCH=amd64 go build -o binaries/cfip-osx64 cfip.go
