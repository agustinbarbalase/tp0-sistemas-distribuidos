package common

import (
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"
	"bufio"
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
		log.Errorf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner == nil {
		log.Errorf("failed to create scanner for file")
		return
	}
	
	batch := NewBatch(c.config.MaxAmount, 8 * 1024)

	for scanner.Scan() {
		line := scanner.Text()
		bet, err := ProcessCSVLine(line)
		if err != nil {
			log.Errorf("failed to process CSV line: %v", err)
			continue
		}

		if !batch.AddBet(c.config.ID, bet) {
			if err := protocol.SendBatchBet(c.config.ID, batch); err != nil {
				log.Errorf("failed to send batch bet: %v", err)
				continue
			}

			status, numOfBets, err := protocol.RecvOKMsg()
			if err != nil {
				log.Errorf("failed to receive OK message: %v", err)
			}

			result := "failed"
			if status {
				result = "success"
			}
			log.Infof("action: apuesta_recibida | result: %s | cantidad: %d", result, numOfBets)

			batch = NewBatch(c.config.MaxAmount, 8 * 1024)
			batch.AddBet(c.config.ID, bet)
		}
	}

	if batch.Amount > 0 {
		if err := protocol.SendBatchBet(c.config.ID, batch); err != nil {
			log.Errorf("failed to send batch bet: %v", err)
		}

		status, numOfBets, err := protocol.RecvOKMsg()
		if err != nil {
			log.Errorf("failed to receive OK message: %v", err)
		}

		result := "failed"
		if status {
			result = "success"
		}
		log.Infof("action: apuesta_recibida | result: %s | cantidad: %d", result, numOfBets)
	}

	protocol.SendFinishMsg()

	log.Infof("action: apuesta_recibida | result: success")
	c.conn.Close()
}
