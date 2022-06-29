package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/whaios/showdocdb/log"
)

func init() {
	mysqlFactory = newMySQL
}

func newMySQL(host, user, pwd, dbName string) Query {
	// "username:password@tcp(127.0.0.1:3306)/dbname?param=value"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pwd, host, dbName)
	return &mySQL{
		dsn:    dsn,
		dbName: dbName,
	}
}

type mySQL struct {
	dsn    string
	dbName string
	db     *sql.DB
}

func (d *mySQL) Open() error {
	db, err := sql.Open(MySQL, d.dsn)
	d.db = db
	return err
}

func (d *mySQL) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *mySQL) Query() ([]*Table, error) {
	tbs, err := d.queryTable()
	if err != nil {
		return nil, err
	}
	cols, err := d.queryColumn()
	if err != nil {
		return nil, err
	}
	tbmap := make(map[string]*Table, len(tbs))
	{
		for _, tb := range tbs {
			tbmap[tb.Name] = tb
			log.Debug("table: %s", tb.Name)
		}
	}
	for _, col := range cols {
		tb, ok := tbmap[col.TableName]
		if !ok {
			continue
		}
		tb.Columns = append(tb.Columns, col)
		log.Debug("column: %s.%s", col.TableName, col.Name)
	}
	return tbs, err
}

func (d *mySQL) queryTable() ([]*Table, error) {
	querySql := "SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?"
	rows, err := d.db.Query(querySql, d.dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tbs := make([]*Table, 0)
	for rows.Next() {
		var name, comment sql.NullString
		if err = rows.Scan(&name, &comment); err != nil {
			return tbs, err
		}
		tbs = append(tbs, &Table{Name: name.String, Comment: comment.String})
	}
	return tbs, nil
}

func (d *mySQL) queryColumn() ([]*Column, error) {
	querySql := `SELECT TABLE_NAME, COLUMN_NAME, COLUMN_DEFAULT, IS_NULLABLE, COLUMN_TYPE, COLUMN_COMMENT
				 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ?
				 ORDER BY ORDINAL_POSITION`
	rows, err := d.db.Query(querySql, d.dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]*Column, 0)
	for rows.Next() {
		var tableName, colName, colDefault sql.NullString
		var isNull, colType, comment string
		if err = rows.Scan(&tableName, &colName, &colDefault, &isNull, &colType, &comment); err != nil {
			return cols, err
		}
		cols = append(cols, &Column{
			TableName: tableName.String,
			Name:      colName.String,
			Default:   colDefault.String,
			IsNull:    isNull,
			Type:      colType,
			Comment:   comment,
		})
	}
	return cols, nil
}
