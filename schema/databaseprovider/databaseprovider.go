package databaseprovider

type DatabaseProvider int8

const (
	Oracle     DatabaseProvider = 1
	PostgreSql DatabaseProvider = 2
	MySql      DatabaseProvider = 3
	Influx     DatabaseProvider = 4
	NotDefined DatabaseProvider = 101
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
	}

	return ""
}
