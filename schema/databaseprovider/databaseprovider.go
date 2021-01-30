package databaseprovider

type DatabaseProvider string

const (
	Sqlite3    DatabaseProvider = "sqlite3"
	Oracle     DatabaseProvider = "godror"
	PostgreSql DatabaseProvider = "postgres"
	MySql      DatabaseProvider = "mysql"
	Influx     DatabaseProvider = "influx"
)
