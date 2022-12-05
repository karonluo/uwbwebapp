package message

import (
	"fmt"
	"uwbwebapp/conf"

	emqx "github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient emqx.Client

func InitMQTTClient() {
	if MQTTClient == nil {
		fmt.Print("初始化MQTT链接")
		MQTTMessageHandler := OnReceiveMessage

		opts := emqx.NewClientOptions()
		dsn := fmt.Sprintf("tcp://%s:%d", conf.WebConfiguration.MQTTServerConf.Broker, conf.WebConfiguration.MQTTServerConf.Port)
		opts.AddBroker(dsn)
		opts.SetClientID("go_mqtt_client")
		opts.SetUsername(conf.WebConfiguration.MQTTServerConf.User)
		opts.SetPassword(conf.WebConfiguration.MQTTServerConf.Password)
		opts.SetDefaultPublishHandler(MQTTMessageHandler)

		MQTTClient = emqx.NewClient(opts)

		if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println("......失败")
			panic(token.Error())

		} else {
			fmt.Println("......成功")
		}
	} else {
		if !MQTTClient.IsConnected() {
			if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
				fmt.Println("......失败")
				panic(token.Error())
			} else {
				fmt.Println("......成功")
			}

		}
	}
}

func OnLostConnection(client emqx.Client) {

}

func OnReceiveMessage(client emqx.Client, msg emqx.Message) {
	fmt.Println("ReceiveMessage")
	if MQTTClient != nil {
		fmt.Printf("Received message: %s from topic: %s\r\n", msg.Payload(), msg.Topic())
	}
}
