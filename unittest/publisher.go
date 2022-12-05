package main

import (
	"fmt"
	"time"
)

func PublishMessage(msg string) {
	// 向 主题 topic/test 发送消息
	topic := "topic/test"
	QoS := byte(0)
	fmt.Println(msg)
	token := MQTTClient.Publish(topic, QoS, false, msg)
	token.Wait()
}
func DoTestPublish() {
	ConnectionMQTTServer("go_mqtt_client_publisher")
	var i int
	for {
		i = i + 1
		msg := fmt.Sprintf("Message %d", i)
		PublishMessage(msg)
		time.Sleep(time.Second)
	}
}
