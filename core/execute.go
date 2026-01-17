package core

import (
	"fmt"
	"strconv"

	"github.com/google/btree"
	"github.com/rs/zerolog/log"
	"github.com/xwb1989/sqlparser"
)

type Row struct {
	Key  int
	Data map[string]string
}

func (r Row) Less(b btree.Item) bool {
	return r.Key < b.(Row).Key
}

var tables = map[string]*btree.BTree{}

func ExecuteStatement(stmt sqlparser.Statement) {
	// bt := btree.New(3)
	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		log.Info().Msg("Executing INSERT statement")
		tableName := stmt.Table.Name.String()
		if tables[tableName] == nil {
			tables[tableName] = btree.New(3)
		}
		rows := stmt.Rows.(sqlparser.Values)
		firstRow := rows[0]
		// Assuming the first column is a primary key
		keyExpr := firstRow[0].(*sqlparser.SQLVal)
		key, _ := strconv.Atoi(string(keyExpr.Val))
		rowData := make(map[string]string)
		for i, expr := range firstRow {
			if v, ok := expr.(*sqlparser.SQLVal); ok {
				rowData[fmt.Sprintf("col%d", i)] = string(v.Val)
			}
		}
		tables[tableName].ReplaceOrInsert(Row{Key: key, Data: rowData})
	case *sqlparser.Select:
		log.Info().Msg("Executing SELECT")
		tableName := stmt.From[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
		comp := stmt.Where.Expr.(*sqlparser.ComparisonExpr)
		val := comp.Right.(*sqlparser.SQLVal)
		key, _ := strconv.Atoi(string(val.Val))

		item := tables[tableName].Get(Row{Key: key})
		if item != nil {
			log.Info().Msgf("Found row: %+v", item.(Row))
		} else {
			log.Info().Msg("Row not found")
		}

	}
}
