package cli

import setprompt "github.com/c-bata/go-prompt"

type suggestionType int

const (
	Commands suggestionType = iota
)

var commandSuggestions = []setprompt.Suggest{
	{Text: "SELECT", Description: "Select data from a table"},
	{Text: "INSERT", Description: "Insert data into a table"},
	{Text: "UPDATE", Description: "Update data in a table"},
	{Text: "DELETE", Description: "Delete data from a table"},
	{Text: "CREATE", Description: "Create a new table or database"},
	{Text: "DROP", Description: "Drop a table or database"},
	{Text: "JOIN", Description: "Join tables"},
	{Text: "exit", Description: "Quit/Exit the prompt"},
	{Text: "quit", Description: "Quit/Exit the prompt"},
}

var suggestionsMap = map[suggestionType][]setprompt.Suggest{
	Commands: commandSuggestions,
}
