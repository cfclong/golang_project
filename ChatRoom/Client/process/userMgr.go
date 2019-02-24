package process

import (
	"ChatRoom/Client/model"
	"ChatRoom/Commen/message"
	"fmt"
)

//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser //再用户登入成功后完成对CurUser的初始化
//在客户端显示当前在线用户
func outputOnlineUser() {
	//遍历onlineUsers
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id：\t", id)
	}
}

//处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	//适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserID]
	if !ok { //原来没有
		user = &message.User{
			UserID: notifyUserStatusMes.UserID,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserID] = user
	outputOnlineUser()
}
