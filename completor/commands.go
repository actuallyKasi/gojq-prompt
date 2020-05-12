package completor

import (
	"github.com/c-bata/go-prompt"
)

var commands = []prompt.Suggest{
	{Text: ".", Description: "Help about any command"},
	{Text: "pr", Description: "Create, view, and checkout pull requests"},
	{Text: "repo", Description: "Create, clone, fork, and view repositories"},
	{Text: "issue", Description: "Create and view issues"},
	// Custom commands.
	{Text: "exit", Description: "Exit this program"},
}
