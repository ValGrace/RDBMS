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

type Table struct {
	Name    string
	Columns []string
	Index   *btree.BTree
}

var tables = map[string]*btree.BTree{}
var catalog = map[string]*Table{}

func ShowTables() {
	fmt.Println("Existing tables:")
	for name, tbl := range catalog {
		log.Info().Msgf("Table: %s, Columns: %v", name, tbl.Columns)
	}
}

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
		table := tables[tableName]
		if table == nil {
			log.Warn().Msg("Table not found")
			return
		}

		if stmt.Where != nil {
			// SELECT with WHERE
			comp := stmt.Where.Expr.(*sqlparser.ComparisonExpr)
			val := comp.Right.(*sqlparser.SQLVal)
			key, _ := strconv.Atoi(string(val.Val))

			item := table.Get(Row{Key: key})
			if item != nil {
				log.Info().Msgf("Found row: %+v", item.(Row))
			} else {
				log.Info().Msg("Row not found")
			}
		} else {
			// SELECT without WHERE → iterate all rows
			table.Ascend(func(i btree.Item) bool {
				row := i.(Row)
				log.Info().Msgf("Row: %+v", row)
				return true
			})
		}
	case *sqlparser.DDL:
		if stmt.Action == sqlparser.CreateStr {
			tableName := stmt.NewName.Name.String()
			log.Info().Msgf("Creating table: %s", tableName)
			cols := []string{}
			if stmt.TableSpec != nil {
				for _, col := range stmt.TableSpec.Columns {
					cols = append(cols, col.Name.String())
				}
			}
			catalog[tableName] = &Table{
				Name:    tableName,
				Columns: cols,
				Index:   btree.New(3),
			}
			log.Info().Msgf("Table %s created with columns %v", tableName, cols)
		}
	case *sqlparser.Show:
		if stmt.Type == "tables" {
			ShowTables()
		}
	}
}
