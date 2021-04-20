set GOOS=windows
set GOARCH=amd64
rem set https_proxy=exp..:mdp@proxyi.msnet.railb.be:80
rem go mod vendor
go test -v ./schema
go test -v ./data


pause