package common

const (
	SEPARATOR         = ";"  // Message separator
)

// Bet represents the data structure of a betting message
type Bet struct {
	FirstName string
	LastName  string
	Document  int
	Birthdate string
	Number    int
}

// serializeBet serializes a Bet struct and an AgencyID into a string using the predefined separator.
// The resulting string contains the following fields in order, separated by SEPARATOR:
//   AgencyID;FirstName;LastName;Document;Birthdate;Number
// This format is used for transmitting bet data over the network.
func (p *bet) serializeBet(agencyID string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%d%s%s%s%d",
		agencyID, SEPARATOR,
		bet.FirstName, SEPARATOR,
		bet.LastName, SEPARATOR,
		bet.Document, SEPARATOR,
		bet.Birthdate, SEPARATOR,
		bet.Number,
	)
}
