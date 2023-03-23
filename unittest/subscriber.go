package main

import (
	"fmt"
	"time"

	emqx "github.com/eclipse/paho.mqtt.golang"
)

var MQTTUWBTagClient emqx.Client

func ReceiveUWBTagMessage(client emqx.Client, msg emqx.Message) {
	fmt.Println(client.IsConnected())
	fmt.Printf("Received message: %s from topic: %s\r\n", string(msg.Payload()), msg.Topic())
}

func ConnectionUWBTagMQTTServer(clientId string) {
	opts := emqx.NewClientOptions()
	dsn := fmt.Sprintf("tcp://%s:%d", "mqtt.4g.zah.dhwork.cn", 18883)
	opts.AddBroker(dsn)
	opts.SetClientID(clientId)
	opts.SetUsername("dhza_mqtt_tmp1")
	opts.SetPassword("doon1eiMaiBeChus")

	// dsn := fmt.Sprintf("tcp://%s:%d", "219.142.82.97", 1883)
	// opts.AddBroker(dsn)
	// opts.SetClientID(clientId)
	// opts.SetUsername("mqtt_adm")
	// opts.SetPassword("eeN*eizai8ah")

	MQTTUWBTagClient = emqx.NewClient(opts)
	if token := MQTTUWBTagClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// client.connect('219.142.82.97', 1883, 60)
// client.username_pw_set('mqtt_adm', 'eeN*eizai8ah')

// 订阅函数
func Subscribe() {

	// fmt.Sprintf("uwb/%s/%s/%s/upload/+", uwbPlatConf.OrganizationId, uwbPlatConf.ApplicationId, uwbPlatConf.ApplicationId)
	topic := "uwb/3/2/2/upload/+" // 设置订阅 topic/test 主题
	QoS := byte(0)                // QoS 设置服务质量等级， 请参阅 QoS 相关介绍
	// 订阅 topic/test 主题、QoS、获取消息的函数
	token := MQTTUWBTagClient.Subscribe(topic, QoS, ReceiveUWBTagMessage)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\r\n", topic)
	for {
		time.Sleep(500 * time.Millisecond)
	}
}

func DoTestSubscribe() {
	ConnectionUWBTagMQTTServer("go_mqtt_client_subscriber")
	Subscribe()
}

// func main() {
// 	DoTestSubscribe()
// }
