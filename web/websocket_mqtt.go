package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	emqx "github.com/eclipse/paho.mqtt.golang"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos/gorilla"
)

// Begin UWB Tag Terminal Service
var MQTTUWBTagClient emqx.Client

func CheckPointInPolygon() bool {
	return true
}

// 获取 MQTT Server 关于 UWB 终端标签的信息
func ReceiveUWBTagMessage(client emqx.Client, msg emqx.Message) {
	infor, err := entities.UnmarshalUWBTerminalMQTTInformation(msg.Payload()) // 获取订阅的 MQTT Server 原始的 UWB Terminal 信息
	var sendData []byte
	if err == nil {
		cons := GlobalUWBWebSocket.conns // 获取所有 Web Socket 链接，方便广播到所有 Web Socket 客户端
		for con, _ := range cons {
			// con.Socket().WriteText(msg.Payload(), 5*time.Second)
			//TODO: 这里需要判断是否符合场地的游泳池区域信息，如果在场地游泳池区域内，才进行累计距离。
			cacheInfor, er := cache.GetUWBTerminalTagFromRedis(infor.Properties.SwimmerId) // 获取缓存中的 UWB Terminal 信息
			if er == nil {
				// 当没有产生任何错误时（注意：参阅自定义错误列表，在 cache.GetUWBTerminalTagFromRedis 中）
				if CheckPointInPolygon() { // 目前默认在泳池内
					infor.SumDistance = cacheInfor.SumDistance + infor.Distance // 将原有存入缓存的 总里程 + 目前收集到的终端距离，累计新的总里程
				} else {
					infor.SumDistance = cacheInfor.SumDistance // 如果不在泳池内，则不该表总里程数，避免增加额外的里程，无法表达真实的游泳运动情况（注意：有一定误差）
				}
				cache.SetUWBTerminalTagToRedis(&infor) // 将处理后的 UWB Terminal 信息存入缓存中
				sendData, _ = json.Marshal(infor)      // 使用 JSON 解析，用于广播到已链接到服务的 WEB SOCKET 客户端
				con.Socket().WriteText(sendData, 5*time.Second)
			} else if er.Error() == "the UWB terminal was closed" {
				// 如果收到错误信息是 UWB Terminal 缓存被关闭（证明该游泳者被显性的签出了）则发送清空 WEB Socket 客户端 UWB 标签信息的命令
				clearWSClientUWBInforCommand := make(map[string]string)
				clearWSClientUWBInforCommand["swimmerId"] = infor.Properties.SwimmerId
				clearWSClientUWBInforCommand["devEui"] = infor.DevEui
				clearWSClientUWBInforCommand["action"] = "clear"
				clearWSClientUWBInforCommand["swimmerDisplayName"] = infor.Properties.SwimmerDisplayName
				clearWSClientUWBInforCommand["message"] = ""
				rctx := context.Background()
				cache.RedisDatabase.Set(rctx, "UWBTag_"+infor.Properties.SwimmerId, "clear", 24*time.Hour) // 向 cache 发送 clear 值
				sendData, _ = json.Marshal(clearWSClientUWBInforCommand)
				con.Socket().WriteText(sendData, 5*time.Second)
			} else if er.Error() != "the UWB terminal was cleared on web socket client" {

				// 其他错误则将当前距离作为初始总里程
				infor.SumDistance = infor.Distance
				tools.ProcessError(`web.ReceiveUWBTagMessage`,
					`oInfor, er := cache.GetUWBTerminalTagFromRedis(infor.Properties.SwimmerDisplayName, infor.Properties.SwimmerId)`, er)
				cache.SetUWBTerminalTagToRedis(&infor)
				sendData, _ = json.Marshal(infor)
				con.Socket().WriteText(sendData, 5*time.Second)

			}
			// the UWB terminal was cleared on web socket client 代表已经从 web socket 客户端清楚了标签信息，不用做任何操作。
		}
	} else {
		fmt.Println(err.Error())
		tools.ProcessError(`web.ReceiveUWBTagMessage`,
			`infor, err := entities.UnmarshalUWBTerminalMQTTInformation(msg.Payload())`, err)
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

// End UWB Tag Terminal Service

// Begin Web Socket Define
var GlobalUWBWebSocket *UWBProjWebSocket

func MQTTWebSocketHandle(app *iris.Application, path string) {
	// 注意: 设置 gorilla.Upgrader 用于 ws:// 协议跨域
	GlobalUWBWebSocket.ws = websocket.New(gorilla.Upgrader(
		gorillaWs.Upgrader{
			CheckOrigin: func(*http.Request) bool {
				return true
			}}),
		websocket.Events{
			websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
				var sendData []byte
				log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
				command := string(msg.Body)
				if strings.Contains(command, "GetSiteSwimmerReport:") {
					// 获取泳池人员统计报告
					report, err := biz.SiteSwimmerReport(strings.Split(command, ":")[1])
					if err == nil {
						sendData, _ = json.Marshal(report)
					} else {
						sendData = []byte(err.Error())
					}
				} else if command == "PING" {
					sendData = []byte("PONG")

				} else if strings.Contains(command, "GetSwimmersSumDistenceOrder:") {
					// 获取游泳距离前五排名
					report, err := biz.GetSwimmersSumDistenceOrder(strings.Split(command, ":")[1])
					if err == nil {
						sendData, _ = json.Marshal(report)
					} else {
						sendData = []byte(err.Error())
					}

				} else {
					sendData = []byte("NoData")
				}
				nsConn.Conn.Socket().WriteText(sendData, 10*time.Second)
				return nil
			},
		})

	// 连接时设置用户信息
	GlobalUWBWebSocket.ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		ctx := websocket.GetContext(c)
		uid := ctx.URLParam("uid")
		GlobalUWBWebSocket.SetUID(c, uid)
		return nil
	}
	GlobalUWBWebSocket.ws.OnDisconnect = func(c *websocket.Conn) {
		GlobalUWBWebSocket.DelConn(c)
		log.Printf("[%s] Disconnected from server", c.ID())
	}
	GlobalUWBWebSocket.ws.OnUpgradeError = func(err error) {
		log.Printf("Upgrade Error: %v", err)
	}
	app.Get(path, websocket.Handler(GlobalUWBWebSocket.ws))
}
func BoardcastToWebSocketClient() {
	uwbPlatConf := conf.WebConfiguration.UWBDevicePlatformConf
	// topic := "uwb/3/2/2/upload/+" // 设置订阅 topic/test 主题
	topic := fmt.Sprintf("uwb/%d/%d/%d/upload/+", uwbPlatConf.OrganizationId, uwbPlatConf.ApplicationId, uwbPlatConf.DefaultRegionId)
	QoS := byte(0) // QoS 设置服务质量等级， 请参阅 QoS 相关介绍
	// 订阅 topic/test 主题、QoS、获取消息的函数
	token := MQTTUWBTagClient.Subscribe(topic, QoS, ReceiveUWBTagMessage)
	token.Wait()

}

// End Web Socket Define
