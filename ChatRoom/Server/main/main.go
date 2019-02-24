package main

import (
	"ChatRoom/Server/model"
	"fmt"
	"net"
	"time"
)

//func readPkg(conn net.Conn) (mes message.Message, err error) {
//	fmt.Println("读取客户端发送的数据")
//	buf := make([]byte, 8096)
//	//conn.Read在conn没有关闭的时候才会阻塞，如果客户端关闭了conn则不会
//	n, err := conn.Read(buf[0:4])
//	if n != 4 || err != nil {
//		//fmt.Println("conn.Read err=",err)
//		return
//	}
//	//fmt.Println("读取到buf=",buf[0:4])
//	//根据buf[:4]转成一个uint32类型
//	var pkgLen uint32
//	pkgLen = binary.BigEndian.Uint32(buf[:4])
//	//根据pkLen读取消息内容
//	n2, err := conn.Read(buf[:pkgLen])
//	if n2 != int(pkgLen) || err != nil {
//		fmt.Println("conn.Read err=", err)
//		return
//	}
//	//反序列化
//	err = json.Unmarshal(buf[:pkgLen], &mes)
//	if err != nil {
//		fmt.Println("json.unmarshal err", err)
//	}
//	return
//}
//
//func writePkg(conn net.Conn,data []byte) (err error) {
//	//先发送长度给客户端
//	var pkgLen uint32
//	pkgLen =uint32(len(data))
//	var buf [4]byte
//	binary.BigEndian.PutUint32(buf[:4], pkgLen)
//	//发送长度
//	n,err :=conn.Write(buf[:4])
//	if n!=4||err !=nil{
//		fmt.Println("conn.Write err=",err)
//		return
//	}
//	//发送data本身
//	n,err =conn.Write(data)
//	if n!=int(pkgLen) ||err!=nil {
//		fmt.Println("conn.Write err=",err)
//		return
//	}
//	return
//}

//func process(conn net.Conn) {
//	defer conn.Close()
//
//	//循环处理客户端发送的消息
//	for {
//		//用封装函数readPkg()读取数据包，返回message，err
//		mes, err := readPkg(conn)
//		if err!=nil{
//			if err == io.EOF {
//				fmt.Println("客户端退出，服务端也退出。")
//				return
//			}else {
//				fmt.Println("readPkg err=",err)
//				return
//			}
//		}
//		fmt.Println("mes=",mes)
//		err = serverProcessMes(conn,&mes)
//		if err!=nil {
//			return
//		}
//	}
//}
////处理登入请求的函数
//func serverProcessLogin(conn net.Conn,mes *message.Message) (err error) {
//	//先从mes中取出mes.Data，并直接反序列化成LoginMes
//	var loginMes message.LoginMes
//	err =json.Unmarshal([]byte(mes.Data),&loginMes)
//	if err !=nil {
//		fmt.Println("json.Unmarshal fail err=",err)
//		return
//	}
//	//1.先声明一个resMes
//	var resMes message.Message
//	resMes.Tyep = message.LoginMesType
//	//2.声明一个LoginResMes
//	var loginResMes message.LoginResMes
//	//暂时设定用户id=100,密码=123456为合法，否则不合法
//	if loginMes.UserID ==100 && loginMes.UserPWD=="123456" {
//		//合法，返回状态码200
//		loginResMes.Code =200
//	}else {
//		//不合法，返回状态码500
//		loginResMes.Code=500
//		loginResMes.Error= "该用户不存在，请注册后使用"
//	}
//	//3.将loginResMes序列化
//	data,err := json.Marshal(loginResMes)
//	if err !=nil{
//		fmt.Println("json.Marshal fail err=",err)
//		return
//	}
//	//4.将data赋给resMes
//	resMes.Data =string(data)
//	//5.对resMes进行序列化，准备发送
//	data,err =json.Marshal(resMes)
//	if err!= nil {
//		fmt.Println("json.Marshal err=",err)
//		return
//	}
//	//6.发送data，封装函数writePkg()
//	err =writePkg(conn, data)
//	return
//}
//func serverProcessMes(conn net.Conn,mes *message.Message) (err error)  {
//	switch mes.Tyep {
//	case message.LoginMesType://处理登入
//		err = serverProcessLogin(conn, mes)
//	case message.RegisterMesType:
//		//处理注册
//	default:
//		fmt.Println("消息类型比存在，无法处理。。")
//	}
//	return
//}

//处理和客户端的通讯
func process(conn net.Conn) {
	defer conn.Close()
	//调用总控，创建一个
	processor := &Processor{
		Conn: conn,
	}
	err :=processor.process2()
	if err!=nil{
		fmt.Println("客户端服务端协程通讯错误err",err)
		return
	}
}
//UserDao初始化函数
func initUserDao() {
//要先初始化initpool,再初始化initUserDao
	model.MyUserDao =model.NewUserDao(pool)
}
func main() {
	//服务启动就初始化
	initPool("localhost:6379",16,0,300*time.Second)
	initUserDao()
	//提示信息
	fmt.Println("服务器在监听8889端口")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close() //释放监听资源
	if err != nil {
		fmt.Println("监听失败。。")
		return
	}
	//等待客户连接。
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//启动一个协
		go process(conn)
	}
}
