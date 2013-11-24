package gox

import (
	"time"
)

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

type SimpleTime struct {
	time.Time
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
