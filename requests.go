package mtgox

import (
	"errors"
	"time"
)

// RequestInfo dispatches a request for `private/info`, returning an info
// payload or timing out after 10 seconds.
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

// RequestOrders fetches the open orders for your account
func (g *Client) RequestOrders() (<-chan []Order, error) {
	reqId, err := g.call("private/orders", nil)
	if err != nil {
		return nil, err
	}

	normalisedReplyChan := make(chan []Order)
	go func() {
		replyChan := make(chan []byte, 1)
		defer close(replyChan)
		defer close(normalisedReplyChan)

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

// GetHistory returns the histroy for the given wallet ID
func (g *Client) GetHistory(walletID string) []Order {
	return []Order{}
}

// QueryOrder returns updated information for the given order
func (g *Client) QueryOrder(orderID string) (Order, error) {
	return Order{}, nil
}

// PlaceOrder adds the given Order to the market, returning the new ID
func (g *Client) PlaceOrder(o Order) (string, error) {
	if o.OrderID != "" {
		return "", errors.New("order should not have an ID yet")
	}
	return "", nil
}

// QuoteOrder retrieves a quote for the given Order
func (g *Client) QuoteOrder(o Order) (string, error) {
	return "", nil
}

// CancelOrder will cancel the given order if already added to the market
func (g *Client) CancelOrder(o Order) (string, error) {
	return "", nil
}
