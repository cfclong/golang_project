package process2

import (
	"ChatRoom/Commen/message"
	"ChatRoom/Server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
	
}
//转发消息的方法
func (this *SmsProcess) SendGroupMes(mes *message.Message)  {
	//遍历服务器端的	onlineUsers map[int]*UserProcess
	//将消息转发取出
	//取出mes的内容SmsMes
	var smsMes message.SmsMes
	err :=json.Unmarshal([]byte(mes.Data),&smsMes)
	if err!=nil{
		fmt.Println("json.Unmarshal err=",err)
		return
	}
	data,err:=json.Marshal(mes)
	if err!=nil{
		fmt.Println("json.Marshal err=",err)
		return
	}
	for id,up :=range userMgr.onlineUsers{
		//这里不要发消息给自己
		if id==smsMes.UserID{
			continue
		}
		this.SendMesToEachOnlineUser(data,up.Conn)
	}
}
func (this *SmsProcess) SendMesToEachOnlineUser(data []byte,conn net.Conn)  {
	//创建transfer实例，发送data
	tf :=utils.Transfer{
		Conn:conn,
	}
	err :=tf.WritePkg(data)
	if err!=nil{
		fmt.Println("转发消息失败 err=",err)
		return
	}
}