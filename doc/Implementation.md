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
    // and send it to the channel
    ConnChan <-Listener.Accept()
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
NB: For clarity, the above example does not include typical error handling 

## Parameter Access and Type Safety






## Pre-Establishment

https://www.ietf.org/archive/id/draft-ietf-taps-impl-10.html#section-3.1 states: 

> The Transport Services system should have a list of supported protocols available, which each have transport features reflecting the capabilities of the protocol. Once an application specifies its Transport Properties, the transport system matches the required and prohibited properties against the transport features of the available protocols.

Our core implementation does _not_ itself provide a list of supported protocols directly. Instead, the applications must explicitly import the corresponding library for each protocol option. These must then be "registered" with the core system, before property matching can occur.

It is of course possible to bundle a sane default selection of protocol libraries into a single "meta" library that can be imported instead.


## Events

The event-driven nature of the spec does not directly translate into idiomatic Go. 
