package process

import (
	"ChatRoom/Client/utils"
	"ChatRoom/Commen/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func ShowMenu() {
	fmt.Println("---------------登入成功---------------------")
	fmt.Println("---------------1.显示在线用户---------------")
	fmt.Println("---------------2.发送消息-------------------")
	fmt.Println("---------------3.消息列表-------------------")
	fmt.Println("---------------4.退出系统-------------------")
	fmt.Println("请选择（1-4）：")
	var key int
	var content string
	//实例化SmsProcess
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("显示在线用户。。")
		outputOnlineUser()
	case 2:
		fmt.Println("你想对大家说什么")
		fmt.Scanf("%s\n",&content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("消息列表。。")
	case 4:
		fmt.Println("您退出了系统。。")
		os.Exit(0)
	default:
		fmt.Println("输入错误。。")
	}
}

//和服务器保持通讯
func ServerProcessMes(conn net.Conn) {
	//创建transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端在等待服务端发送消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		//读取到消息,下一步处理
		switch mes.Tyep {
		case message.NotifyUserStatusMesType: //有人上线
			//1.取出.NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2.把这个用户消息，状态保存到客户map[int]User中
			updateUserStatus(&notifyUserStatusMes)
		//处理
		case message.SmsMesType://有人群发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务端返回了未知消息类型")
		}
		//fmt.Println("mes=",mes)
	}
}
