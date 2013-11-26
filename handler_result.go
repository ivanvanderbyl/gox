package mtgox

import (
	"encoding/json"
	"fmt"
	"time"
)

type ResultHeader struct {
	Id string `json:"id"`
	Op string `json:"op"`
}

type ResultPayload struct {
	ResultHeader
	Result json.RawMessage `json:"result"`
}

func (g *Client) handleResult(data []byte) {
	var header ResultHeader
	json.Unmarshal(data, &header)

	if ch, ok := g.requestListeners[header.Id]; ok {
		ch <- data
	} else {
		// Handle Order and other result data here
		fmt.Printf("RESULT: %s\n", header.Id)

		var payload map[string]interface{}
		json.Unmarshal(data, &payload)
		fmt.Println(string(PrettyPrintJson(payload)))
	}
}

func (g *Client) handleDebug(data []byte) {
	var payload map[string]interface{}
	json.Unmarshal(data, &payload)
	fmt.Printf("DEBUG:\n%s\n", PrettyPrintJson(payload))
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
