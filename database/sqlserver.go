package database

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/whaios/showdocdb/log"
)

func init() {
	sqlServerFactory = newSQLServer
}

func newSQLServer(host, user, pwd, dbName string) Query {
	// "sqlserver://username:password@host/instance?database=value"
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", user, pwd, host, dbName)
	return &sqlServer{
		dsn:    dsn,
		dbName: dbName,
	}
}

type sqlServer struct {
	dsn    string
	dbName string
	db     *sql.DB
}

func (d *sqlServer) Open() error {
	db, err := sql.Open(SQLServer, d.dsn)
	d.db = db
	return err
}

func (d *sqlServer) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *sqlServer) Query() ([]*Table, error) {
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

func (d *sqlServer) queryTable() ([]*Table, error) {
	querySql := fmt.Sprintf(`
			SELECT TABLE_NAME
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_TYPE = 'BASE TABLE' AND TABLE_CATALOG = '%s' ORDER BY TABLE_NAME`,
		d.dbName)
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

func (d *sqlServer) queryColumn() ([]*Column, error) {
	querySql := fmt.Sprintf(`
			SELECT C.TABLE_NAME, C.COLUMN_NAME, C.COLUMN_DEFAULT, C.IS_NULLABLE, 
				C.DATA_TYPE + 
					CASE 
					WHEN C.CHARACTER_MAXIMUM_LENGTH IS NOT NULL THEN '('+ CONVERT(VARCHAR(10), C.CHARACTER_MAXIMUM_LENGTH) +')' 
					WHEN C.NUMERIC_SCALE > 0 THEN '('+ CONVERT(VARCHAR(10), C.NUMERIC_PRECISION) +', ' + CONVERT(VARCHAR(10), C.NUMERIC_SCALE) +')'
					WHEN C.NUMERIC_PRECISION > 0 THEN '('+ CONVERT(VARCHAR(10), C.NUMERIC_PRECISION) +')'
					ELSE '' END 
				AS DATA_TYPE,
				ISNULL(EP.COLUMN_COMMENT, '') AS COLUMN_COMMENT
			FROM INFORMATION_SCHEMA.COLUMNS AS C
			LEFT JOIN (
				SELECT TB.name AS TABLE_NAME, COL.name AS COLUMN_NAME, EP.value AS COLUMN_COMMENT
				FROM sys.extended_properties AS EP, sys.tables AS TB, sys.columns AS COL
				WHERE EP.major_id = TB.object_id AND EP.minor_id = COL.column_id AND TB.object_id = COL.object_id
			) AS EP ON EP.TABLE_NAME = C.TABLE_NAME AND EP.COLUMN_NAME = C.COLUMN_NAME
			WHERE TABLE_CATALOG = '%s'
			ORDER BY TABLE_NAME, ORDINAL_POSITION`,
		d.dbName)
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
