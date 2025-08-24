package common

import (
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	conn      net.Conn
	signalChannel chan os.Signal
	isClosed  bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		signalChannel: make(chan os.Signal, 1),
		isClosed: false,
	}
	signal.Notify(client.signalChannel, syscall.SIGTERM)
	signal.Notify(client.signalChannel, syscall.SIGINT)
	go client.handleSignal()
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// Shutdown the client gracefully
func (c *Client) shutdown() {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
	c.isClosed = true
	if c.conn != nil {
		c.conn.Close()
	}
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
}

// handleSignal listens for termination signals 
// and shuts down the client gracefully
func (c *Client) handleSignal() {
	<-c.signalChannel
	close(c.signalChannel)
	c.shutdown()
	os.Exit(0)
}

// StartClientLoop Send bets
func (c *Client) StartClientLoop() {
	if err := c.createClientSocket(); err != nil {
		return
	}

	protocol := NewProtocol(c.conn)

	document := os.Getenv("DOCUMENTO")
	number, err := strconv.Atoi(os.Getenv("NUMERO"))
	if err != nil {
		if c.isClosed { 
			return 
		}
		log.Errorf("action: parse_number | result: fail | client_id: %v | error: %v", c.config.ID, err)
		number = 0
	}

	err = protocol.SendBet(
		&Bet{
			FirstName:  os.Getenv("NOMBRE"),
			LastName:   os.Getenv("APELLIDO"),
			Document:   document,
			Birthdate:  os.Getenv("NACIMIENTO"),
			Number:     number,
		},
	)

	if err != nil {
		if c.isClosed { 
			return 
		}
		log.Error("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v", document, number, err)
		return
	}

	err = protocol.RecvAckMsg()
	if err != nil {
		if c.isClosed { 
			return 
		}
		log.Error("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v", document, number, err)
		return
	}

	
	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", document, number)
}
