package mtgox

import (
	"encoding/json"
	"strconv"
	"time"
)

// Order represents a market order from your account.
type Order struct {
	// OrderId is the unique order identifier
	OrderId string `json:"oid,string"`

	// Currency represents the external fiat currency
	Currency string `json:"currency,string"`

	// Instrument is the object being traded
	Instrument string `json:"item,string"`

	// OrderType in the market, one of `bid` or `ask`
	OrderType string `json:"type,string"`

	// Amount is the requested trade amount in BTC
	Amount float64 `json:"amount"`

	// EffectiveAmount is the actual amount traded in BTC
	EffectiveAmount float64 `json:"effective_amount"`

	// InvalidAmount is the amount which could not traded in BTC
	InvalidAmount float64 `json:"invalid_amount"`

	// Price is the amount paid in fiat currency
	Price float64 `json:"price"`

	// Status is the current status of the order
	Status string `json:"status,string"`

	// Timestamp of the order taking place
	Timestamp time.Time `json:"date,string"`

	// Priority is a unique microsecond timestamp of the order
	// (Not sure what the actual use of this is)
	Priority uint64 `json:"priority,string"`
}

func (o *Order) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		switch vv := v.(type) {
		case string:
			switch k {
			case "oid":
				o.OrderId = vv
			case "currency":
				o.Currency = vv
			case "item":
				o.Instrument = vv
			case "type":
				o.OrderType = vv
			case "status":
				o.Status = vv
			case "priority":
				o.Priority, err = strconv.ParseUint(vv, 10, 64)
				if err != nil {
					return err
				}
			}
		case map[string]interface{}:
			switch k {
			case "amount":
				if val, ok := vv["value_int"].(string); ok {
					valFloat, err := strconv.ParseFloat(val, 64)
					if err != nil {
						return err
					}
					if currency, ok := vv["currency"].(string); ok {
						o.Amount = valFloat / currencyDivisions[currency]
					}
				}

			case "effective_amount":
				if val, ok := vv["value_int"].(string); ok {
					valFloat, err := strconv.ParseFloat(val, 64)
					if err != nil {
						return err
					}
					if currency, ok := vv["currency"].(string); ok {
						o.EffectiveAmount = valFloat / currencyDivisions[currency]
					}
				}

			case "invalid_amount":
				if val, ok := vv["value_int"].(string); ok {
					valFloat, err := strconv.ParseFloat(val, 64)
					if err != nil {
						return err
					}
					if currency, ok := vv["currency"].(string); ok {
						o.InvalidAmount = valFloat / currencyDivisions[currency]
					}
				}

			case "price":
				if val, ok := vv["value_int"].(string); ok {
					valFloat, err := strconv.ParseFloat(val, 64)
					if err != nil {
						return err
					}

					if currency, ok := vv["currency"].(string); ok {
						o.Price = valFloat / currencyDivisions[currency]
					}
				}
			}

		case float64:
			switch k {
			case "date":
				o.Timestamp = time.Unix(int64(vv), 0)
			}
			// default:
			// fmt.Printf("Got unknown type: %v\n", vv)
		}
	}

	return nil

}
