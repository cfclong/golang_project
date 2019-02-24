package process

import (
	"ChatRoom/Client/utils"
	"ChatRoom/Commen/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct{}

func (this *UserProcess) Register(userID int, userPWD string, UserName string) (err error) {
	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()
	//2.通过conn发送消息给服务器
	var mes message.Message
	mes.Tyep = message.RegisterMesType
	//3.创建一个loginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserID = userID
	registerMes.User.UserPWD = userPWD
	registerMes.User.UserName = UserName
	//4.将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//5.把data赋给mes.Data字段
	mes.Data = string(data)
	//6.将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Mashal err=", err)
		return
	}
	//创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	//发送给服务的
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送消息错误 err", err)
	}
	mes, err = tf.ReadPkg() //mes就是RegisterResMes
	if err != nil {
		fmt.Println("ReadPkg(conn) err=", err)
		return
	}
	//将mes的data部分反序列化为RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登入")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
func (this *UserProcess) Login(id int, pwd string) (err error) {
	//1.连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	defer conn.Close()
	//2.向服务器发送消息
	var mes message.Message
	mes.Tyep = message.LoginMesType
	//3.创建一个loginmes结构体
	var loginMes message.LoginMes
	loginMes.UserID = id
	loginMes.UserPWD = pwd
	//4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("jsonMarshal err=", err)
		return
	}
	//5.把data赋给mes.data字段
	mes.Data = string(data)
	//将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("jsonMarshal err=", err)
		return
	}
	//7.data就是要发送的消息
	//先取data长度->转成表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write() err=", err)
		return
	}
	//fmt.Printf("客户端发送消息长度=%d 内容%s",len(data),string(data))
	//return
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write err", err)
	}
	//time.Sleep(time.Second *20)
	//fmt.Println("休眠20秒。。")
	//创建transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}
	//将mes的data部分反序列化成loginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserID = id
		CurUser.UserStatus = message.UserOnline
		//fmt.Println("登入成功")
		//可以显示当前用户列表，遍历LoginResMes.UserID
		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UsersID {
			//可以让自己不显示在线
			if v == id {
				continue
			}
			fmt.Println("用户id:\t", v)
			//完成客户端的onlineUsers初始化
			user := &message.User{
				UserID:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Println("\n\n")
		//z这里再客户端启动一个协程
		//该协程保持与服务器的通讯，如果服务器有消息发送给客户端
		//则接受并显示在客户端终端
		go ServerProcessMes(conn)
		//显示登入成功后的菜单
		for {
			ShowMenu()
		}
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	} else {
		fmt.Println("其他错误。。。。")
	}
	return
}
