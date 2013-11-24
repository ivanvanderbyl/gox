package gox

import (
	"testing"
)

var depthPayload = []byte(`{
  "channel": "296ee352-dd5d-46f3-9bea-5e39dede2005",
  "channel_name": "depth.BTCAUD",
  "op": "private",
  "origin": "broadcast",
  "private": "depth",
  "ticker": null,
  "depth": {
    "type": 1,
    "type_str": "ask",
    "volume": "-7.95",
    "volume_int": "-795000000",
    "now": "2013-11-24T16:47:39.960998+11:00",
    "price": "919.80358",
    "price_int": "91980358",
    "item": "BTC",
    "currency": "AUD",
    "total_volume_int": "1590000000"
  },
  "info": null
}`)

func newTestClient(t *testing.T) *Gox {
	client, err := NewWithConnection("123abcde-4567-8910-1112-74e2ef79f40d", "VE9QU0VDUkVU", nil)
	if err != nil {
		t.Error(err.Error())
	}
	return client
}

func TestHandleDepthData(t *testing.T) {
	client := newTestClient(t)

	client.handle(depthPayload)
}
