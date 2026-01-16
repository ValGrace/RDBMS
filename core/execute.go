package core

import (
	"github.com/rs/zerolog/log"
	"github.com/xwb1989/sqlparser"
	"gitlab.ipsyn.net/bergacorp/oakdb"
)

func ExecuteStatement(stmt sqlparser.Statement) {
	bt := oakdb.NewBTree(3)
	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		log.Info().Msg("Executing INSERT statement")
		bt.Insert(stmt)
	case *sqlparser.Select:
		log.Info().Msg("Executing SELECT statement")
		value, found := bt.Search(stmt)
		if found {
			log.Info().Msgf("Found value: %v", value)
		}
	}
}
