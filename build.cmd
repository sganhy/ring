set GOOS=windows
set GOARCH=amd64
set GOROOT=c:\go\
rem go get vendor
rem go mod vendor
go build -o ./ ./cmd/main.go
pause