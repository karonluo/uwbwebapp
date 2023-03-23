package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
	gorillaWs "github.com/gorilla/websocket"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos/gorilla"
)

var GlobalAlertWebSocket *UWBProjWebSocket

// Begin Web Socket Define
func AlertWebSocketHandle(app *iris.Application, path string) {
	// 注意: 设置 gorilla.Upgrader 用于 ws:// 协议跨域
	GlobalAlertWebSocket.ws = websocket.New(gorilla.Upgrader(
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
					report, err := biz.SiteSwimmerReport(strings.Split(command, ":")[1])
					if err == nil {
						sendData, _ = json.Marshal(report)
					} else {
						sendData = []byte(err.Error())
					}
				} else if command == "PING" {
					sendData = []byte("PONG")

				} else if strings.Contains(command, "ClearAlert:") {
					//ClearAlret:swimmerId:userDisplayName 清空命令需要提供游泳者编号和发布命令的用户姓名
					suInfor := strings.Split(command, ":")
					swimmerId := suInfor[1]
					userDisplayName := suInfor[2]
					cache.ClearAlertInformationFromRedis(swimmerId) // 清除告警信息
					var calendar entities.SwimmerCalendar
					calendar.SwimmerID = swimmerId
					calendar.ExitDateTime = time.Now()
					calendar.Modifier = userDisplayName
					biz.SwimmerExitFromSite(&calendar) // 将游泳者签出，避免继续产生告警信息
					sendData = []byte(fmt.Sprintf("Cleared Alert: %s", swimmerId))
				} else {
					sendData = []byte("NoData")
				}
				nsConn.Conn.Socket().WriteText(sendData, 10*time.Second)
				return nil
			},
		})

	// 连接时设置用户信息
	GlobalAlertWebSocket.ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		ctx := websocket.GetContext(c)
		uid := ctx.URLParam("uid")
		GlobalAlertWebSocket.SetUID(c, uid)
		return nil
	}
	GlobalAlertWebSocket.ws.OnDisconnect = func(c *websocket.Conn) {
		GlobalAlertWebSocket.DelConn(c)
		log.Printf("[%s] Disconnected from server", c.ID())
	}
	GlobalAlertWebSocket.ws.OnUpgradeError = func(err error) {
		log.Printf("Upgrade Error: %v", err)
	}
	app.Get(path, websocket.Handler(GlobalAlertWebSocket.ws))
}
func BoardcastToAlertWebSocketClient() {
	for {
		time.Sleep(1 * time.Second)
		cons := GlobalAlertWebSocket.conns
		for conn, _ := range cons {
			infor := entities.AlertInformation{Message: "Hello,World!!!", SwimmerId: uuid.NewString()}
			sendData, _ := json.Marshal(infor)
			conn.Socket().WriteText(sendData, 5*time.Second)
		}
	}
}

// End Web Socket Define
