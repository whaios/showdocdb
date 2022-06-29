package database

import "fmt"

const (
	MySQL      = "mysql"
	PostgreSQL = "postgres"
	SQLServer  = "sqlserver"
	SQlite     = "sqlite3"
	Oracle     = "oracle"
)

var (
	mysqlFactory     func(host, user, pwd, dbName string) Query
	postgreFactory   func(host, user, pwd, dbName, schema string) Query
	sqlServerFactory func(host, user, pwd, dbName string) Query
	sqliteFactory    func(file string) Query
	oracleFactory    func(host, user, pwd, dbName string) Query
)

func NewQuery(driver, host, user, pwd, db, schema string) (query Query, err error) {
	switch driver {
	case MySQL:
		if mysqlFactory != nil {
			query = mysqlFactory(host, user, pwd, db)
		}
	case PostgreSQL:
		if postgreFactory != nil {
			query = postgreFactory(host, user, pwd, db, schema)
		}
	case SQLServer:
		if sqlServerFactory != nil {
			query = sqlServerFactory(host, user, pwd, db)
		}
	case SQlite:
		if sqliteFactory != nil {
			query = sqliteFactory(host)
		}
	case Oracle:
		if oracleFactory != nil {
			query = oracleFactory(host, user, pwd, db)
		}
	}
	if query == nil {
		err = fmt.Errorf("不支持的数据库类型 %s", driver)
	}
	return
}
