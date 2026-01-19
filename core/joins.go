package core

import (
	"github.com/google/btree"
	"github.com/xwb1989/sqlparser"
)

func executeJoin(left *Table, right *Table, on *sqlparser.ComparisonExpr) []map[string]string {
	results := []map[string]string{}

	leftCol := on.Left.(*sqlparser.ColName).Name.String()
	rightCol := on.Right.(*sqlparser.ColName).Name.String()

	left.Index.Ascend(func(item btree.Item) bool {
		lrow := item.(Row)

		right.Index.Ascend(func(ritem btree.Item) bool {
			rrow := ritem.(Row)

			if lrow.Data[leftCol] == rrow.Data[rightCol] {
				combined := map[string]string{}
				for k, v := range lrow.Data {
					combined[left.Name+"."+k] = v
				}
				for k, v := range rrow.Data {
					combined[right.Name+"."+k] = v
				}
				results = append(results, combined)
			}
			return true
		})
		return true
	})
	return results
}
