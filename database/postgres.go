package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/whaios/showdocdb/log"
)

func init() {
	postgreFactory = newPostgreSQL
}

func newPostgreSQL(host, user, pwd, dbName, schema string) Query {
	// "postgres://username:password@localhost:5432/database_name?sslmode=disable"
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pwd, host, dbName)
	if schema == "" {
		schema = "public"
	}
	return &postgreSQL{
		dsn:    dsn,
		dbName: dbName,
		schema: schema,
	}
}

type postgreSQL struct {
	dsn    string
	dbName string
	schema string // 模式，多个表的集合，默认 public
	db     *sql.DB
}

func (d *postgreSQL) Open() error {
	db, err := sql.Open(PostgreSQL, d.dsn)
	d.db = db
	return err
}

func (d *postgreSQL) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *postgreSQL) Query() ([]*Table, error) {
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

func (d *postgreSQL) queryTable() ([]*Table, error) {
	querySql := "SELECT table_name FROM information_schema.TABLES WHERE TABLE_SCHEMA = $1"
	rows, err := d.db.Query(querySql, d.schema)
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

func (d *postgreSQL) queryColumn() ([]*Column, error) {
	querySql := `SELECT table_name, column_name, column_default, is_nullable, udt_name
				 FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = $1
				 ORDER BY ordinal_position`
	rows, err := d.db.Query(querySql, d.schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols := make([]*Column, 0)
	for rows.Next() {
		var tableName, colName, colDefault sql.NullString
		var isNull, colType string
		if err = rows.Scan(&tableName, &colName, &colDefault, &isNull, &colType); err != nil {
			return cols, err
		}
		cols = append(cols, &Column{
			TableName: tableName.String,
			Name:      colName.String,
			Default:   colDefault.String,
			IsNull:    isNull,
			Type:      colType,
			Comment:   "",
		})
	}
	return cols, nil
}
