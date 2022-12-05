package main

import (
	"fmt"
	"time"
)

// 订阅函数
func Subscribe() {
	topic := "topic/test" // 设置订阅 topic/test 主题
	QoS := byte(0)        // QoS 设置服务质量等级， 请参阅 QoS 相关介绍
	// 订阅 topic/test 主题、QoS、获取消息的函数
	token := MQTTClient.Subscribe(topic, QoS, ReceiveMessage)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\r\n", topic)
	for {
		time.Sleep(time.Second)
	}
}

func DoTestSubscribe() {
	ConnectionMQTTServer("go_mqtt_client_subscriber")
	Subscribe()
}
