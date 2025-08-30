package common

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
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
	BATCH_HEADER = 0x04 // Header indicating a batch of bets
	FINISH_HEADER = 0x05 // Header indicating the end of transmission
	ID_HEADER = 0x06 // Header indicating identification message
)

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

// ntohs converts a 2-byte slice in network byte order (big-endian) to a uint16 in host byte order.
func ntohs(value []byte) uint16 {
	return binary.BigEndian.Uint16(value)
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
	// messageHeader := []byte{BET_HEADER}
	// if err := p.writeAll(messageHeader, SIZE_HEADER_BYTES); err != nil {
	// 	return err
	// }

	betSerialize := bet.serializeBet(agencyID)
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

// SendBatchBet reads bets from the provided reader and sends them one by one using SendBet.
// Each line in the reader should represent a bet in the expected format.
// Returns an error if any bet fails to send.
func (p *Protocol) SendBatchBet(ID string, batch *Batch) error {
	messageHeader := []byte{BATCH_HEADER}
	if err := p.writeAll(messageHeader, SIZE_HEADER_BYTES); err != nil {
		return err
	}

	// Send length
	batchLength := len(batch.Bets)
	messageLength := []byte(htons(uint16(batchLength)))
	if err := p.writeAll(messageLength, SIZE_LENGTH_BYTES); err != nil {
		return err
	}

	for i := range batch.Bets {
		if err := p.SendBet(ID, batch.Bets[i]); err != nil {
			return err
		}
	}

	return nil
}

// RecvOKMsg waits for an acknowledgment message from the server.
// The server should reply with a single-byte header = OK_HEADER (0x02).
// Returns an error if the message is invalid or not received.
func (p *Protocol) RecvOKMsg() (bool, int, error) {
	ack := make([]byte, SIZE_HEADER_BYTES)
	if err := p.readAll(ack, SIZE_HEADER_BYTES); err != nil {
		return false, 0, err
	}

	numOfBets := make([]byte, SIZE_LENGTH_BYTES)
	if err := p.readAll(numOfBets, SIZE_LENGTH_BYTES); err != nil {
		return false, 0, err
	}

	if ack[0] != OK_HEADER {
		return false, 0, fmt.Errorf("invalid header")
	}

	return ack[0] == OK_HEADER, int(ntohs(numOfBets)), nil
}


// SendFinishMsg sends a finish message to the remote endpoint using the protocol.
// It constructs a message header with the FINISH_HEADER byte and writes it to the connection.
// If the message fails to send, an error is logged.
func (p *Protocol) SendFinishMsg() {
	messageHeader := []byte{FINISH_HEADER}
	if err := p.writeAll(messageHeader, SIZE_HEADER_BYTES); err != nil {
		log.Error("failed to send finish message: %v", err)
	}
}


func (p *Protocol) SendIdentification(agencyID string) error {
	messageHeader := []byte{ID_HEADER}
	if err := p.writeAll(messageHeader, SIZE_HEADER_BYTES); err != nil {
		return err
	}

	agencyIDInt, err := strconv.Atoi(agencyID)
	if err != nil {
		return err
	}

	agencyIDBytes := []byte(htons(uint16(agencyIDInt)))
	if err := p.writeAll(agencyIDBytes, SIZE_LENGTH_BYTES); err != nil {
		return err
	}

	return nil
}

	
func (p *Protocol) RecvWinner() (*Bet, error) {
	header := make([]byte, SIZE_HEADER_BYTES)
	if err := p.readAll(header, SIZE_HEADER_BYTES); err != nil {
		return nil, err
	}

	if header[0] == FINISH_HEADER {
		return nil, nil
	} else if header[0] != BET_HEADER {
		return nil, fmt.Errorf("invalid header")
	}

	lengthBytes := make([]byte, SIZE_LENGTH_BYTES)
	if err := p.readAll(lengthBytes, SIZE_LENGTH_BYTES); err != nil {
		return nil, err
	}

	length := ntohs(lengthBytes)

	winnerData := make([]byte, length)
	if err := p.readAll(winnerData, int(length)); err != nil {
		return nil, err
	}

	winner := &Bet{}
	if err := winner.Deserialize(winnerData); err != nil {
		return nil, err
	}

	return winner, nil
}
