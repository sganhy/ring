package databaseprovider

type DatabaseProvider int8

const (
	Oracle     DatabaseProvider = 1
	PostgreSql DatabaseProvider = 2
	MySql      DatabaseProvider = 3
	Influx     DatabaseProvider = 4
)

/*
	Oracle     DatabaseProvider = "godror"
	PostgreSql DatabaseProvider = "postgres"
	MySql      DatabaseProvider = "mysql"
	Influx     DatabaseProvider = "influx"
*/

func (provider *DatabaseProvider) ToString() string {
	prov := *provider
	switch prov {
	case PostgreSql:
		return "postgres"
	case MySql:
		return "mysql"
	}

	return ""
}
