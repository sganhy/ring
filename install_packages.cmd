set GOOS=windows
set GOARCH=amd64
set GOROOT=c:\go\
rem set https_proxy=exp..:mdp@proxyi.msnet....
rem go mod init .
go get github.com/DATA-DOG/go-sqlmock
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get golang.org/x/text
go get github.com/denisenkom/go-mssqldb
pause