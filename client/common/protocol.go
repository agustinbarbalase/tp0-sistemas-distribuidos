package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

// Sizes of the different fields
const (
	SIZE_HEADER_BYTES = 1  // Size of the message header
	SIZE_LENGTH_BYTES = 2  // Size for message length
)

// Header values for messages
const (
	BET_HEADER  = 0x01  // Header indicating a bet message
	OK_HEADER   = 0x02  // Header indicating a success response
	FAIL_HEADER = 0x03  // Header indicating a failure response
)

// Consts for message formatting
const (
	SEPARATOR         = ";"  // Message separator
	NUM_OF_ATTRIBUTES =  6   // Number of attributes in a bet message
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
func htons(value uint16) []byte {
	msg := make([]byte, 2)
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

// serializeBet serializes a Bet struct and an AgencyID into a string using the predefined separator.
// The resulting string contains the following fields in order, separated by SEPARATOR:
//   AgencyID;FirstName;LastName;Document;Birthdate;Number
// This format is used for transmitting bet data over the network.
func (p *Protocol) serializeBet(agencyID string, bet *Bet) string {
	return fmt.Sprintf("%s%s%s%s%s%s%d%s%s%s%d",
		agencyID, SEPARATOR,
		bet.FirstName, SEPARATOR,
		bet.LastName, SEPARATOR,
		bet.Document, SEPARATOR,
		bet.Birthdate, SEPARATOR,
		bet.Number,
	)
}

// SendBet sends a bet message to the server using the protocol format.
// The message consists of:
//   - A single-byte header indicating a bet message.
//   - A 2-byte length of the serialized bet.
//   - The serialized bet data as a string, with fields separated by SEPARATOR.
//
// The serialized bet fields are in the following order:
//   AgencyID;FirstName;LastName;Document;Birthdate;Number
//
// Returns an error if any part of the message fails to send.
func (p *Protocol) SendBet(agencyID string, bet *Bet) error {
	// Send header
	messageHeader := []byte{BET_HEADER}
	if err := p.writeAll(messageHeader, SIZE_HEADER_BYTES); err != nil {
		return err
	}

	betSerialize := p.serializeBet(agencyID, bet)
	betSerializeLength := len(betSerialize)

	// Send length
	messageLength := []byte(htons(uint16(betSerializeLength)))
	if err := p.writeAll(messageLength, SIZE_LENGTH_BYTES); err != nil {
		return err
	}

	// Send bet serialized
	messageBetSerialized := []byte(betSerialize)
	if err := p.writeAll(messageBetSerialized, betSerializeLength); err != nil {
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
