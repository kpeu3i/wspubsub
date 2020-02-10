// Package wspubsub provides an easy way to publish/subscribe and receive messages
// over WebSocket protocol hiding the details of the underlying transport.
// The library based on the idea that all messages are published to channels.
// Subscribers will receive messages published to the channels to which they subscribe.
// The hub is responsible for defining the channels to which subscribers can subscribe.
package wspubsub
