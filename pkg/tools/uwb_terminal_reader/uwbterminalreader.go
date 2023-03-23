package main

import (
	"fmt"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	emqx "github.com/eclipse/paho.mqtt.golang"
)

var MQTTUWBTagClient emqx.Client

func ReceiveUWBTagMessage(client emqx.Client, msg emqx.Message) {
	infor, err := entities.UnmarshalUWBTerminalMQTTInformation(msg.Payload())
	if err == nil {
		fmt.Printf("标签编号: %s, 距离上次定位差：%f, 当前位置: X=%f, Y=%f\r\n------------------------------------------------\r\n",
			infor.DevEui, infor.Distance, infor.X, infor.Y)
		for _, station := range infor.StationInfos {
			fmt.Printf("基站编号: %s, 基站离标签距离：%f\r\n", station.DevEui, station.Distance)
		}
		fmt.Println("------------------------------------------------")
	} else {
		fmt.Println(err.Error())
	}
}

func ConnectionUWBTagMQTTServer(clientId string) {

	opts := emqx.NewClientOptions()
	mqttConf := conf.WebConfiguration.MQTTServerConf
	dsn := fmt.Sprintf("tcp://%s:%d", mqttConf.Broker, mqttConf.Port)
	opts.AddBroker(dsn)
	opts.SetClientID(clientId)
	opts.SetUsername(mqttConf.User)
	opts.SetPassword(mqttConf.Password)

	MQTTUWBTagClient = emqx.NewClient(opts)
	if token := MQTTUWBTagClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Subscribe() {
	uwbPlatConf := conf.WebConfiguration.UWBDevicePlatformConf
	// topic := "uwb/3/2/2/upload/+" // 设置订阅 topic/test 主题
	topic := fmt.Sprintf("uwb/%d/%d/%d/upload/+", uwbPlatConf.OrganizationId, uwbPlatConf.ApplicationId, uwbPlatConf.ApplicationId)
	QoS := byte(0) // QoS 设置服务质量等级， 请参阅 QoS 相关介绍
	// 订阅 topic/test 主题、QoS、获取消息的函数
	token := MQTTUWBTagClient.Subscribe(topic, QoS, ReceiveUWBTagMessage)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\r\n", topic)
	for {
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	conf.LoadWebConfig("../conf/WebConfig.json")
	dao.InitDatabase()

	if cache.InitRedisDatabase() {
		ConnectionUWBTagMQTTServer("go_mqtt_client_subscriber")
		fmt.Println("UWB 终端标签数据读取器初始化成功。")
		Subscribe()

		// rctx := context.Background()
		// for {

		// 	strCmd := cache.RedisDatabase.BLPop(rctx, time.Hour*24, "SystemLog")
		// 	if strCmd.Err() != nil {
		// 		fmt.Println(strCmd.Err().Error())
		// 	} else {
		// 		var log entities.SystemLog
		// 		strLog := strCmd.Val()[1]
		// 		err := json.Unmarshal([]byte(strLog), &log)
		// 		if err != nil {
		// 			fmt.Println(err.Error())
		// 		} else {
		// 			fmt.Println(log.String())
		// 			// 暂时使用数据库方式记录。
		// 			dao.WriteSystemLogToDB(&log)
		// 		}
		// 	}
		// }
	}
}
