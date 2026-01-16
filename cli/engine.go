package cli

import (
	"os"
	"strings"

	"github.com/ValGrace/rdbms/compiler"
	"github.com/ValGrace/rdbms/core"
	prompt "github.com/c-bata/go-prompt"
	"github.com/rs/zerolog/log"
)

func completer(d prompt.Document) []prompt.Suggest {
	// prompt suggestions
	s := []prompt.Suggest{}
	for _, sug := range suggestionsMap {
		s = append(s, sug...)
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func getExecutor(string) func(string) {
	// TODO: setup file for database
	return func(prStr string) {
		prStr = strings.TrimSpace(prStr) // remove the trailing space/newline
		prStr = strings.ToLower(prStr)

		switch prStr {
		case "":
			return
		case "exit", "quit":
			log.Info().Msg("Exiting PesaPal RDBMS prompt. Goodbye! 👋")
			os.Exit(0)
		default:
			if strings.HasPrefix(prStr, ".") {
				log.Error().Msgf("Unknown command: %s", prStr)
				break
			}
			// TODO: Prepare statement with sql compiler
			stmt, res := compiler.PrepareStatement(prStr)
			if res == compiler.SUCCESS {
				core.ExecuteStatement(stmt)
			} else {
				log.Error().Msgf("Failed to execute statement: %s", res)
			}
		}
	}

}

func NewPrompt(file string) {
	p := prompt.New(
		getExecutor(file),
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("Welcome to PesaPal RDBMS💸💳"),
		prompt.OptionInputTextColor(prompt.DarkRed),
		prompt.OptionPrefixBackgroundColor(prompt.Black),
	)
	p.Run()
}
