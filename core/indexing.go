package core

import (
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
