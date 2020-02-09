package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
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
	publishChannels := flag.String("channels", "general,public,private", "Publishing channels")

	publishChannelList := strings.Split(*publishChannels, ",")
	publishIntervalDuration, err := time.ParseDuration(*publishInterval)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().Unix())
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
			hub.LogInfof("Subscribed: client_id=%s, channels=%s", clientID, strings.Join(m.Channels, ","))
		case "UNSUBSCRIBE":
			err := hub.Unsubscribe(clientID, m.Channels...)
			if err != nil {
				hub.LogError(errors.Wrap(err, "hub.on_receive.unsubscribe"))
			}
			hub.LogInfof("Unsubscribed: client_id=%s, channels=%s", clientID, strings.Join(m.Channels, ","))
		}
	})

	go func() {
		err := hub.ListenAndServe(*addr, *path)
		if err != nil {
			hub.LogPanic(err)
		}
	}()

	go func() {
		for range publishTicker.C {
			// Pick a random channel
			channel := publishChannelList[rand.Intn(len(publishChannelList))]
			message := wspubsub.NewTextMessageFromString(fmt.Sprintf(`{"now": %d}`, time.Now().Unix()))
			_, err := hub.Publish(message, channel)
			if err != nil {
				hub.LogPanic(err)
			}
			hub.LogInfof("Published: channel=%s, message=%s", channel, string(message.Payload))
		}
	}()

	<-quit
}
