# PANAPI's implementation of the TAPS API (draft)

Considerations for mapping the specification laid out in  https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html to our Go implementation. In particular, we look at the [Usage Examples](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#name-usage-examples) to think about reasonable ways to implement the API in Go.

## Preconnection specification

The [Server Example](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#name-server-example) from the draft spec envisions setting up a Preconnection looks as follows:

```
LocalSpecifier := NewLocalEndpoint()
LocalSpecifier.WithInterface("any")
LocalSpecifier.WithService("https")

TransportProperties := NewTransportProperties()
TransportProperties.Require(preserve-msg-boundaries)
// Reliable Data Transfer and Preserve Order are Required by default

SecurityParameters := NewSecurityParameters()
SecurityParameters.Set(identity, myIdentity)
SecurityParameters.Set(key-pair, myPrivateKey, myPublicKey)

// Specifying a Remote Endpoint is optional when using Listen()
Preconnection := NewPreconnection(LocalSpecifier,
                                  TransportProperties,
                                  SecurityParameters)
```

There are different ways to map this to Go

### Option 1: String arguments

The easiest way, but probably also the most error-prone way for the developer, would be to use strings here. 

```Go
LocalSpecifier := panapi.NewLocalEndpoint()
LocalSpecifier.WithInterface("any")
LocalSpecifier.WithService("https")

TransportProperties := panapi.NewTransportProperties()
TransportProperties.Require("preserve-msg-boundaries")

SecurityParameters := panapi.NewSecurityParameters()
SecurityParameters.Set("identity", myIdentity)
SecurityParameters.Set("key-pair", myPrivateKey, myPublicKey)

Preconnection := panapi.NewPreconnection(LocalSpecifier,
                                  TransportProperties,
                                  SecurityParameters)
```



### Option 2: Constants
```Go
LocalSpecifier := panapi.NewLocalEndpoint()
LocalSpecifier.WithInterface(panapi.ANY)
LocalSpecifier.WithService(panapi.HTTPS)

TransportProperties := panapi.NewTransportProperties()
TransportProperties.Require(panapi.PRESERVE_MSG_BOUNDARIES)
// Reliable Data Transfer and Preserve Order are Required by default

SecurityParameters := panapi.NewSecurityParameters()
SecurityParameters.Set(panapi.IDENTITY, myIdentity)
SecurityParameters.Set(panapi.KEY_PAIR, myPrivateKey, myPublicKey)

// Specifying a Remote Endpoint is optional when using Listen()
Preconnection := panapi.NewPreconnection(LocalSpecifier,
                                  TransportProperties,
                                  SecurityParameters)
```


### Option 3: Composition of structs
This option probably departs the most from the example in the draft. The question is, to what degree would the _intent_ of the draft still be met.

```Go
import (
    "panapi/service/https"
)

// ...

LocalSpecifier := panapi.NewLocalEndpoint()
// all possible interfaces, equivalent to "any" 
LocalSpecifier.WithInterface(panapi.NetworkInterfaces())
LocalSpecifier.WithService(https.NewService())

TransportProperties := panapi.NewTransportProperties()
TransportProperties.Require(panapi.PRESERVE_MSG_BOUNDARIES)

SecurityParameters := panapi.SecurityParameters{
    Identity: myIdentity,
    PrivateKey: myPrivateKey,
    PublicKey: myPublicKey,
}

// Specifying a Remote Endpoint is optional when using Listen()
Preconnection := panapi.NewPreconnection(LocalSpecifier,
                                  TransportProperties,
                                  SecurityParameters)
```

## Asynchronous Operation

The [Server Example](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#name-server-example) from the draft spec envisions asynchronous connection handling like this:

```
Listener := Preconnection.Listen()

Listener -> ConnectionReceived<Connection>

// Only receive complete messages in a Conn.Received handler
Connection.Receive()

Connection -> Received<messageDataRequest, messageContext>

//---- Receive event handler begin ----
Connection.Send(messageDataResponse)
Connection.Close()

// Stop listening for incoming Connections
// (this example supports only one Connection)
Listener.Stop()
//---- Receive event handler end ----
```

Again, there are numerous ways to do this in Go.

### Option 1: Heavy use of Go channels for events

This would associate each different type of possible event parameter to its own channel.

```Go
Listener := Preconnection.Listen()

//---- Loop to handle multiple connections begin ----

// New connections are sent on the "ConnectionReceived" channel provided by Listener
// (Execution blocks, until a Connection is received)
Connection := <- Listener.ConnectionReceived

//---- Go-routine for asynchronous handling of a received connection begin ----
go func(connection panapi.Connection) {

    // Calling "Receive()" on a Connection indicates our desire to
    // later read a complete messages from the "Received" channel.
    // (Nonblocking)
    connection.Receive()

    // Block, until a complete message is received
    // (Message is a struct containing the fields "Request", "Context", etc)
    Message := <- connection.Received

    // Process message

    connection.Send(messageDataResponse)
    connection.Close()
}(Connection)
//---- Go-routine end ----

//---- Loop end ----

Listener.Stop()
```

### Option 2: Use callback handlers

```Go
Listener := Preconnection.Listen()

// the inner callback gets executed, when a complete message is available
Listener.ReceiveFunc(func(conn panapi.Connection, msg panapi.Message) {
    messageDataRequest, messageContext := msg.Request, msg.Context
    ...
    conn.Send(messageDataResponse)
    conn.Close()
}

Listener.Stop()
```

