package mtgox

import (
	"testing"
	"time"
)

var lagPayload = []byte(`
{
  "id": "5e399504b8918c295ca54ac66d28785ccbe16e21",
  "op": "result",
  "result": {
    "lag": 2.26162e+06,
    "lag_secs": 2.26162,
    "lag_text": "2.26162 seconds",
    "length": "10"
  }
}`)

func TestProcessLagResult(t *testing.T) {
	client := newTestClient(t)

	lag, err := client.processLagResult(lagPayload)
	if err != nil {
		t.Error(err.Error())
	}

	if expected := time.Microsecond * 2.26162e+06; lag != expected {
		t.Errorf("Expected lag to = %v, got %v", expected, lag)
	}
}
