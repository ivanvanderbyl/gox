package mtgox

import (
	"testing"
	"time"
)

var orderPayload = []byte(`{
  "id": "aed43fbd20929d4288c17a2c729b9f060af45606",
  "op": "result",
  "result": [
    {
      "actions": [],
      "amount": {
        "currency": "BTC",
        "display": "0.01000000 BTC",
        "display_short": "0.01 BTC",
        "value": "0.01000000",
        "value_int": "1000000"
      },
      "currency": "AUD",
      "date": 1.385257384e+09,
      "effective_amount": {
        "currency": "BTC",
        "display": "0.00998098 BTC",
        "display_short": "0.01 BTC",
        "value": "0.00998098",
        "value_int": "998098"
      },
      "invalid_amount": {
        "currency": "BTC",
        "display": "0.00001902 BTC",
        "display_short": "0.00 BTC",
        "value": "0.00001902",
        "value_int": "1902"
      },
      "item": "BTC",
      "oid": "bc1fb5dc-450e-438e-a80b-96cdaf6ba86c",
      "price": {
        "currency": "AUD",
        "display": "AU$923.00000",
        "display_short": "AU$923.00",
        "value": "923.00000",
        "value_int": "92300000"
      },
      "priority": "1385257384731527",
      "status": "open",
      "type": "ask"
    }
  ]
}`)

func TestHandleResultDataFromOrderRequest(t *testing.T) {
	client := newTestClient(t)

	go client.handle(orderPayload)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for depth data")
	case orders := <-client.Orders:
		t.Logf("Received: %v", orders)
	}
}
