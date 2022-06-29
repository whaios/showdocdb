package database

import "strings"

// Query 数据库查询器
type Query interface {
	Open() error
	Close() error
	Query() ([]*Table, error)
}

// Table 数据表结构
type Table struct {
	Name    string
	Comment string

	Columns []*Column
}

// Markdown 生成 markdown 内容
func (t *Table) Markdown() string {
	buf := strings.Builder{}
	if t.Comment != "" {
		buf.WriteString("- " + t.Comment)
	}
	buf.WriteString("\n\n")
	buf.WriteString("| 字段 | 类型 | 允许空 | 默认 | 注释 |\n")
	buf.WriteString("| :---- | :---- | :---- | ---- | ---- |\n")
	for _, col := range t.Columns {
		if col.Comment == "" {
			col.Comment = "无"
		}
		buf.WriteString("| " + col.Name)
		buf.WriteString(" | " + col.Type)
		buf.WriteString(" | " + col.IsNull)
		buf.WriteString(" | " + col.Default)
		buf.WriteString(" | " + col.Comment)
		buf.WriteString(" |\n")
	}
	return buf.String()
}

type Column struct {
	TableName string
	Name      string
	Comment   string
	Default   string
	IsNull    string
	Type      string
}
