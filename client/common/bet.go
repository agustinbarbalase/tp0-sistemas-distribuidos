package common

import (
	"fmt"
	"strconv"
	"strings"
)

const SEPARATOR = ";"  // Message separator

// Bet represents the data structure of a betting message
type Bet struct {
	FirstName string
	LastName  string
	Document  int
	Birthdate string
	Number    int
}

// NewBet creates and returns a new Bet instance.
func NewBet(firstName, lastName string, document int, birthdate string, number int) *Bet {
	return &Bet{
		FirstName: firstName,
		LastName:  lastName,
		Document:  document,
		Birthdate: birthdate,
		Number:    number,
	}
}

// serializeBet serializes a Bet struct and an AgencyID into a string using the predefined separator.
// The resulting string contains the following fields in order, separated by SEPARATOR:
//   AgencyID;FirstName;LastName;Document;Birthdate;Number
// This format is used for transmitting bet data over the network.
func (b *Bet) serializeBet(agencyID string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%d%s%s%s%d",
		agencyID, SEPARATOR,
		b.FirstName, SEPARATOR,
		b.LastName, SEPARATOR,
		b.Document, SEPARATOR,
		b.Birthdate, SEPARATOR,
		b.Number,
	)
}

// SerializedBetLength returns the length of the serialized bet string for a given AgencyID.
func (b *Bet) SerializedBetLength(agencyID string) int {
	return len(b.serializeBet(agencyID))
}

// ProcessCSVLine parses a CSV line string and returns a Bet.
// The expected order is: FirstName,LastName,Document,Birthdate,Number.
func ProcessCSVLine(line string) (*Bet, error) {
	fields := strings.Split(line, ",")
	
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid number of fields: expected 5, got %d", len(fields))
	}
	
	document, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, fmt.Errorf("invalid document: %v", err)
	}
	
	number, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, fmt.Errorf("invalid number: %v", err)
	}

	return &Bet{
		FirstName: fields[0],
		LastName:  fields[1],
		Document:  document,
		Birthdate: fields[3],
		Number:    number,
	}, nil
}
