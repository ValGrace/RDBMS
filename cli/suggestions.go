package cli

import setprompt "github.com/c-bata/go-prompt"

type suggestionType int

const (
	Commands suggestionType = iota
	Keywords
)

var commandSuggestions = []setprompt.Suggest{

	{Text: "UPDATE", Description: "Update data in a table"},
	{Text: "DELETE", Description: "Delete data from a table"},
	{Text: "CREATE", Description: "Create a new table or database"},
	{Text: "DROP", Description: "Drop a table or database"},
	{Text: "JOIN", Description: "Join tables"},
	{Text: "exit", Description: "Quit/Exit the prompt"},
	{Text: "quit", Description: "Quit/Exit the prompt"},
}

var keywordSuggestions = []setprompt.Suggest{
	{Text: "select", Description: "read data from a table"},
	{Text: "insert", Description: "add data to a table"},
}

var suggestionsMap = map[suggestionType][]setprompt.Suggest{
	Commands: commandSuggestions,
	Keywords: keywordSuggestions,
}
