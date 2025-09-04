package common

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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
	config        ClientConfig
	protocol      *Protocol
	signalChannel chan os.Signal
	isClosed      bool
}

const dateFormat = "2006-01-02"

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		signalChannel: make(chan os.Signal, 1),
		isClosed: false,
	}
	signal.Notify(client.signalChannel, syscall.SIGTERM)
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
	c.protocol = NewProtocol(conn)
	return nil
}

// Shutdown the client gracefully
func (c *Client) shutdown() {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
	c.isClosed = true
	if c.protocol != nil {
		c.protocol.Close()
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

// createBet constructs a Bet instance using environment variables for its fields.
// It retrieves and parses the following environment variables:
//   - NUMERO: the bet number (must be an integer)
//   - DOCUMENTO: the user's document number (must be an integer)
//   - NOMBRE: the user's first name
//   - APELLIDO: the user's last name
//   - NACIMIENTO: the user's birthdate
// If NUMERO or DOCUMENTO cannot be parsed as integers, it logs an error and returns nil with the parsing error.
// Returns a pointer to the created Bet and an error if any occurred during parsing.
func (c *Client) createBet() (*Bet, error) {
	number, err := strconv.Atoi(os.Getenv("CLI_NUMERO"))
	if err != nil {
		log.Error("action: numero_invalido | numero: %s | error: %s", os.Getenv("CLI_NUMERO"), err)
		return nil, err
	}

	document, err := strconv.Atoi(os.Getenv("CLI_DOCUMENTO"))
	if err != nil {
		log.Error("action: documento_invalido | documento: %s | error: %s", os.Getenv("CLI_DOCUMENTO"), err)
		return nil, err
	}

	birthdate := os.Getenv("CLI_NACIMIENTO")
	if _, err := time.Parse(dateFormat, birthdate); err != nil {
		log.Error("action: nacimiento_invalido | nacimiento: %s | error: %s", birthdate, err)
		return nil, err
	}

	return &Bet{
		FirstName: os.Getenv("CLI_NOMBRE"),
		LastName:  os.Getenv("CLI_APELLIDO"),
		Document:  document,
		Birthdate: birthdate,
		Number:    number,
	}, nil
}

// StartClient Send bets
func (c *Client) StartClient() {
	if err := c.createClientSocket(); err != nil || c.isClosed {
		return
	}
	
	bet, err := c.createBet()
	if err != nil {
		if !c.isClosed {
			c.protocol.Close()
		}
		return
	}	

	if err := c.protocol.SendBet(c.config.ID, bet); err != nil {
		if !c.isClosed {
			log.Error("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v", bet.Document, bet.Number, err)
			c.protocol.Close()
		}
		return
	}
	
	if err := c.protocol.RecvOKMsg(); err != nil {
		if !c.isClosed {
			log.Error("action: apuesta_enviada | result: fail | dni: %v | numero: %v | error: %v", bet.Document, bet.Number, err)
			c.protocol.Close()
		}
		return
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.Document, bet.Number)
	c.protocol.Close()
	log.Info("action: closed_connection | result: success")
}
