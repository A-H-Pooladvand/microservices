package Filter

import "strings"

// Filter represents the filters that can be applied to a query.
type Filter struct {
	Select []string `json:"select"`
	Search string   `json:"search"`
}

// Selects returns the fields that should be selected.
func (f Filter) Selects() string {
	return strings.Join(f.Select, ", ")
}
