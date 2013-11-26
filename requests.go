package mtgox

import (
	"errors"
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

func (g *Client) RequestOrders() (<-chan []Order, error) {
	reqId, err := g.call("private/orders", nil)
	if err != nil {
		return nil, err
	}

	normalisedReplyChan := make(chan []Order)
	go func() {
		replyChan := make(chan []byte, 1)
		defer close(replyChan)
		g.enqueuePendingRequest(reqId, replyChan)
		data := <-replyChan
		result, err := g.processOrderResult(data)

		if err != nil {
			return
		}

		normalisedReplyChan <- result
	}()

	return normalisedReplyChan, nil
}

// RequestOrderLag returns the lag time for executing orders
// This method will block for up to 5 seconds before timing out
func (g *Client) RequestOrderLag() (time.Duration, error) {
	reqId, err := g.call("order/lag", nil)
	if err != nil {
		return 0, nil
	}

	replyChan := make(chan []byte, 1)
	g.enqueuePendingRequest(reqId, replyChan)

	select {
	case <-time.After(5 * time.Second):
		return 0, errors.New("Failed to receive lag response within 5 seconds")
	case lagData := <-replyChan:
		close(replyChan)
		return g.processLagResult(lagData)
	}
}
