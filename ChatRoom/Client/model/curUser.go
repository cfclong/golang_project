package model

import (
	"ChatRoom/Commen/message"
	"net"
)

//再客户端很多地方使用，定义全局
type CurUser struct {
	Conn net.Conn
	message.User
}
