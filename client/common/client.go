package common

import (
	"bufio"
	"net"
	"os"
	"os/signal"
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
	MaxAmount     int
	DataFilePath  string
}

// Client Entity that encapsulates how
type Client struct {
	config        ClientConfig
	conn          net.Conn
	signalChannel chan os.Signal
	isClosed      bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:        config,
		signalChannel: make(chan os.Signal, 1),
		isClosed:      false,
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

// sendABatchBet sends a batch of bets to the server using the established protocol.
// It transmits the batch, waits for an acknowledgment message, and logs the result.
// Returns an error if sending the batch or receiving the acknowledgment fails.
func (c *Client) sendABatchBet(protocol *Protocol, batch *Batch) error {
	if err := protocol.SendBatchBet(c.config.ID, batch); err != nil {
		return err
	}

	status, numOfBets, err := protocol.RecvOKMsg()
	if err != nil {
		return err
	}

	result := "failed"
	if status {
		result = "success"
	}
	log.Infof("action: apuesta_enviada | result: %s | cantidad: %d", result, numOfBets)
	return nil
}

// sendAllBatches reads data from the file specified in the configuration,
// splits the data into batches using a batch iterator, and sends each batch to the server.
func (c *Client) sendAllBatches(p *Protocol) {
	file, err := os.Open(c.config.DataFilePath)
	if err != nil {
		log.Errorf("Failed to open file: %v", err)
		c.conn.Close()
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner == nil {
		log.Errorf("failed to create scanner for file")
		c.conn.Close()
		return
	}

	batchIterator := CreateIteratorByScanner(c.config.ID, scanner, c.config.MaxAmount)
	
	for {
		batch, ok := batchIterator.Next()
		if !ok {
			break
		}
		if err := c.sendABatchBet(p, batch); err != nil {
			if c.isClosed {
				return
			}
			log.Errorf("failed to send batch bet: %v", err)
		}
	}

	if err := p.SendFinishMsg(); err != nil {
		if !c.isClosed {
			c.conn.Close()
			log.Errorf("failed to send finish message: %v", err)
		}
	}
}

// recvWinners receives winner bets from the provided Protocol.
func (c *Client) recvWinners(p *Protocol) []*Bet {
	winners := make([]*Bet, 0)

	for {
		winner, err := p.RecvWinner()
		if err != nil {
			log.Errorf("failed to receive winner: %v", err)
			break
		}
		if winner == nil {
			break
		}
		winners = append(winners, winner)
	}
	
	return winners
}

// StartClient Send bets
func (c *Client) StartClient() {
	if err := c.createClientSocket(); err != nil || c.isClosed {
		return
	}

	protocol := NewProtocol(c.conn)
	if err := protocol.SendIdentification(c.config.ID); err != nil {
		log.Errorf("Failed to send identification: %v", err)
		c.conn.Close()
		return
	}

	c.sendAllBatches(protocol)
	winners := c.recvWinners(protocol)

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))

	time.Sleep(5 * time.Second)

	c.conn.Close()

	log.Info("action: closed_connection | result: success")
}
