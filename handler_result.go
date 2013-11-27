package mtgox

import (
	"encoding/json"
	"fmt"
	"time"
)

// ResultHeader defines a simple header structure for partially unmarshalling
// replies from Mt.Gox to determine their type
type ResultHeader struct {
	ID string `json:"id"`
	Op string `json:"op"`
}

// ResultPayload defines a structure for unmarshalling complete replies with
// unknown response data formats to be parsed later
type ResultPayload struct {
	ResultHeader
	Result json.RawMessage `json:"result"`
}

func (c *Client) handleResult(data []byte) {
	var header ResultHeader
	json.Unmarshal(data, &header)

	if ch, ok := c.requestListeners[header.ID]; ok {
		ch <- data
	} else {
		// Handle Order and other result data here
		fmt.Printf("RESULT: %s\n", header.ID)

		var payload map[string]interface{}
		json.Unmarshal(data, &payload)
		fmt.Println(string(prettyPrintJSON(payload)))
	}
}

func (c *Client) handleDebug(data []byte) {
	var payload map[string]interface{}
	json.Unmarshal(data, &payload)
	fmt.Printf("DEBUG:\n%s\n", prettyPrintJSON(payload))
}

func (c *Client) processOrderResult(data []byte) ([]Order, error) {
	var p ResultPayload
	err := json.Unmarshal(data, &p)
	if err != nil {
		return []Order{}, err
	}

	var orders []Order
	err = json.Unmarshal(p.Result, &orders)
	if err != nil {
		return []Order{}, err
	}

	return orders, nil
}

type lagResult struct {
	Lag float64 `json:"lag"`
}

func (c *Client) processLagResult(data []byte) (time.Duration, error) {
	var p ResultPayload
	var lag lagResult
	err := json.Unmarshal(data, &p)
	if err != nil {
		return time.Duration(0), nil
	}

	err = json.Unmarshal(p.Result, &lag)
	if err != nil {
		return time.Duration(0), nil
	}

	return time.Duration(lag.Lag) * time.Microsecond, nil
}
