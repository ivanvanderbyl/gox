// Copyright 2013 Ivan Vanderbyl. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package mtgox implements a complete MtGox streaming API client.

Overview

The Client type represents a connection to the Mt.Gox streaming API using
WebSockets (Powered by the modern Gorilla WebSocket package). After creating
a Client instance, you should call client.Start to begin handling received
messages.

  client, err := mtgox.New( "KEY", "SECRET", "Currencies"...)
  if err != nil {
    // Handle connection error
  }

  // Start message receive routine
  client.Start()

From here all messages will be sent to their respective channels.
*/
package mtgox
