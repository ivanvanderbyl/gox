# MtGox API Client

[![Build Status](https://travis-ci.org/ivanvanderbyl/mtgox.png?branch=master)](https://travis-ci.org/ivanvanderbyl/mtgox)
[![GoDoc](http://godoc.org/github.com/ivanvanderbyl/mtgox?status.png)](http://godoc.org/github.com/ivanvanderbyl/mtgox)

Mt.Gox Streaming API implementation in Go.

### Documentation

* [Reference](http://godoc.org/github.com/ivanvanderbyl/mtgox)

## Usage

Basic usage requires an API token from your Mt.Gox account.

```go
client, err := mtgox.New( "KEY", "SECRET", "AUD", "USD")
if err != nil {
  // Handle connection error
}

// Start message receive routine
client.Start()

for {
  select {
  case orders := <-client.Orders:
    fmt.Printf("Orders: %s\n", PrettyPrintJson(orders))

  case trade := <-client.Trades:
    fmt.Printf("Trade: %s\n", PrettyPrintJson(trade))

  case tick := <-client.Ticker:
    fmt.Printf("Tick: %s\n", PrettyPrintJson(tick))

  case depth := <-client.Depth:
    fmt.Printf("Depth: %s\n", PrettyPrintJson(depth))
  }
}
```

## Features

- Receives streaming events for `depth`, `trade`, `ticker`
- All API methods to placing orders (in the works)
- All Account and Wallet API methods (in the works)

