package mtgox

/*
Depth payload handler

Some code taken from https://github.com/ryanslade/mtgox/
*/

import (
	"encoding/json"
	"strconv"
	"time"
)

type DepthPayload struct {
	StreamHeader
	Depth Depth `json:"depth"`
}

// Represents a Market Depth payload
type Depth struct {
	// Ask or Bid
	TypeString string

	// The price at which volume change happened
	Price int64

	// The volume change
	Volume int64

	// BTC
	Instrument string

	// The currency affected
	Currency string

	// Total volume at this price, after applying the depth update, can be used as a starting point before applying subsequent updates.
	TotalVolume int64

	// When the change happened
	Timestamp time.Time
}

func (d *Depth) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		switch vv := v.(type) {
		case string:
			switch k {
			case "type_str":
				d.TypeString = vv
			case "price_int":
				d.Price, err = strconv.ParseInt(vv, 10, 64)
				if err != nil {
					return err
				}
			case "volume_int":
				d.Volume, err = strconv.ParseInt(vv, 10, 64)
				if err != nil {
					return err
				}
			case "item":
				d.Instrument = vv
			case "currency":
				d.Currency = vv
			case "total_volume_int":
				d.TotalVolume, err = strconv.ParseInt(vv, 10, 64)
				if err != nil {
					return err
				}
			case "now":
				nowInt, err := strconv.ParseInt(vv, 10, 64)
				if err != nil {
					return err
				}
				d.Timestamp = timeFromUnixMicro(nowInt)
			}
		}
	}

	return nil
}

// Handles a depth payload
func (g *Client) handleDepth(data []byte) {
	var depthPayload DepthPayload
	err := json.Unmarshal(data, &depthPayload)
	if err != nil {
		select {
		case g.Errors <- err:
		default:
			// Ignore error if nothing is handling errors so we don't block
		}
	}
	select {
	case g.Depth <- &depthPayload:
	default:
		// Ignore if blocked
	}
}
