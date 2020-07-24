package types

import (
	"fmt"
)

// GenesisState - all nameservice state that must be provided at genesis
type GenesisState struct {
	WhoisRecords []Whois `json:"whois_records"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(whoIsRecords []Whois) GenesisState {
	return GenesisState{
		WhoisRecords: whoIsRecords,
	}
}

// DefaultGenesisState - default GenesisState provides minimal state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords: []Whois{},
	}
}

// ValidateGenesis validates the nameservice genesis parameters
func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WhoisRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
		}
		if record.Value == "" {
			return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
		}
		if record.Price == nil {
			return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
		}
	}
	return nil
}
