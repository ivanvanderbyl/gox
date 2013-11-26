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

func TestOrderPayloadUnmarshal(t *testing.T) {
	client := newTestClient(t)

	orders, err := client.processOrderResult(orderPayload)
	if err != nil {
		t.Error(err.Error())
	}

	if len(orders) == 0 {
		t.Error("Expected orders, not none")
	}

	valueExpectations := []struct {
		Prop     string
		Expected float64
		Actual   float64
	}{
		{"Amount", 0.01, orders[0].Amount},
		{"EffectiveAmount", 0.00998098, orders[0].EffectiveAmount},
		{"InvalidAmount", 0.00001902, orders[0].InvalidAmount},
		{"Price", 923.0, orders[0].Price},
	}

	for _, e := range valueExpectations {
		if e.Expected != e.Actual {
			t.Errorf("Failed to parse %s property. Got %v != %v", e.Prop, e.Expected, e.Actual)
		}
	}

	stringExpectations := []struct {
		Prop     string
		Expected string
		Actual   string
	}{
		{"Status", "open", orders[0].Status},
		{"Currency", "AUD", orders[0].Currency},
		{"Instrument", "BTC", orders[0].Instrument},
		{"OrderType", "ask", orders[0].OrderType},
		{"OrderId", "bc1fb5dc-450e-438e-a80b-96cdaf6ba86c", orders[0].OrderId},
	}

	for _, e := range stringExpectations {
		if e.Expected != e.Actual {
			t.Errorf("Failed to parse %s property. Got %v != %v", e.Prop, e.Expected, e.Actual)
		}
	}

	if expected, actual := time.Unix(1.385257384e+09, 0), orders[0].Timestamp; expected != actual {
		t.Errorf("Failed to parse Timestamp property. Got %v", actual)
	}
}

// func TestHandleResultDataFromOrderRequest(t *testing.T) {
// 	client := newTestClient(t)

// 	go client.handle(orderPayload)

// 	select {
// 	case <-time.After(100 * time.Millisecond):
// 		t.Error("Timed out waiting for depth data")
// 	case orders := <-client.Orders:
// 		t.Logf("Received: %v", orders)
// 	}
// }
