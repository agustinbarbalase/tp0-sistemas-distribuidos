package common

import (
	"bufio"
	"fmt"
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
	go client.handleSignal(client.signalChannel)
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

// handleSignal listens for termination signals 
// and shuts down the client gracefully
func (c *Client) handleSignal(sigs chan os.Signal) {
	<-sigs
	c.Shutdown()
	os.Exit(0)
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount && c.isClosed; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()
		if c.isClosed {
			return
		}

		// TODO: Modify the send to avoid short-write
		_, err := fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		if err != nil && c.isClosed {
			if !c.isClosed {
				log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)
			}
			return
		}

		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if c.isClosed {
			return
		}
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

// Shutdown the client gracefully
func (c *Client) Shutdown() {
	log.Infof("action: shutdown | result: in_progress | client_id: %v", c.config.ID)
	c.isClosed = true
	if c.conn != nil {
		c.conn.Close()
	}
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
}
