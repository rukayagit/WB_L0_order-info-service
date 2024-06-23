package nats

import (
	"github.com/nats-io/stan.go"
)

func ConnectNATS(clusterID, clientID, url string) (stan.Conn, error) {
	nc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, err
	}
	return nc, nil
}

func Subscribe(nc stan.Conn, channelName, queueName string, handler func(*stan.Msg)) (stan.Subscription, error) {
	sub, err := nc.QueueSubscribe(channelName, queueName, handler, stan.DurableName("my-durable"))
	if err != nil {
		return nil, err
	}
	return sub, nil
}
