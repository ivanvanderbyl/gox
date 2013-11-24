package gox

import (
	"time"
)

// Dispatches a request for private/info, returning an info payload or timing out
func (g *Gox) RequestInfo() *Info {
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

// func (g *Gox) RequestOrders() []Order {
// 	g.call("private/orders", nil)

// 	select {
// 	case <-time.After(10 * time.Second):
// 		return nil
// 	case orders := <-g.Orders:
// 		return orders
// 	}
// }

func (g *Gox) RequestOrders() {
	g.call("private/orders", nil)
}
