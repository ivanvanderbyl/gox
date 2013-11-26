package mtgox

import (
	"encoding/json"
	"strconv"
	"time"
)

// Order represents a market order from your account.
type Order struct {
	OrderId         string    `json:"oid"`
	Currency        string    `json:"currency"`
	Instrument      string    `json:"item"`
	OrderType       string    `json:"type"`
	Amount          float64   `json:"amount"`
	EffectiveAmount float64   `json:"effective_amount"`
	InvalidAmount   float64   `json:"invalid_amount"`
	Price           float64   `json:"price"`
	Status          string    `json:"status"`
	Date            time.Time `json:"date,string"`
	Priority        uint64    `json:"priority"`
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var raw map[string]string
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		switch k {
		case "oid":
			o.OrderId = v
		case "currency":
			o.Currency = v
		case "item":
			o.Instrument = v
		case "type":
			o.OrderType = v
		case "status":
			o.Status = string(v)
		case "priority":
			o.Priority, err = strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
		case "amount":
			var value Value
			err = json.Unmarshal([]byte(v), &value)
			if err != nil {
				return err
			}

			o.Amount = value.Value
		case "effective_amount":
			var value Value
			err = json.Unmarshal([]byte(v), &value)
			if err != nil {
				return err
			}

			o.EffectiveAmount = value.Value
		case "invalid_amount":
			var value Value
			err = json.Unmarshal([]byte(v), &value)
			if err != nil {
				return err
			}

			o.InvalidAmount = value.Value
		case "price":
			var value Value
			err = json.Unmarshal([]byte(v), &value)
			if err != nil {
				return err
			}

			o.Price = value.Value
		}
	}

	return nil

}
