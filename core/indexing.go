package core

import (
	"github.com/rs/zerolog/log"

	"github.com/google/btree"
)

type Constraint struct {
	Type       string
	Columns    []string
	RefTable   string
	RefColumns []string
}

func enforcePrimaryKey(tbl *Table, row Row) error {
	for _, c := range tbl.Constraint {
		if c.Type == "PRIMARY" {
			keyCols := c.Columns
			tbl.Index.Ascend(func(item btree.Item) bool {
				existingRow := item.(Row)
				match := true
				for _, col := range keyCols {
					if existingRow.Data[col] != row.Data[col] {
						match = false
						break
					}
				}
				if match {
					return false
				}
				return true
			})
		}
	}
	return nil
}

func enforceForeignKey(tbl *Table, row Row, db map[string]*Table) error {
	for _, fn := range tbl.Constraint {
		if fn.Type == "FOREIGN" {
			refTable, ok := db[fn.RefTable]
			if !ok {
				return nil
			}
			found := false
			refTable.Index.Ascend(func(item btree.Item) bool {
				refRow := item.(Row)
				match := true
				for i, col := range fn.Columns {
					if row.Data[col] != refRow.Data[fn.RefColumns[i]] {
						match = false
						break
					}
				}
				if match {
					found = true
				}
				return true
			})
			if !found {
				log.Error().Msgf("foreign key constraint violation: row does not match referenced table %s", fn.RefTable)

			}
		}
	}
	return nil
}
