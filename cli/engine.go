package cli

import (
	"os"
	"strings"

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
			//TODO: process commands here
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
