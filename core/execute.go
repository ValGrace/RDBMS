package core

import (
	"github.com/ValGrace/rdbms/compiler"
	"github.com/rs/zerolog/log"
)

func ExecuteStatement(stmt compiler.Statement) {
	switch stmt.Type {
	case compiler.Insert:
		log.Info().Msg("Executing INSERT statement")
	case compiler.Select:
		log.Info().Msg("Executing SELECT statement")
	default:
		log.Error().Msg("Cannot execute unknown statement type")
	}
}
