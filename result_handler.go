package mtgox

import (
	"encoding/json"
	"fmt"
)

type ResultHeader struct {
	Id string `json:"id"`
	Op string `json:"op"`
}

func (g *Client) handleResult(data []byte) {
	// Handle Order and other result data here
	fmt.Println("RESULT")

	var payload map[string]interface{}
	json.Unmarshal(data, &payload)
	fmt.Println(string(PrettyPrintJson(payload)))
}

func (g *Client) handleDebug(data []byte) {
	var payload map[string]interface{}
	json.Unmarshal(data, &payload)
	fmt.Printf("DEBUG:\n%s\n", PrettyPrintJson(payload))
}
