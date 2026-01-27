package core

import (
	"fmt"
	"regexp"
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

// Enforce unique contraints on string columns
type IndexItem struct {
	Value string
	Pk    int
}

// Less implements btree.Item for the Unique index
func (i IndexItem) Less(b btree.Item) bool {
	return i.Value < b.(IndexItem).Value
}

type Table struct {
	Name    string
	Columns []string

	PkColumn string

	// Index Stores the actual data sorted by primary key
	Index *btree.BTree

	// UniqueIndexes maps column name to its unique index BTree
	UniqueIndexes map[string]*btree.BTree
	Constraint    []Constraint
}

var tables = map[string]*btree.BTree{}
var catalog = map[string]*Table{}

func ShowTables() {
	fmt.Println("Existing tables:")
	for name, tbl := range catalog {
		log.Info().Msgf("Table: %s, Columns: %v", name, tbl.Columns)
	}
}
func VerifyTable(tableName string) {
	for name, tbl := range catalog {
		log.Info().Msgf("Catalog contains table: %s", name)

		log.Info().Msgf("Table: %s", name)
		log.Info().Msgf("Primary Key: %s", tbl.PkColumn)
		log.Info().Msgf("Unique Indexes count: %d", len(tbl.UniqueIndexes))
		for colName := range tbl.UniqueIndexes {
			log.Info().Msgf("Unique Index on column: %s", colName)
		}
	}

}
func ExecuteStatement(stmt sqlparser.Statement, prStr string) {
	// bt := btree.New(3)
	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		log.Info().Msg("Executing INSERT statement")
		tableName := stmt.Table.Name.String()
		// Prefer the catalog's index if the table exists in catalog,
		// otherwise ensure there's a btree in `tables` map.
		var tree *btree.BTree
		if meta, ok := catalog[tableName]; ok && meta.Index != nil {
			tree = meta.Index
		} else {
			if tables[tableName] == nil {
				tables[tableName] = btree.New(3)
			}
			tree = tables[tableName]
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
		if err := enforceForeignKey(catalog[tableName], Row{Key: key, Data: rowData}, catalog); err != nil {
			log.Error().Msgf("Foreign key constraint violation: %s", err)
			return
		}
		tree.ReplaceOrInsert(Row{Key: key, Data: rowData})
	case *sqlparser.Select:
		log.Info().Msg("Executing SELECT")
		tableName := stmt.From[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
		// Find the underlying btree either from catalog or tables map
		var table *btree.BTree
		if meta, ok := catalog[tableName]; ok && meta.Index != nil {
			table = meta.Index
		} else {
			table = tables[tableName]
		}
		if table == nil {
			log.Warn().Msg("Table not found")
			return
		}

		join := stmt.From[0].(*sqlparser.JoinTableExpr)
		left := join.LeftExpr.(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
		right := join.RightExpr.(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
		results := executeJoin(catalog[left], catalog[right], join.Condition.On.(*sqlparser.ComparisonExpr))

		for _, row := range results {
			log.Info().Msgf("Joined Row: %+v", row)
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

			// Initialize the new table
			newTable := &Table{
				Name:          tableName,
				Columns:       []string{},
				Index:         btree.New(2),
				UniqueIndexes: make(map[string]*btree.BTree),
			}

			if stmt.TableSpec != nil {
				for _, col := range stmt.TableSpec.Columns {
					newTable.Columns = append(newTable.Columns, col.Name.String())

				}
				for _, index := range stmt.TableSpec.Indexes {

					if index.Info.Primary {
						colName := index.Columns[0].Column.String()
						newTable.PkColumn = colName
						log.Info().Msgf("Set Primary Key to Column: %s", colName)
					}
					if index.Info.Unique {

						colName := index.Columns[0].Column.String()
						// Create a separate BTree to track uniqueness for this column
						newTable.UniqueIndexes[colName] = btree.New(2)
						log.Info().Msgf("Created Unique Index for column: %s", colName)

					}

				}

			}
			catalog[tableName] = newTable
			log.Info().Msgf("Table %s created with columns %v", tableName, newTable.Columns)

		}

		if stmt.Action == sqlparser.DropStr {
			tableName := stmt.Table.Name.String()
			if _, ok := catalog[tableName]; ok {
				delete(catalog, tableName)
				// also remove from runtime tables map
				delete(tables, tableName)
				log.Info().Msgf("Table %s dropped", tableName)
			} else {
				log.Warn().Msgf("Table %s does not exist", tableName)
			}
		}
		if stmt.Action == sqlparser.AlterStr {
			tableName := stmt.Table.Name.String()
			log.Warn().Msgf("Alter table operation detected on %s ", tableName)

			alterOp, err := parseAlterDetails(prStr, tableName)
			if err != nil {
				log.Error().Msgf("Failed to parse ALTER statement: %s", err)
				return
			}

			log.Info().Msgf("ALTER operation parsed: %+v", alterOp)

			tblMeta, ok := catalog[tableName]
			if !ok {
				log.Error().Msgf("Table %s not found in catalog", tableName)
				return
			}

			switch alterOp.Type {
			case "ADD":
				if alterOp.ObjectType == "COLUMN" {
					// add column to catalog if not exists
					col := alterOp.ColumnName
					found := false
					for _, c := range tblMeta.Columns {
						if c == col {
							found = true
							break
						}
						found = false
					}

					if !found {
						tblMeta.Columns = append(tblMeta.Columns, col)
						log.Info().Msgf("Added column %s to table %s", col, tableName)
					} else {
						log.Warn().Msgf("Column %s already exists on table %s", col, tableName)
					}

					// update existing rows to include the new column with empty value
					tree := tables[tableName]
					if tree != nil {
						var toReplace []Row
						tree.Ascend(func(i btree.Item) bool {
							r := i.(Row)
							if r.Data == nil {
								r.Data = map[string]string{}
							}
							if _, has := r.Data[col]; !has {
								r.Data[col] = ""
							}
							toReplace = append(toReplace, r)
							return true
						})
						for _, nr := range toReplace {
							tree.ReplaceOrInsert(nr)
						}
					}
				}

			case "DROP":
				if alterOp.ObjectType == "COLUMN" {
					col := alterOp.ColumnName
					// remove from catalog
					newCols := []string{}
					for _, c := range tblMeta.Columns {
						if c != col {
							newCols = append(newCols, c)
						}
					}
					if len(newCols) == len(tblMeta.Columns) {
						log.Warn().Msgf("Column %s does not exist on table %s", col, tableName)
					} else {
						tblMeta.Columns = newCols
						log.Info().Msgf("Dropped column %s from table %s", col, tableName)
					}

					// remove the column key from existing rows
					tree := tables[tableName]
					if tree != nil {
						var toReplace []Row
						tree.Ascend(func(i btree.Item) bool {
							r := i.(Row)
							if r.Data != nil {
								delete(r.Data, col)

							}
							toReplace = append(toReplace, r)
							return true
						})
						for _, nr := range toReplace {
							tree.ReplaceOrInsert(nr)
						}
					}
				}
			default:
				log.Warn().Msgf("Unhandled ALTER operation: %+v", alterOp)
			}

			return

		}

	case *sqlparser.Show:
		if stmt.Type == "tables" {
			ShowTables()
		}
		if stmt.Type == "index" {
			tableName := stmt.OnTable.Name.String()

			VerifyTable(tableName)
		}
	case *sqlparser.Delete:
		log.Info().Msg("Executing DELETE")
		tableName := stmt.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
		comp := stmt.Where.Expr.(*sqlparser.ComparisonExpr)
		val := comp.Right.(*sqlparser.SQLVal)
		key, _ := strconv.Atoi(string(val.Val))

		// Prefer catalog index if present
		if meta, ok := catalog[tableName]; ok && meta.Index != nil {
			meta.Index.Delete(Row{Key: key})
		} else if tables[tableName] != nil {
			tables[tableName].Delete(Row{Key: key})
		} else {
			log.Warn().Msgf("Table %s not found for DELETE", tableName)
			break
		}
		log.Info().Msgf("Deleted row with key %d", key)
	case *sqlparser.Update:
		tableName := stmt.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
		tbl, ok := catalog[tableName]
		if !ok {
			log.Warn().Msgf("Table %s not found", tableName)
			break
		}

		// Extract WHERE key (assume primary key = first column)
		var key int
		if stmt.Where != nil {
			comp := stmt.Where.Expr.(*sqlparser.ComparisonExpr)
			val := comp.Right.(*sqlparser.SQLVal)
			key, _ = strconv.Atoi(string(val.Val))
		} else {
			log.Warn().Msg("UPDATE without WHERE not supported yet")
			break
		}

		// Find existing row
		item := tbl.Index.Get(Row{Key: key})
		if item == nil {
			log.Warn().Msgf("Row with key %d not found", key)
			break
		}
		row := item.(Row)

		// Apply updates
		for _, expr := range stmt.Exprs {
			colName := expr.Name.Name.String()
			if v, ok := expr.Expr.(*sqlparser.SQLVal); ok {
				row.Data[colName] = string(v.Val)
			}
		}

		// Replace updated row back into BTree
		tbl.Index.ReplaceOrInsert(row)
		log.Info().Msgf("Updated row in %s: %+v", tableName, row.Data)

	}

}

type AlterOperation struct {
	Type       string // "ADD", "DROP", etc.
	ObjectType string // "COLUMN", "INDEX", etc.
	ColumnName string // for column operations
	Table      string // table name from sqlparser
}

func parseAlterDetails(sql string, tableName string) (*AlterOperation, error) {
	// Simple regex-based parsing for common ALTER patterns
	// This is a basic example - you'd need more sophisticated parsing for production

	addColumnRegex := regexp.MustCompile(`(?i)ALTER\s+TABLE\s+\w+\s+ADD\s+(?:COLUMN\s+)?(\w+)`)
	dropColumnRegex := regexp.MustCompile(`(?i)ALTER\s+TABLE\s+\w+\s+DROP\s+(?:COLUMN\s+)?(\w+)`)

	if matches := addColumnRegex.FindStringSubmatch(sql); matches != nil {
		return &AlterOperation{
			Type:       "ADD",
			ObjectType: "COLUMN",
			ColumnName: matches[1],
			Table:      tableName,
		}, nil
	}

	if matches := dropColumnRegex.FindStringSubmatch(sql); matches != nil {
		return &AlterOperation{
			Type:       "DROP",
			ObjectType: "COLUMN",
			ColumnName: matches[1],
			Table:      tableName,
		}, nil
	}

	return &AlterOperation{
		Type:  "UNKNOWN",
		Table: tableName,
	}, nil
}
