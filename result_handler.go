package gox

import (
	"encoding/json"
	"fmt"
)

func (g *Gox) handleResult(data []byte) {
	// Handle Order and other result data here
	fmt.Println("RESULT")

	var payload map[string]interface{}
	json.Unmarshal(data, &payload)
	fmt.Println(string(PrettyPrintJson(payload)))
}

func (g *Gox) handleDebug(data []byte) {
	var payload map[string]interface{}
	json.Unmarshal(data, &payload)
	fmt.Printf("DEBUG:\n%s\n", PrettyPrintJson(payload))
}
