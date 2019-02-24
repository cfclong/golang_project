package process

import (
	"ChatRoom/Client/utils"
	"ChatRoom/Commen/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个Mes
	var mes message.Message
	mes.Tyep = message.SmsMesType
	//2.创建SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserID = CurUser.UserID
	smsMes.UserStatus = CurUser.UserStatus
	//3.序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}
	mes.Data = string(data)
	//4.对Mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}
	//5.将mes发送给服务器
	tf :=utils.Transfer{
		Conn:CurUser.Conn,
	}
	//6.发送
	err =tf.WritePkg(data)
	if err!=nil{
		fmt.Println("SendGroupMes err=",err.Error())
		return
	}
	return
}
