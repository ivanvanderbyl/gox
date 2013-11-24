package gox

import (
	"strconv"
	"time"
)

type Trade struct {
	Type           string
	Trade_type     string
	Properties     string
	Now            time.Time
	Amount         float64
	Amount_int     int64
	Primary        string
	Price          float64
	Price_int      int64
	Item           string
	Price_currency string
}

type Ticker struct {
	Volume       Value     `json:"vol"`
	Instrument   string    `json:"item"`
	High         Value     `json:"high"`          //Highest value
	Low          Value     `json:"low"`           // Lowest Value
	Last         Value     `json:"last"`          // == Last_local
	LastLocal    Value     `json:"last_local"`    // Last trade in auxilary currency
	LastAux      Value     `json:"last_aux"`      // Last trade converted to auxilary currency
	LastOriginal Value     `json:"last_original"` // Last trade in any currency
	Buy          Value     `json:"buy"`
	Sell         Value     `json:"sell"`
	VWAP         Value     `json:"vwap"` // Volume weighted average price
	Avg          Value     `json:"avg"`  // Averaged price
	Timestamp    EpochTime `json:"now,string"`
}

type Value struct {
	Value        float64 `json:"value,string"`
	ValueInteger int64   `json:"value_int,string"`
	Display      string  `json:"display"`
	DisplayShort string  `json:"display_short"`
	Currency     string  `json:"currency"`
}

type Info struct {
	Created       SimpleTime `json:",string"`
	Id            string
	Index         string
	Language      string
	LastLogin     SimpleTime `json:"last_login,string"`
	Link          string
	Login         string
	Montly_Volume Value
	Trade_fee     float64
	Rights        []string
	Wallets       map[string]Wallet
}

type Wallet struct {
	Balance              Value
	Daily_Withdraw_Limit Value
	Max_Withdraw         Value
	// Monthly_Withdraw_Limit nil
	Open_Orders Value
	Operations  int64
}

type Rate struct {
	To   string
	From string
	Rate float64
}

// Represents a Market Depth payload
type Depth struct {
	ActionId           int       `json:"type"`
	Action             string    `json:"type_str"`
	Volume             float64   `json:"volume,string"`
	VolumeInteger      int64     `json:"volume_int,string"`
	Timestamp          EpochTime `json:"now,string"`
	Price              float64   `json:"price,string"`
	PriceIneteger      int64     `json:"price_int,string"`
	Instrument         string    `json:"item"`
	Currency           string    `json:"currency"`
	TotalVolumeInteger int64     `json:"total_volume_int,string"`
}

type Order struct {
	Oid              string
	Currency         string
	Item             string
	Type             string
	Amount           Value
	Effective_amount Value
	Price            Value
	Status           string
	Date             EpochTime
	Priority         EpochTime `json:",string"`
	Actions          []string
}

type EpochTime struct {
	time.Time
}

type SimpleTime struct {
	time.Time
}

func (t *EpochTime) UnmarshalJSON(b []byte) error {
	result, err := strconv.ParseInt(string(b), 0, 64)
	if err != nil {
		return err
	}
	// convert the unix epoch to a Time object
	*t = EpochTime{time.Unix(0, result*1000)}
	return nil
}
func (t *SimpleTime) UnmarshalJSON(b []byte) error {
	layout := "2006-01-02 15:04:05"
	time, err := time.Parse(layout, string(b))
	if err != nil {
		return err
	}
	*t = SimpleTime{time}
	return nil

}
