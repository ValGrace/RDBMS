package core

import (
	"github.com/rs/zerolog/log"
	"github.com/xwb1989/sqlparser"
)

func ExecuteStatement(stmt sqlparser.Statement) {
	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		log.Info().Msg("Executing INSERT statement")
		_ = stmt
	case *sqlparser.Select:
		log.Info().Msg("Executing SELECT statement")
		_ = stmt
	}
}
