//go:build cgo
// +build cgo

// 因为 go-sqlite3 库是一个 cgo 库，编译代码时需要 gcc 环境。

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/whaios/showdocdb/log"
)

func init() {
	sqliteFactory = newSqlite
}

func newSqlite(file string) Query {
	// "file:test.db?cache=shared&mode=memory"
	dsn := fmt.Sprintf("file:%s", file)
	return &sqlite{
		dsn: dsn,
	}
}

type sqlite struct {
	dsn string
	db  *sql.DB
}

func (d *sqlite) Open() error {
	db, err := sql.Open(SQlite, d.dsn)
	d.db = db
	return err
}

func (d *sqlite) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *sqlite) Query() ([]*Table, error) {
	tbs, err := d.queryTable()
	if err != nil {
		return nil, err
	}
	for _, tb := range tbs {
		log.Debug("table: %s", tb.Name)
		cols, err := d.queryColumn(tb.Name)
		if err != nil {
			return nil, err
		}
		tb.Columns = cols
	}
	return tbs, nil
}

func (d *sqlite) queryTable() ([]*Table, error) {
	querySql := "SELECT name FROM sqlite_master WHERE type ='table' ORDER BY name"
	rows, err := d.db.Query(querySql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tbs := make([]*Table, 0)
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return tbs, err
		}
		tbs = append(tbs, &Table{Name: name})
	}
	return tbs, nil
}

func (d *sqlite) queryColumn(tableName string) ([]*Column, error) {
	querySql := fmt.Sprintf("PRAGMA table_info([%s])", tableName)
	rows, err := d.db.Query(querySql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]*Column, 0)
	for rows.Next() {
		var dflt_value sql.NullString
		var cid, name, tpe, notnull, pk string
		if err = rows.Scan(&cid, &name, &tpe, &notnull, &dflt_value, &pk); err != nil {
			return cols, err
		}
		col := &Column{
			TableName: tableName,
			Name:      name,
			Default:   dflt_value.String,
			IsNull:    notnull,
			Type:      tpe,
			Comment:   "",
		}
		cols = append(cols, col)
		log.Debug("column: %s.%s", col.TableName, col.Name)
	}
	return cols, nil
}
