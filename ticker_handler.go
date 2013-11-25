package mtgox

import (
	"encoding/json"
	"strconv"
)

type CurrencyValue struct {
	Value        int64
	Display      string
	DisplayShort string
	Currency     string
}

func (tv *CurrencyValue) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		switch k {
		case "display":
			tv.Display = v
		case "display_short":
			tv.DisplayShort = v
		case "currency":
			tv.Currency = v
		case "value_int":
			tv.Value, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Ticker struct {
	High       CurrencyValue `json:"high"`
	Low        CurrencyValue `json:"low"`
	Average    CurrencyValue `json:"avg"`
	Vwap       CurrencyValue `json:"vwap"`
	Volume     CurrencyValue `json:"vol"`
	LastLocal  CurrencyValue `json:"last_local"`
	LastOrigin CurrencyValue `json:"last_orig"`
	LastAll    CurrencyValue `json:"last_all"`
	Last       CurrencyValue `json:"last"`
	Buy        CurrencyValue `json:"buy"`
	Sell       CurrencyValue `json:"sell"`
	Instrument string        `json:"item"`
	Timestamp  EpochTime     `json:"now,string"`
}

type TickerPayload struct {
	StreamHeader
	Ticker Ticker `json:"ticker"`
}

func (g *Client) handleTicker(data []byte) {
	var payload TickerPayload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		select {
		case g.Errors <- err:
		default:
		}
	}

	select {
	case g.Ticker <- &payload:
	default:
	}
}
