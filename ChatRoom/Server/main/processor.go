package main

import (
	"ChatRoom/Commen/message"
	"ChatRoom/Server/process2"
	"ChatRoom/Server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}
//根据客户端发送消息的种类的不同决定用哪个函数处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error)  {
	fmt.Println("测试mes=",mes)

	switch mes.Tyep {
	case message.LoginMesType://处理登入
	//创建一个userprocess实例
		up :=&process2.UserProcess{
			Conn:this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up:=&process2.UserProcess{
			Conn:this.Conn,
		}
		err =up.ServerProcessRegister(mes)//type:data
	case message.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息
		smsProcess :=&process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型比存在，无法处理。。")
	}
	return
}
func (this *Processor) process2 () (err error)  {
	//循环处理客户端发送的消息
	for {
		//用封装函数readPkg()读取数据包，返回message，err
		//创建结构体对象
		tf := &utils.Transfer{
			Conn:this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err!=nil{
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出。",err)
				return err
			}else {
				fmt.Println("readPkg err=",err)
				return err
			}
		}
		fmt.Println("mes=",mes)
		err = this.serverProcessMes(&mes)
		if err!=nil {
			return err
		}
	}
}
