//go:build cgo

// 因为 godror 库是一个 cgo 库，编译代码时需要 gcc 环境。

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"github.com/whaios/showdocdb/log"
)

func init() {
	oracleFactory = newOracle
}

func newOracle(host, user, pwd, dbName string) Query {
	// `user="scott" password="tiger" connectString="dbhost:1521/orclpdb1"`
	dsn := fmt.Sprintf(`user="%s" password="%s" connectString="%s/%s"`, user, pwd, host, dbName)
	return &oracle{
		dsn: dsn,
	}
}

type oracle struct {
	dsn string
	db  *sql.DB
}

func (d *oracle) Open() error {
	db, err := sql.Open("godror", d.dsn)
	d.db = db
	return err
}

func (d *oracle) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *oracle) Query() ([]*Table, error) {
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

func (d *oracle) queryTable() ([]*Table, error) {
	querySql := "SELECT T.TABLE_NAME, C.COMMENTS FROM USER_TABLES T, USER_TAB_COMMENTS C WHERE T.TABLE_NAME = C.TABLE_NAME ORDER BY T.TABLE_NAME"
	rows, err := d.db.Query(querySql)
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

func (d *oracle) queryColumn() ([]*Column, error) {
	querySql := `SELECT T.TABLE_NAME, T.COLUMN_NAME, T.DATA_DEFAULT, T.NULLABLE,
       				T.DATA_TYPE || 
						CASE 
						WHEN T.CHAR_LENGTH > 0 THEN '(' || T.CHAR_LENGTH || ')'
						WHEN T.DATA_SCALE > 0 THEN '(' || T.DATA_PRECISION || ', ' || T.DATA_SCALE || ')'
						WHEN T.DATA_PRECISION > 0 THEN '(' || T.DATA_PRECISION || ')'
						ELSE '' END
					AS DATA_TYPE,
					C.COMMENTS
				 FROM USER_TAB_COLUMNS T, USER_COL_COMMENTS C
				 WHERE T.TABLE_NAME = C.TABLE_NAME AND T.COLUMN_NAME = C.COLUMN_NAME
				 ORDER BY T.TABLE_NAME, T.COLUMN_ID`
	rows, err := d.db.Query(querySql)
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
