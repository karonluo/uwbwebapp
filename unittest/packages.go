package main

import (
	"fmt"

	emqx "github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient emqx.Client

func ReceiveMessage(client emqx.Client, msg emqx.Message) {
	fmt.Println(client.IsConnected())
	fmt.Printf("Received message: %s from topic: %s\r\n", string(msg.Payload()), msg.Topic())
}

func ConnectionMQTTServer(clientId string) {
	opts := emqx.NewClientOptions()
	dsn := fmt.Sprintf("tcp://%s:%d", "127.0.0.1", 1883)
	opts.AddBroker(dsn)
	opts.SetClientID(clientId)
	opts.SetUsername("admin")
	opts.SetPassword("pass@@word123")
	MQTTClient = emqx.NewClient(opts)
	if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
