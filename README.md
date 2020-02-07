# WSPubSub

[![CircleCI Build Status](https://circleci.com/gh/kpeu3i/wspubsub.svg?style=shield)](https://circleci.com/gh/kpeu3i/wspubsub)
[![Coverage Status](https://codecov.io/gh/kpeu3i/wspubsub/graphs/badge.svg)](https://codecov.io/gh/kpeu3i/wspubsub/graphs/badge.svg)
[![GoDoc](https://godoc.org/github.com/kpeu3i/wspubsub?status.svg)](https://godoc.org/github.com/kpeu3i/wspubsub)
[![MIT Licensed](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

WSPubSub library is a Go implementation of channels based [pub/sub pattern](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern) over WebSocket protocol.

This library completely hides interaction with the transport layer like: connection upgrading, disconnecting, ping/pong etc.
Thanks to this, you can focus only on your tasks.

Client interaction with the library mainly occurs through the `hub` API. The only two steps to publish messages are required:
* Register receive handler using `hub.OnReceive`
* Subscribe clients using `hub.Subscribe`

Now you are ready to publish messages to different channels using `hub.Publish`.
Users who have been subscribed to those channels will receive the messages!

## Install

Use go get to install the latest version of the library:

    go get github.com/kpeu3i/wspubsub

Next, include `wspubsub` in your application:

```go
import "github.com/kpeu3i/wspubsub"
```

## Usage

A minimal working example is listed below:

```go
hub := wspubsub.NewDefaultHub()
defer hub.Close()

publishTicker := time.NewTicker(publishIntervalDuration)
defer publishTicker.Stop()

hub.OnReceive(func(clientID wspubsub.UUID, message wspubsub.Message) {
    m := Message{}
    err := json.Unmarshal(message.Payload, &m)
    if err != nil {
        hub.LogError(errors.Wrap(err, "hub.on_receive.unmarshal"))
        return
    }

    switch m.Command {
    case "SUBSCRIBE":
        err := hub.Subscribe(clientID, m.Channels...)
        if err != nil {
            hub.LogError(errors.Wrap(err, "hub.on_receive.subscribe"))
        }
    case "UNSUBSCRIBE":
        err := hub.Unsubscribe(clientID, m.Channels...)
        if err != nil {
            hub.LogError(errors.Wrap(err, "hub.on_receive.unsubscribe"))
        }
    }
})

go func() {
    err := hub.ListenAndServe("localhost:8080", "/")
    if err != nil {
        hub.LogPanic(err)
    }
}()

go func() {
    message := wspubsub.NewTextMessageFromString(`{"demo": true}`)
    for range publishTicker.C {
        _, err := hub.Publish(message, publishChannelList...)
        if err != nil {
            hub.LogPanic(err)
        }
    }
}()
````
More examples you can find in [examples](https://github.com/kpeu3i/wspubsub/tree/master/examples/) directory.