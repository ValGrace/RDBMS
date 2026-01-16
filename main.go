package main

import (
	"os"
	"time"

	"github.com/ValGrace/rdbms/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Set time in UNIX format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: false})

	cli.NewPrompt("testing the prompt database.db")

}
