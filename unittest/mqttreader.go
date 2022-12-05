package main

import (
	"flag"
	"fmt"
	"time"
	"uwbwebapp/conf"

	emqx "github.com/eclipse/paho.mqtt.golang"
)

func DoTestMQTTReader() {

	mqttClientId := flag.String("clientid", "go_mqtt_client_reader", "MQTT 客户端编号")
	flag.Parse()
	conf.LoadWebConfig("./conf/WebConfig.json")
	opts := emqx.NewClientOptions()
	dsn := fmt.Sprintf("tcp://%s:%d", conf.WebConfiguration.MQTTServerConf.Broker, conf.WebConfiguration.MQTTServerConf.Port)
	opts.AddBroker(dsn)
	opts.SetClientID(*mqttClientId)
	opts.SetUsername(conf.WebConfiguration.MQTTServerConf.User)
	opts.SetPassword(conf.WebConfiguration.MQTTServerConf.Password)
	MQTTClient = emqx.NewClient(opts)
	if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub()

}

func sub() {
	topic := "topic/test"
	token := MQTTClient.Subscribe(topic, 1, ReceiveMessage)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\r\n", topic)
	for {
		time.Sleep(time.Second)
	}
}
