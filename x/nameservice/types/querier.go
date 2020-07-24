package types

import "strings"

// query endpoints supported by the nameservice Querier
const (
	QueryResolve = "resolve"
	QueryWhois   = "whois"
	QueryNames   = "names"
)

// query results must implement fmt.Stringer interface, and JSON marshalable

// QueryResResolve Queries Result Payload for a resolve query
type QueryResResolve struct {
	Value string `json:"value"`
}

// implement fmt.Stringer
func (r QueryResResolve) String() string {
	return r.Value
}

// QueryResNames Queries Result Payload for a names query
type QueryResNames []string

// implement fmt.Stringer
func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}
