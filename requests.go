package mtgox

import (
	"time"
)

// Dispatches a request for private/info, returning an info payload or timing out
func (g *Client) RequestInfo() *Info {
	g.call("private/info", nil)

	select {
	case <-time.After(10 * time.Second):
		return nil
	case info := <-g.Info:
		return info
	}
}

// func (api *StreamingApi) RequestOrders() (c chan []Order) {
// 	api.call("private/orders", nil)
// 	return api.Orders
// }

// func (g *Client) RequestOrders() []Order {
// 	g.call("private/orders", nil)

// 	select {
// 	case <-time.After(10 * time.Second):
// 		return nil
// 	case orders := <-g.Orders:
// 		return orders
// 	}
// }

func (g *Client) RequestOrders() {
	g.call("private/orders", nil)
}
