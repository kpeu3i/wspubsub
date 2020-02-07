package main

import (
	"encoding/json"
	"expvar"
	"flag"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/kpeu3i/wspubsub"
	"github.com/pkg/errors"
)

var (
	publishMessageCount int64
	publishTime         int64
	publishCount        int64
)

type Message struct {
	Command  string   `json:"command"`
	Channels []string `json:"channels"`
}

func init() {
	rand.Seed(time.Now().Unix())
	expvar.Publish("Goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
}

func main() {
	addr := flag.String("addr", "localhost:8080", "Server host and port")
	path := flag.String("path", "/", "URL")
	publishInterval := flag.String("publish", "50ms", "Messages publishing interval")
	publishChannels := flag.String("channels", "X,Y,Z", "Publishing channels")
	messageSize := flag.Int("message", 5*1024, "Message size (bytes)")

	publishChannelList := strings.Split(*publishChannels, ",")
	publishIntervalDuration, err := time.ParseDuration(*publishInterval)
	if err != nil {
		panic(err)
	}

	publishTicker := time.NewTicker(publishIntervalDuration)
	payload := messagePayload(*messageSize)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	hub := wspubsub.NewDefaultHub()

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

	http.Handle(*path, hub)

	go func() {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			hub.LogPanic(err)
		}
	}()

	go publish(hub, publishTicker, publishChannelList, payload)

	<-quit

	publishTicker.Stop()

	err = hub.Close()
	if err != nil {
		hub.LogPanic(err)
	}

	if publishTime > 0 && publishCount > 0 {
		hub.LogPrintf("Total publish time: %s\n", time.Duration(publishTime))
		hub.LogPrintf("Total messages published: %d\n", publishMessageCount)
		hub.LogPrintf("Single message publish time: %s\n", time.Duration(publishTime/publishCount))
		hub.LogPrintf("Publish RPS: %d\n", int(float64(publishMessageCount)/(float64(publishTime)/1e9)))
	}
}

func publish(hub *wspubsub.Hub, ticker *time.Ticker, channels []string, payload []byte) {
	hub.LogInfoln("Starting publishing...")
	message := wspubsub.NewBinaryMessage(payload)
	for range ticker.C {
		channel := channels[rand.Intn(len(channels))]
		now := time.Now()
		clientCount, err := hub.Publish(message, channel)
		if err != nil {
			hub.LogPanic(err)
		}

		if clientCount > 0 {
			publishMessageCount += int64(clientCount)
			publishTime += time.Since(now).Nanoseconds()
			publishCount++
		}
	}
}

func messagePayload(size int) []byte {
	message := make([]byte, 0, size)
	for i := 0; i < size; i++ {
		message = append(message, byte(rand.Intn(math.MaxUint8+1)))
	}

	return message
}
