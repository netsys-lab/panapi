# Notes on Implementation

We implement PANAPI in Go, a statically typed, compiled programming
language with memory safety and powerful native primitives for
concurrency via message passing between light-weight processes
("goroutines"). Go is not object oriented but instead relies on type
composition and implements dynamic dispatch via interfaces to enable
effective code reuse.

Over the years, these language properties, together with certain
agreed-upon "best practices", have caused the community of Go
developers to converge on a coding style consensus that we now usually
refer to as "idiomatic Go".

Some conflicts exist between the TAPS specifications on the one hand
and our implementation on the other. We are not sure of the severity
in each case, but wish to point out any departures from the TAPS
vision here.

## Asnychronous, Event-driven Interaction Patterns

https://www.ietf.org/archive/id/draft-ietf-taps-impl-10.html#abstract states:

> The Transport Services system [...] defines a protocol-independent
> Transport Services Application Programming Interface (API) that is
> based on an asynchronous, event-driven interaction pattern.

https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#section-1.1 states:

> Events are sent to the application or application-supplied code
> [...] for processing; the details of event processing are platform-
> and implementation-specific.

Finally, https://www.ietf.org/archive/id/draft-ietf-taps-arch-10.html#section-2.1 states:

>  2.1. Event-Driven API
>
> Originally, sockets presented a blocking interface for establishing
> connections and transferring data. However, most modern applications
> interact with the network asynchronously. Emulation of an
> asynchronous interface using sockets generally uses a try-and-fail
> model. If the application wants to read, but data has not yet been
> received from the peer, the call to read will fail. The application
> then waits and can try again later.
> 
> In contrast to sockets, all interaction with a Transport Services
> system is expected to be asynchronous, and use an event-driven model
> (see Section 4.1.6). For example, if the application wants to read,
> its call to read will not complete immediately, but will deliver an
> event containing the received data once it is available. Error
> handling is also asynchronous; a failure to send results in an
> asynchronous send error as an event.
>
> The Transport Services API also delivers events regarding the
> lifetime of a connection and changes in the available network links,
> which were not previously made explicit in sockets.
> 
> Using asynchronous events allows for a more natural interaction model
> when establishing connections and transferring data. Events in time
> more closely reflect the nature of interactions over networks, as
> opposed to how sockets represent network resources as file system
> objects that may be temporarily unavailable.
>
> Separate from events, callbacks are also provided for asynchronous
> interactions with the API not directly related to events on the
> network or network interfaces.

Consider the [server
example](https://www.ietf.org/archive/id/draft-ietf-taps-interface-13.html#name-server-example)
from the Draft:

```
Listener := Preconnection.Listen()

Listener -> ConnectionReceived<Connection>

// Only receive complete messages in a Conn.Received handler
Connection.Receive()

Connection -> Received<messageDataRequest, messageContext>
```

Here, `Preconnection.Listen()` and `Connection.Receive()` are implied
to be non-blocking calls that initiate some background processes which
at some point emit the associated Events once they are completed.

Let's assume we have the _blocking_ calls `Listener.Accept()` and
`Connection.Receive()` instead.

The semantically equivalent, _synchronous_ variant of the above code
example could look as follows:

```Go
// does not block (same as above)
Listener := Preconnection.Listen()

// blocks until a Connection can be accepted
Connection := Listener.Accept()

// blocks until a complete Message can be received
messageDataRequest, messageContext := Connection.Receive()
```

In Go, _asynchronous_ processing is usually implemented
explicitly on the application side using Go's native concurrency
primitives as needed.

Applications that can not afford to halt execution to wait for a
blocking function call to succeed would simply dispatch that call in a
concurrent goroutine and read the eventual result from a channel
constructed for that purpose.

Without the need for any dedicated Event system, we could turn the
above code into an asynchronous variant:

```Go
Listener := Preconnection.Listen()

// create a channel to later read Connections from
ConnChan := make(chan Connection)

// dispatch an asynchronous goroutine that writes the next
// 'Accept()'ed Connection to ConnChan
go func() {
    // wait until Connection can be accepted
    Connection := Listener.Accept()
    // and send it to the channel
    ConnChan <- Connection
}()

// dispatch an asynchronous goroutine that waits for a 
// Connection, receives a Message and processes it...
go func() {
    // wait until a Connection is sent on ConnChan
    Connection := <-ConnChan
    //wait until a complete Message can be received
    messageDataRequest, messageContext := Connection.Receive()
    // process message
    ...
}()

// regardless of any calls still blocking inside the goroutines
// execution immediately resumes here
...
```

For clarity, the above example does not include typical error
handling. Also, in this case, asynchronously accepting connections
makes little sense and the use of a channel to pass Connection objects
is therefore slight overkill. "proper" implementation in idiomatic Go
would rather look something like this:

```Go
listener, err := Preconnection.Listen()
if err != nil {
    // handle error, exit gracefully
    log.Fatalf("Could not open listener: %s", err)
}

// server main loop
for {
    // block until a new connection is received
    conn, err := Listener.Accept()
    if err != nil {
        // handle error
        log.Printf("could not accept connection: %s", err)
        // re-enter the server main loop from the top
        continue
    }

    // dispatch goroutine to handle connection asynchronously
    go func(c Connection) {
        // wait until a complete Message can be received
        messageDataRequest, messageContext := c.Receive()
        // process message
        ...
    }(conn)
}

```

We suspect that TAPS is specified as an event-based asynchronous API
because most programming languages don't offer such high-level
concurrency features, at least not on a comparable level of
convenience.

With this in mind, we have (tentatively) converged on an
implementation strategy of our TAPS-like API in Go that is centered
around *blocking* calls, which can nevertheless safely be put into
asynchronous goroutines as needed. For now, we will simply not include
any kind of event system, simultaneously getting rid of any associated
requirements for handcrafted event types and necessarily idiosyncratic
Error-handling. Instead, an application can selectively decide to
"handle events" by calling the corresponding blocking function in an
asynchronous goroutine and thereby stay informed about, e.g., path
changes:

```Go
listener, err := Preconnection.Listen()
if err != nil {
    // handle error, exit gracefully
    log.Fatalf("Could not open listener: %s", err)
}

// server main loop
for {
    // block until a new connection is received
    conn, err := Listener.Accept()
    if err != nil {
        // handle error
        log.Printf("could not accept connection: %s", err)
        // re-enter the server main loop from the top
        continue
    }
    
    // dispatch goroutine to worry about the underlying 
    // path changing for conn
    go func(c Connection) {
        // wait for a path change ocurring
        err := c.PathChange()
        if err != nil {
            // handle error
            if err == ErrClosed {
                // connection closed without any path change
            } else {
                // some other error occured
                log.Printf("while waiting for path change: %s", err)
            }
            return
        }
        // react to path change
        ...
        
    }(conn)

    // dispatch goroutine to handle connection asynchronously
    go func(c Connection) {
        // wait until a complete Message can be received
        messageDataRequest, messageContext := c.Receive()
        // process message
        ...
    }(conn)
}



```
### Events and blocking Functions covering them

`func (*Preconnection) Initiate() (Connection, error)` and `func(*Preconnection) Rendezvous()`(Connection, error)` cover the Events:

 - `Ready<>`: When `Initiate()` or `Rendezvous` return a `Connection` and no `error`
 - `EstablishmentError<>`: When `Initiate` returns an `error`
 - `RendezvousDone<Connection>`: When `Rendezvous` returns a `Connection` and no `error`

`func (*Listener) Accept() (Connection, error)` covers the following Events:

 - `ConnectionReceived<Connection>`: When `Accept` returns a Connection but no `error`
 - `Stopped<>`: When `Accept` returns `StoppedError`
 - `EstablishmentError<>`: When `Accept` returns any other `error`

`func (*Connection) Send(Message) error` covers the following Events:

 - `Sent<messageContext>`: When `Send` returns no `error`
 - `Expired<messageContext>`: When `Send` returns `ErrorExpired`
 - `ConnectionError<>`: When `Send` returns an `error`, because the underlying `Connection` closed due to the `error`
 - `SendError<messageContext, reason?>`: When `Send` returns any other `error`
 
`func (*Connection) Receive() (Message, error)` covers the following Events:
 
 - `Received<messageData, messageContext>`: When `Receive` returns a `Message` but no `error`
 - `ReceivedPartial<messageData, messageContext, endOfMessage>`: When `Receive` returns a `Message` with the `EndOfMessage` flag set to `false`.
 - `ConnectionError<>`: When `Receive` returns an `error`, because the underlying `Connection` closed due to the `error`
 - `ReceiveError<messageContext, reason?>`: When `Receive` returns any other `error`

`func (*Connection) Close() error` covers the following Events:

 - `Closed<>`: When `Close` returns no `error`
 - `ConnectionError<>`: When `Close` returns an `error`, because the underlying `Connection` closed due to the `error`


The following Events are covered by dedicated blocking functions for this purpose
 
 - `SoftError<>`: covered by `func (*Connection) SoftError() error`, returns an ICMP `error` if one is received on the underlying `Connection`
 - `PathChange<>`: covered by `func (*Connection) PathChange() error`, returns an `error` if `Connection` closed without a Path Change

This should completely cover all potential control flow patters that are enabled by the Events from the TAPS Spec.

## Property/Parameter Access and Type Safety

To manage Properties, we decided to use Go-native struct fields
instead of more error-prone string-based approaches.

While it _is_ possible to access struct fields using strings in Go,
the process has some overhead and ignores the benefits of Go's static
type system "in favor of" unneccessary runtime errors. Using
reflection, it _would_ be possible to have fuzzy string matching
against struct field names. A `Set` function could store a value for a
property name, which is stripped of case and non-alphabetic characters
before being matched against the (equally stripped) exported field
names of the struct. The type of value can be checked for
"assignability" to the type of the targeted property field and
otherwise return an error. This function would allow you to say:

```Go
tp := NewTransportProperties()
err := tp.Set("preserve-msg-boundaries", Require)
if err != nil {
    ... // handle runtime error
}
```
In idiomatic Go, you would (and should) instead say:

```
tp.PreserveMsgBoundaries = Require
```

For this reason, and to keep the API as concise as possible, we
have currently disabled `Get`ing and `Set`ing Properties via strings.

## Pre-Establishment

https://www.ietf.org/archive/id/draft-ietf-taps-impl-10.html#section-3.1
states:

> The Transport Services system should have a list of supported protocols available, which each have transport features reflecting the capabilities of the protocol. Once an application specifies its Transport Properties, the transport system matches the required and prohibited properties against the transport features of the available protocols.

Our core implementation does _not_ itself provide a list of supported protocols directly. Instead, the applications must explicitly import the corresponding library for each protocol option. These must then be "registered" with the core system, before property matching can occur.

It is of course possible to bundle a sane default selection of protocol libraries into a single "meta" library that can be imported instead.
