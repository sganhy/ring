package databaseprovider

type DatabaseProvider int8

const (
	Oracle     DatabaseProvider = 1
	PostgreSql DatabaseProvider = 2
	MySql      DatabaseProvider = 3
	Influx     DatabaseProvider = 4
	SqlServer  DatabaseProvider = 5
	Undefined  DatabaseProvider = 101
)

/*
	Oracle     DatabaseProvider = "godror"
	PostgreSql DatabaseProvider = "postgres"
	MySql      DatabaseProvider = "mysql"
	Influx     DatabaseProvider = "influx"
*/

func (provider *DatabaseProvider) String() string {
	prov := *provider
	switch prov {
	case PostgreSql:
		return "postgres"
	case MySql:
		return "mysql"
	case SqlServer:
		return "sqlserver"
	}
	return ""
}

func GetDatabaseProviderById(providerId int) DatabaseProvider {
	if providerId <= 127 && providerId >= -128 {
		var newId = DatabaseProvider(providerId)
		if newId == Oracle || newId == PostgreSql || newId == MySql || newId == Influx {
			return newId
		}
	}
	return Undefined
}
