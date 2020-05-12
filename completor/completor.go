package completor

import (
	"github.com/c-bata/go-prompt"
)

// Completer .
type Completer struct {
}

// NewCompleter .
func NewCompleter(version string) (*Completer, error) {
	return &Completer{}, nil
}

// Complete .
func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {

	if d.TextBeforeCursor() == "" {
		return prompt.FilterHasPrefix(commands, d.GetWordBeforeCursor(), true)
	}

}
