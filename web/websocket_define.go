package web

import "github.com/kataras/neffos"

type UWBProjWebSocket struct {
	ws    *neffos.Server
	conns map[*neffos.Conn]string //map的value存储uid，用于区分用户
}

// SetUID 设置用户信息
func (socket *UWBProjWebSocket) SetUID(con *neffos.Conn, uid string) {
	//myWebSocket.conns[con] = uid
	socket.conns[con] = uid

}

// DelConn 移除连接
func (socket *UWBProjWebSocket) DelConn(c *neffos.Conn) {
	delete(socket.conns, c)

}

func NewSocket() *UWBProjWebSocket {
	ws := UWBProjWebSocket{
		conns: make(map[*neffos.Conn]string),
	}
	return &ws
}
