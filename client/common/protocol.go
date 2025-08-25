package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Sizes of the different fields (in bytes)
const (
	SIZE_HEADER_BYTES    = 1  // Size of the message header
	SIZE_FIELD_BYTES     = 2  // Size used to encode variable-length fields (e.g., names)
	SIZE_DOCUMENT_BYTES  = 8  // Size of the document field (DNI)
	SIZE_DATE_BYTES      = 10 // Size of the birthdate field
	SIZE_NUMBER_BYTES    = 4  // Size of the bet number
)

// Header values for messages
const (
	BET_HEADER  = 0x01  // Header indicating a bet message
	OK_HEADER   = 0x02  // Header indicating a success response
	FAIL_HEADER = 0x03  // Header indicating a failure response
)

// Bet represents the data structure of a betting message
type Bet struct {
	FirstName string
	LastName  string
	Document  int
	Birthdate string
	Number    int
}

// Protocol encapsulates the communication mechanism over a socket
type Protocol struct {
	Conn net.Conn
}

// htons converts an integer to a 2-byte slice in big-endian order.
// Used for encoding variable-length fields.
func htons(value uint16) []byte {
	msg := make([]byte, SIZE_FIELD_BYTES)
	binary.BigEndian.PutUint16(msg, value)
	return msg
}

// NewProtocol initializes a new Protocol given a socket connection
func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		Conn: conn,
	}
}

// readAll ensures that exactly `totalLength` bytes are read from the connection.
// Keeps reading until the expected number of bytes is retrieved.
func (p *Protocol) readAll(msg []byte, totalLength int) error {
	readed := 0
	for readed < totalLength {
		n, err := p.Conn.Read(msg[readed:])
		if err != nil {
			return err
		}
		readed += n
	}
	return nil
}

// writeAll ensures that exactly `totalLength` bytes are written to the connection.
// Keeps writing until the expected number of bytes is sent.
func (p *Protocol) writeAll(msg []byte, totalLength int) error {
	writed := 0
	for writed < totalLength {
		n, err := p.Conn.Write(msg[writed:])
		if err != nil {
			return err
		}
		writed += n
	}
	return nil
}

// SendBet sends a bet message to the server following the defined protocol.
//
// Message format:
//   - Header (1 byte)
//   - First name length (2 bytes) + string
//   - Last name length (2 bytes) + string
//   - Document (8 bytes)
//   - Birthdate (10 bytes)
//   - Number (4 bytes, zero-padded string)
func (p *Protocol) SendBet(bet *Bet) error {
	// Header
	header := []byte{BET_HEADER}
	if err := p.writeAll(header, SIZE_HEADER_BYTES); err != nil {
		return err
	}

	// First name
	sizeFirstName := htons(uint16(len(bet.FirstName)))
	if err := p.writeAll(sizeFirstName, SIZE_FIELD_BYTES); err != nil {
		return err
	}
	if err := p.writeAll([]byte(bet.FirstName), len(bet.FirstName)); err != nil {
		return err
	}

	// Last name
	sizeLastName := htons(uint16(len(bet.LastName)))
	if err := p.writeAll(sizeLastName, SIZE_FIELD_BYTES); err != nil {
		return err
	}
	if err := p.writeAll([]byte(bet.LastName), len(bet.LastName)); err != nil {
		return err
	}

	// Document
	if err := p.writeAll([]byte(bet.Document), SIZE_DOCUMENT_BYTES); err != nil {
		return err
	}

	// Birthdate
	if err := p.writeAll([]byte(bet.Birthdate), SIZE_DATE_BYTES); err != nil {
		return err
	}

	// Number (zero-padded to 4 chars)
	numberStr := fmt.Sprintf("%04d", bet.Number)
	if err := p.writeAll([]byte(numberStr), SIZE_NUMBER_BYTES); err != nil {
		return err
	}

	return nil
}

// RecvOKMsg waits for an acknowledgment message from the server.
// The server should reply with a single-byte header = OK_HEADER (0x02).
// Returns an error if the message is invalid or not received.
func (p *Protocol) RecvOKMsg() error {
	ack := make([]byte, SIZE_HEADER_BYTES)
	if err := p.readAll(ack, SIZE_HEADER_BYTES); err != nil {
		return err
	}

	if ack[0] != OK_HEADER {
		return fmt.Errorf("invalid OK message")
	}
	return nil
}
