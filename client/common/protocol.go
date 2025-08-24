package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Bet struct {
	FirstName string
	LastName  string
	Document  string
	Birthdate string
	Number    int
}

type Protocol struct {
	Conn net.Conn
}

func htons(value int) []byte {
	msg := make([]byte, 2)
	binary.BigEndian.PutUint16(msg, uint16(value))
	return msg
}

func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{
		Conn: conn,
	}
}

func (p *Protocol) readAll(msg []byte, total_length int) (error) {
	readed := 0

	for readed < total_length {
		n, err := p.Conn.Read(msg[readed:])
		if err != nil {
			return err
		}
		readed += n
	}

	return nil
}

func (p *Protocol) writeAll(msg []byte, total_length int) error {
	writed := 0

	for writed < total_length {
		n, err := p.Conn.Write(msg[writed:])
		if err != nil {
			return err
		}
		writed += n
	}

	return nil
}

func (p *Protocol) SendBet(bet *Bet) error {
	header := []byte{0x01}
	if err := p.writeAll(header, len(header)); err != nil {
		return err
	}

	sizeFirstName := htons(len(bet.FirstName))
	if err := p.writeAll(sizeFirstName, len(sizeFirstName)); err != nil {
		return err
	}
	if err := p.writeAll([]byte(bet.FirstName), len(bet.FirstName)); err != nil {
		return err
	}

	sizeLastName := htons(len(bet.LastName))
	if err := p.writeAll(sizeLastName, len(sizeLastName)); err != nil {
		return err
	}
	if err := p.writeAll([]byte(bet.LastName), len(bet.LastName)); err != nil {
		return err
	}

	if err := p.writeAll([]byte(bet.Document), 8); err != nil {
		return err
	}

	if err := p.writeAll([]byte(bet.Birthdate), 10); err != nil {
		return err
	}

	numberStr := fmt.Sprintf("%04d", bet.Number)
	if err := p.writeAll([]byte(numberStr), 4); err != nil {
		return err
	}

	return nil
}

func (p *Protocol) RecvAckMsg() error {
	ack := make([]byte, 1)
	if err := p.readAll(ack, 1); err != nil {
		return err
	}

	if ack[0] != 0x02 {
		return fmt.Errorf("invalid ACK message")
	}

	return nil
}
