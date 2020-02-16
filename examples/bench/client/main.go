package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kpeu3i/wspubsub"
)

var (
	logger wspubsub.Logger
)

func init() {
	rand.Seed(time.Now().Unix())
	logger = wspubsub.NewLogrusLogger(wspubsub.NewLogrusLoggerOptions())
}

func main() {
	url := flag.String("addr", "ws://localhost:8080/", "Server URL")
	connectionCount := flag.Int("connections", 5000, "Total connection count")
	workerCount := flag.Int("workers", 1, "Total worker count")
	channels := flag.String("channels", "X,Y,Z", "One of possible subscription channels will be used")
	wait := flag.String("wait", "10s", "Wait time after all connections will be established")
	readTimeout := flag.String("read_timeout", "60s", "Read timeout")
	writeTimeout := flag.String("write_timeout", "10s", "Write timeout")

	channelList := strings.Split(*channels, ",")

	readTimeoutDuration, err := time.ParseDuration(*readTimeout)
	if err != nil {
		panic(err)
	}

	writeTimeoutDuration, err := time.ParseDuration(*writeTimeout)
	if err != nil {
		panic(err)
	}

	waitDuration, err := time.ParseDuration(*wait)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	connectionsMutex := sync.RWMutex{}
	connections := make([]*websocket.Conn, 0, *connectionCount)
	tasks := make(chan func() error)

	logger.Println("Starting...")

	for i := 1; i <= *workerCount; i++ {
		go runWorker(tasks)
	}

	now := time.Now()

	for i := 1; i <= *connectionCount; i++ {
		wg.Add(1)
		tasks <- func() error {
			defer wg.Done()

			connection, err := subscribe(*url, channelList, readTimeoutDuration, writeTimeoutDuration)
			if err != nil {
				return err
			}

			connectionsMutex.Lock()
			connections = append(connections, connection)
			connectionsMutex.Unlock()

			if i%100 == 0 {
				logger.Printf("Connected: %d\n", i)
			}

			return nil
		}
	}

	wg.Wait()

	took := time.Since(now)

	logger.Println("Done!")
	logger.Printf("Waiting %s...", waitDuration)

	time.Sleep(waitDuration)

	logger.Println("Stopping...")

	close(tasks)

	for _, connection := range connections {
		_ = connection.Close()
	}

	logger.Println("Connections:", len(connections))
	logger.Println("Took:", took)
	logger.Println("RPS:", math.Round(float64(len(connections))/took.Seconds()))
}

func subscribe(url string, channels []string, readTimeout, writeTimeout time.Duration) (*websocket.Conn, error) {
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	err = connection.SetReadDeadline(time.Now().Add(readTimeout))
	if err != nil {
		return nil, err
	}

	err = connection.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return nil, err
	}

	connection.SetPongHandler(func(string) error {
		return connection.SetReadDeadline(time.Now().Add(readTimeout))
	})

	err = connection.WriteMessage(websocket.TextMessage, subscribeMessagePayload(channels))
	if err != nil {
		return nil, err
	}

	// Reader
	go func() {
		for {
			err = connection.SetReadDeadline(time.Now().Add(60 * time.Second))
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					logger.Errorf("Read error: %s\n", err)
					return
				}

				return
			}

			_, _, err := connection.ReadMessage()
			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					logger.Errorf("Read error: %s\n", err)
					return
				}

				return
			}
		}
	}()

	return connection, nil
}

func runWorker(tasks <-chan func() error) {
	for task := range tasks {
		err := task()
		if err != nil {
			logger.Errorf("Failed to process a task: %s\n", err)
		}
	}
}

func subscribeMessagePayload(channels []string) []byte {
	channel := channels[rand.Intn(len(channels))]
	message := fmt.Sprintf(`{"command": "SUBSCRIBE", "channels": ["%s"]}`, channel)

	return []byte(message)
}
