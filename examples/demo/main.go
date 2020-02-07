package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kpeu3i/wspubsub"
	"github.com/pkg/errors"
)

type Message struct {
	Command  string   `json:"command"`
	Channels []string `json:"channels"`
}

func main() {
	addr := flag.String("addr", "localhost:8080", "Server host and port")
	path := flag.String("path", "/", "URL")
	publishInterval := flag.String("publish", "1s", "Messages publishing interval")
	publishChannels := flag.String("channels", "X,Y,Z", "Publishing channels")
	message := flag.String("message", `{"demo": true}`, "Message")

	publishChannelList := strings.Split(*publishChannels, ",")
	publishIntervalDuration, err := time.ParseDuration(*publishInterval)
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

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
		err := hub.ListenAndServe(*addr, *path)
		if err != nil {
			hub.LogPanic(err)
		}
	}()

	go func() {
		message := wspubsub.NewTextMessageFromString(*message)
		for range publishTicker.C {
			_, err := hub.Publish(message, publishChannelList...)
			if err != nil {
				hub.LogPanic(err)
			}
		}
	}()

	<-quit
}
