package mtgox

/*
  package mtgox provides a streaming implementation of Mt. Client's bitcoin trading API
  built on the Gorilla Websocket library
*/

import (
	"crypto/hmac"
	"crypto/sha512"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	secureAPIHost string = "wss://websocket.mtgox.com:443"
	apiHost       string = "ws://websocket.mtgox.com:80"
	apiPath       string = "/mtgox"
	httpEndpoint  string = "http://mtgox.com/api/2"
	originURL     string = "http://websocket.mtgox.com"

	// TODO: Move this into Config
	secureConn bool = true

	// BitcoinDivision represents the current integer division of 1 bitcoin.
	BitcoinDivision = 1e8
)

var (
	currencyDivisions = map[string]float64{
		"BTC": BitcoinDivision,
		"USD": 1e5,
		"GBP": 1e5,
		"EUR": 1e5,
		"JPY": 1e3,
		"AUD": 1e5,
		"CAD": 1e5,
		"CHF": 1e5,
		"CNY": 1e5,
		"DKK": 1e5,
		"HKD": 1e5,
		"PLN": 1e5,
		"RUB": 1e5,
		"SEK": 1e3,
		"SGD": 1e5,
		"THB": 1e5,
	}
)

// ErrorHandlerFunc is a function type to use as an error callback.
type ErrorHandlerFunc func(error)

// Client represents the public type for interacing with the MtGox streaming API.
type Client struct {
	key    []byte
	secret []byte
	conn   *websocket.Conn

	Ticker chan *TickerPayload
	Info   chan *Info
	Depth  chan *DepthPayload
	Trades chan *TradePayload
	Orders chan []Order
	errors chan error
	done   chan bool

	errHandler ErrorHandlerFunc

	requestListeners map[string]chan []byte
}

// Config represents a configuration type to be used when configuring the Client.
type Config struct {
	Currencies []string
	Key        string
	Secret     string
	SecureConn bool
}

type payload struct {
	messageType int
	data        []byte
}

// StreamHeader represents the header of a payload message from MtGox before
// being parsed.
type StreamHeader struct {
	Channel     string `json:"channel"`
	ChannelName string `json:"channel_name"`
	Op          string `json:"op"`
	Origin      string `json:"origin"`
	Private     string `json:"private"`
}

// New constructs a new instance of Client, returning an error
func New(key, secret string, currencies ...string) (*Client, error) {
	var streamURL string
	if secureConn {
		streamURL = fmt.Sprintf("%s%s?Currency=%s", secureAPIHost, apiPath, strings.Join(currencies, ","))
	} else {
		streamURL = fmt.Sprintf("%s%s?Currency=%s", apiHost, apiPath, strings.Join(currencies, ","))
	}

	u, err := url.Parse(streamURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %s", err.Error())
	}

	var netConn net.Conn

	if secureConn {
		netConn, err = tls.Dial("tcp", u.Host, nil)
	} else {
		netConn, err = net.Dial("tcp", u.Host)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting: %s", err.Error())
	}

	conn, _, err := websocket.NewClient(netConn, u, http.Header{"Origin": {originURL}}, 1024, 1024)
	if err != nil {
		return nil, fmt.Errorf("opening websocket: %v", err)
	}

	return NewWithConnection(key, secret, conn)
}

// NewWithConnection constructs a new client using an existing connection, useful for testing
func NewWithConnection(key, secret string, conn *websocket.Conn) (g *Client, err error) {
	g = &Client{
		conn:             conn,
		Ticker:           make(chan *TickerPayload, 1),
		Info:             make(chan *Info, 1),
		Depth:            make(chan *DepthPayload, 1),
		Trades:           make(chan *TradePayload, 1),
		Orders:           make(chan []Order, 1),
		errors:           make(chan error, 1),
		done:             make(chan bool, 1),
		requestListeners: make(map[string]chan []byte),
	}

	g.key, err = hex.DecodeString(strings.Replace(key, "-", "", -1))
	if err != nil {
		return nil, err
	}

	g.secret, err = base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (c *Client) SetErrorHandler(handlerFunc ErrorHandlerFunc) {
	c.errHandler = handlerFunc
}

// Start begins the internal routines for processing messages and errors
func (g *Client) Start() {
	// Handle incoming messages
	go func() {
		for p := range g.messages() {
			g.handle(p)
		}
	}()

	g.startErrorHandler()
}

func (g *Client) startErrorHandler() {
	// Handle errors
	go func() {
		for err := range g.errors {
			if g.errHandler != nil {
				g.errHandler(err)
			}
		}
	}()
}

// Close disconnects the client from the streaming API.
func (g *Client) Close() {
	g.done <- true
}

// Conn returns the raw websocket connection to Mt.Gox
// (This may be removed in the future)
func (g *Client) Conn() *websocket.Conn {
	return g.conn
}

// Reads messages into a channel so we can select on them later
func (g *Client) messages() <-chan []byte {
	msgs := make(chan []byte, 10)

	go func(msgs chan []byte) {
		for {
			messageType, data, err := g.conn.ReadMessage()
			if err != nil {
				g.errors <- err
				break
			}

			if messageType == websocket.TextMessage {
				msgs <- data
			} else {
				g.errors <- fmt.Errorf("received unknown message type: %d", messageType)
			}
		}
	}(msgs)

	return msgs
}

func (g *Client) sign(body []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, g.secret)
	_, err := mac.Write(body)
	if err != nil {
		return nil, err
	}

	return mac.Sum(nil), nil
}

func (g *Client) authenticatedSend(msg map[string]interface{}) error {
	if g.key == nil || g.secret == nil {
		return errors.New("key or secret is invalid or missing")
	}

	req, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	signedReq, err := g.sign(req)
	if err != nil {
		return err
	}

	requestID := msg["id"]

	fullReq := append(append(g.key, signedReq...), req...)
	encodedReq := base64.StdEncoding.EncodeToString(fullReq)

	reqBody := map[string]interface{}{
		"op":      "call",
		"id":      requestID,
		"call":    encodedReq,
		"context": "mtgox.com",
	}

	reqJSON, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	return g.conn.WriteMessage(websocket.TextMessage, reqJSON)
}

// Handler function for processing responses from mtgox
func (g *Client) handle(data []byte) {
	var header StreamHeader
	json.Unmarshal(data, &header)

	switch header.Private {
	case "debug":
		g.handleDebug(data)

	case "ticker":
		g.handleTicker(data)

	case "trade":
		g.handleTrade(data)

	case "depth":
		g.handleDepth(data)

	default:
		if header.Op == "result" {
			g.handleResult(data)
		} else {
			fmt.Printf("HANDLE: %v\n", header.Op)

			var payload map[string]interface{}
			json.Unmarshal(data, &payload)
			fmt.Println(string(prettyPrintJSON(payload)))
		}
	}
}

func prettyPrintJSON(p interface{}) []byte {
	formattedJSON, err := json.MarshalIndent(&p, "", "  ")
	if err != nil {
		return []byte("{}")
	}
	return formattedJSON
}

func (g *Client) call(endpoint string, params map[string]interface{}) (string, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	id := <-ids

	msg := map[string]interface{}{
		"call":   endpoint,
		"item":   "BTC",
		"params": params,
		"id":     id,
		"nonce":  <-nonces,
	}

	return id, g.authenticatedSend(msg)
}

func (c *Client) enqueuePendingRequest(id string, ch chan []byte) {
	c.requestListeners[id] = ch
}
