package common

import (
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxAmount 	  int
	DataFilePath  string
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	conn      net.Conn
	signalChannel chan os.Signal
	isClosed  bool
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
	if err := c.createClientSocket(); err != nil || c.isClosed {
		return
	}
	
	protocol := NewProtocol(c.conn)

	file, err := os.Open(c.config.DataFilePath)
	if err != nil {
		log.Error("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	if err := protocol.SendBatchBet(c.config.ID, c.config.MaxAmount, file); err != nil {
		if !c.isClosed {
			log.Error("action: crear_batch | result: fail | error: %v", err)
			c.conn.Close()
		}
		return
	}
	
	if err := protocol.RecvOKMsg(); err != nil {
		if !c.isClosed {
			log.Error("action: apuesta_recibida | result: fail | error: %v", err)
			c.conn.Close()
		}
		return
	}

	log.Infof("action: apuesta_recibida | result: success")
	c.conn.Close()
}
