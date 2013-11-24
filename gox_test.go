package gox

import (
	"testing"
	"time"
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

var tickerPayload = []byte(`{
  "channel": "eb6aaa11-99d0-4f64-9e8c-1140872a423d",
  "channel_name": "ticker.BTCAUD",
  "op": "private",
  "origin": "broadcast",
  "private": "ticker",
  "ticker": {
    "avg": {
      "currency": "AUD",
      "display": "AU$896.93607",
      "display_short": "AU$896.94",
      "value": "896.93607",
      "value_int": "89693607"
    },
    "buy": {
      "currency": "AUD",
      "display": "AU$870.01350",
      "display_short": "AU$870.01",
      "value": "870.01350",
      "value_int": "87001350"
    },
    "high": {
      "currency": "AUD",
      "display": "AU$970.00000",
      "display_short": "AU$970.00",
      "value": "970.00000",
      "value_int": "97000000"
    },
    "item": "BTC",
    "last": {
      "currency": "AUD",
      "display": "AU$869.99000",
      "display_short": "AU$869.99",
      "value": "869.99000",
      "value_int": "86999000"
    },
    "last_all": {
      "currency": "AUD",
      "display": "AU$883.58189",
      "display_short": "AU$883.58",
      "value": "883.58189",
      "value_int": "88358189"
    },
    "last_local": {
      "currency": "AUD",
      "display": "AU$869.99000",
      "display_short": "AU$869.99",
      "value": "869.99000",
      "value_int": "86999000"
    },
    "last_orig": {
      "currency": "USD",
      "display": "$810.00000",
      "display_short": "$810.00",
      "value": "810.00000",
      "value_int": "81000000"
    },
    "low": {
      "currency": "AUD",
      "display": "AU$820.00000",
      "display_short": "AU$820.00",
      "value": "820.00000",
      "value_int": "82000000"
    },
    "now": "1385273097593509",
    "sell": {
      "currency": "AUD",
      "display": "AU$890.42781",
      "display_short": "AU$890.43",
      "value": "890.42781",
      "value_int": "89042781"
    },
    "vol": {
      "currency": "BTC",
      "display": "573.37262553 BTC",
      "display_short": "573.37 BTC",
      "value": "573.37262553",
      "value_int": "57337262553"
    },
    "vwap": {
      "currency": "AUD",
      "display": "AU$894.22141",
      "display_short": "AU$894.22",
      "value": "894.22141",
      "value_int": "89422141"
    }
  }
}`)

var tradePayload = []byte(`{
  "channel": "dbf1dee9-4f2e-4a08-8cb7-748919a71b21",
  "channel_name": "trade.BTC",
  "op": "private",
  "origin": "broadcast",
  "private": "trade",
  "trade": {
    "amount": 0.02454351,
    "amount_int": "2454351",
    "date": 1.385273108e+09,
    "item": "BTC",
    "price": 810,
    "price_currency": "USD",
    "price_int": "81000000",
    "primary": "Y",
    "properties": "limit",
    "tid": "1385273108844621",
    "trade_type": "bid",
    "type": "trade"
  }
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

	go client.handle(depthPayload)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for depth data")
	case data := <-client.Depth:
		t.Logf("Received Depth: %v", data.Depth.Instrument)
	}
}

func TestHandleDepthDataNonBlocking(t *testing.T) {
	client := newTestClient(t)

	go client.handle(depthPayload)
	go client.handle(depthPayload)
	go client.handle(depthPayload)
	go client.handle(depthPayload)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for depth data")
	case <-client.Depth:
		t.Logf("Received Depth")
	}
}

func TestHandleTickerData(t *testing.T) {
	client := newTestClient(t)

	go client.handle(tickerPayload)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for ticker data")
	case data := <-client.Ticker:
		t.Logf("Received Tick: %v", data.Ticker.Instrument)
		t.Logf("Tick timestamp: %v", data.Ticker.Timestamp)
	}
}

func TestHandleTradeData(t *testing.T) {
	client := newTestClient(t)

	go client.handle(tradePayload)

	select {
	case <-time.After(100 * time.Millisecond):
		t.Error("Timed out waiting for trade data")
	case data := <-client.Trades:
		t.Logf("Received Tick: %v", data.Trade.Amount)
	}
}
