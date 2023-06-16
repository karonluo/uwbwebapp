package web

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

type LoginInformation struct {
	GWAddress string `json:"gw_address"`
	GWId      string `json:"gw_id"`
	GWPort    string `json:"gw_port"`
	MAC       string `json:"mac"`
	URL       string `json:"url"`
}

type AuthorizedInformatin struct {
	Stage    string `json:"stage"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	TOKEN    string `json:"token"`
	Incoming string `json:"incoming"`
	Outgoing string `json:"outgoing"`
}

type HeartBeatInformatin struct {
	// http://auth_sever/ping/?gw_id=%s&sys_uptime=%lu&sys_memfree=%u&sys_load=%.2f&wifidog_uptime=%lu
	GWId          string `json:"gw_id"`
	SysUptime     string `json:"sys_uptime"`
	SysMemFree    string `json:"sys_memfree"`
	SysLoad       string `json:"sys_load"`
	WifiDogUptime string `json:"wifidog_uptime"`
}

// Login Process
func LoginScriptPathFragment(ctx iris.Context) {

	info := LoginInformation{
		GWAddress: ctx.FormValue("gw_address"),
		GWId:      ctx.FormValue("gw_id"),
		GWPort:    ctx.FormValue("gw_port"),
		MAC:       ctx.FormValue("mac"),
		URL:       ctx.FormValue("url"),
	}
	token := tools.SHA1(info.MAC)
	// TODO: 将本次登录动作的数据存入 Redis 进行记录。
	gwhttp := fmt.Sprintf("http://%s:%s/cgi-bin/auth?token=%s", info.GWAddress, info.GWPort, token)
	bytesInfo, _ := json.Marshal(info)
	strInfo := string(bytesInfo)
	rctx := context.Background()
	cache.RedisDatabase.Set(rctx, "wifidog_token_"+token, strInfo, time.Hour*24)
	fmt.Println(string(strInfo))
	// 这里跳转相当于跳过 需要访问验证的页面直接进入设备认证地址， 否则应该先跳到其他页面，用户确定后再转到设备认证地址。
	// TODO: 增加读取广告列表，随机抽取广告页面，广告页面中应该有上述地址，隐藏，直到N秒后，允许用户继续后开放访问所有的网站。
	ctx.Redirect(gwhttp, 301) // 301 == GET

	/*
		gw_address=%s&gw_port=%d&gw_id=%s&mac=%s&url=%s
	*/
}

// View 页面
func PortalScriptPathFragment(ctx iris.Context) {
	// 认证成功后的跳转
	ctx.HTML("Authorized success, thanks. <br />You can access any website now.")

}

// 返回信息
func MsgScriptPathFragment(ctx iris.Context) {
	ctx.Text("Msg")
}

// 认证
func AuthScriptPathFragment(ctx iris.Context) {

	// auth_server:/auth/auth.php?stage=%s&ip=%s&mac=%s&token=%s&incoming=%s&outgoing=%s
	// ctx.Text("Auth: 0") // 失败

	/*
		0 - AUTH_DENIED - User firewall users are deleted and the user removed.
		6 - AUTH_VALIDATION_FAILED - User email validation timeout has occured and user/firewall is deleted
		1 - AUTH_ALLOWED - User was valid, add firewall rules if not present
		5 - AUTH_VALIDATION - Permit user access to email to get validation email under default rules
		-1 - AUTH_ERROR - An error occurred during the validation process
	*/
	// 当 stage=login 时为第一次登录，以后都是计数器，代表保持链接。
	// TODO: 登录/计数/离线
	info := AuthorizedInformatin{
		Stage:    ctx.FormValue("stage"),
		IP:       ctx.FormValue("ip"),
		MAC:      ctx.FormValue("mac"),
		TOKEN:    ctx.FormValue("token"),
		Incoming: ctx.FormValue("incoming"),
		Outgoing: ctx.FormValue("outgoing"),
	}
	strInfo, _ := json.Marshal(info)
	fmt.Println(string(strInfo))
	// 从 cache 中找到相同的 token 证明登录成功
	rctx := context.Background()
	res, err := cache.RedisDatabase.Exists(rctx, "wifidog_token_"+info.TOKEN).Result()
	if err == nil && res == 1 {
		ctx.Text("Auth: 1") // 认证成功
	} else {
		ctx.Text("Auth: 0") // 认证失败
	}

}

// Ping Pong
func DevicePing(ctx iris.Context) {
	// GW 向 认证中心发送 设备 心跳信息。
	// http://auth_sever/ping/?gw_id=%s&sys_uptime=%lu&sys_memfree=%u&sys_load=%.2f&wifidog_uptime=%lu
	// TODO: 对设备进行管理，可以看到所有路由器的链接（在线离线）状态。

	// type HeartBeatInformatin struct {
	// 	// http://auth_sever/ping/?gw_id=%s&sys_uptime=%lu&sys_memfree=%u&sys_load=%.2f&wifidog_uptime=%lu
	// 	GWId          string `json:"gw_id"`
	// 	SysUptime     string `json:"sys_uptime"`
	// 	SysMemFree    string `json:"sys_memfree"`
	// 	SysLoad       string `json:"sys_load"`
	// 	WifiDogUptime string `json:"wifidog_uptime"`
	// }
	info := HeartBeatInformatin{
		GWId:          ctx.FormValue("gw_id"),
		SysUptime:     ctx.FormValue("sys_uptime"),
		SysMemFree:    ctx.FormValue("sys_memfree"),
		SysLoad:       ctx.FormValue("sys_load"),
		WifiDogUptime: ctx.FormValue("wifidog_uptime"),
	}
	strInfo, _ := json.Marshal(info)
	fmt.Println(string(strInfo))
	ctx.Text("Pong")
}

func NoDogSplashAuthorizedService(ctx iris.Context) {
	/*
		- `clientip`：客户端设备的 IP 地址。
		- `clientmac`：客户端设备的 MAC 地址。
		- `gatewayname`：Nodogsplash 网关的名称。
		- `tok`：一个唯一的身份验证令牌，用于识别客户端。
		- `redir`：客户端在认证成功后应重定向的原始请求 URL。
	*/
	clientIP := ctx.FormValue("ip")
	clientMAC := ctx.FormValue("mac")
	gatewayName := ctx.FormValue("gn")
	tok := ctx.FormValue("tok")
	redir := ctx.FormValue("url")
	// http://$gatewayname:2050/nodogsplash_auth?token=$tok&redir=$redir
	// 记录 clientIP\clientMAC
	fmt.Println(clientIP)    // 访问者IP地址
	fmt.Println(clientMAC)   // 访问者客户端MAC地址
	fmt.Println(redir)       // 访问者访问网站地址
	fmt.Println(gatewayName) // 证明 路由器ID, 比如泰国某个餐馆

	// 重定向边缘软路由器NodogSplash认证地址
	// ctx.Redirect(fmt.Sprintf("http://%s:2050/nodogsplash_auth?tok=%s&redir=%s", gatewayName, tok, redir), 302)
	// ctx.Text(fmt.Sprintf("http://%s:2050/nodogsplash_auth/?token=%s&redir=%s", tools.GetGatewayIP(clientIP), tok, redir))
	ctx.Redirect(fmt.Sprintf("http://%s:2050/nodogsplash_auth/?token=%s&redir=%s", tools.GetGatewayIP(clientIP), tok, redir),
		302)
	// http://路由器IP:2050/nodogsplash_auth/?tok=12345&mac=AA:BB:CC:DD:EE:FF&ip=192.168.1.2
	// http: //192.168.8.1:2050/nodogsplash_auth/?token=56183885&redir=http%3a%2f%2f192.168.8.1%2f
}
