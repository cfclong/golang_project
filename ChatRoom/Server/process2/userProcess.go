package process2

import (
	"ChatRoom/Commen/message"
	"ChatRoom/Server/model"
	"ChatRoom/Server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加字段表示Conn是那个用户的
	UserID int
}
//通知所以在线用户的方法
func (this *UserProcess) NotifyOtherOnlineUser(userId int)  {
	//遍历OnlineUsers,然后一个一个发送NotifyUserStatusMes
	for id,up :=range userMgr.onlineUsers{
		//过滤自己
		if id ==userId{
			continue
		}
		//开始通知，方法封装
		up.NotifyMeOnline(userId)
	}
}
func (this *UserProcess) NotifyMeOnline(userId int)  {
	//组装NotifyUserStatusMes
	var mes message.Message
	mes.Tyep=message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID=userId
	notifyUserStatusMes.Status=message.UserOnline
	//将notifyUserStatusMes序列化
	data,err:=json.Marshal(notifyUserStatusMes)
	if err!=nil{
		fmt.Println("json.Marshal err=",err)
		return
	}
	//将序列化后的notifyUserStatusMes赋给mes.Data
	mes.Data=string(data)
	//对mes序列化发送
	data,err=json.Marshal(mes)
	if err!=nil{
		fmt.Println("json.Marshal err=",err)
		return
	}
	//发送，创建Transfer实例
	tf:=&utils.Transfer{
		Conn:this.Conn,
	}
	err=tf.WritePkg(data)
	if err!=nil{
		fmt.Println("NotifyMeOnline err",err)
		return
	}
}
//处理登入请求的方法
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//先从mes中取出mes.Data，并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//1.先声明一个resMes
	var resMes message.Message
	resMes.Tyep = message.LoginMesType
	//2.声明一个LoginResMes
	var loginResMes message.LoginResMes
	//到redis数据库完成验证
	//使用model.MyUserDao到redis完成验证
	user, err := model.MyUserDao.Login(loginMes.UserID, loginMes.UserPWD)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
			fmt.Println("500")
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
			fmt.Println("403")
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
			fmt.Println("505")
		}
	} else {
		loginResMes.Code = 200
		//把登入成功的放到userMgr中，将登入成功的userid放入this
		this.UserID =loginMes.UserID
		userMgr.AddOnlineUser(this)
		//通知其他在线用户，我上线了
		this.NotifyOtherOnlineUser(loginMes.UserID)
		//将当前用户id放到loginResMes.UserId
		//遍历userMgr.OnlineUsers
		for id,_:=range userMgr.onlineUsers{
			loginResMes.UsersID=append(loginResMes.UsersID,id)
		}
		fmt.Println(user, "登入成功")
	}
	////暂时设定用户id=100,密码=123456为合法，否则不合法
	//if loginMes.UserID ==100 && loginMes.UserPWD=="123456" {
	//	//合法，返回状态码200
	//	loginResMes.Code =200
	//}else {
	//	//不合法，返回状态码500
	//	loginResMes.Code=500
	//	loginResMes.Error= "该用户不存在，请注册后使用"
	//}
	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//4.将data赋给resMes
	resMes.Data = string(data)
	//5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//6.发送data，封装函数writePkg()
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1.先从mes中取出mes.Data，并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	//1.先声明一个resMes
	var resMes message.Message
	resMes.Tyep = message.RegisterResMesType
	var registerResMes message.RegisterResMes
	//用到redis数据库去完成注册
	//1.使用model.MyUserDao到redis验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误。。"
		}
	} else {
		registerResMes.Code=200
	}
	data,err :=json.Marshal(registerResMes)
	if err!=nil{
		fmt.Println("json.Marshal err",err)
		return
	}

	//4.将data赋给resMes
	resMes.Data=string(data)
	//5.对resMes序列化，准备发送
	data,err =json.Marshal(resMes)
	if err!=nil{
		fmt.Println("json.Marshal err=",err)
		return
	}
	//6.发送data，将其封装到writPkg函数
	//因为使用了分层模式（mvc），我们先创建一个transfer实例，然后读取
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	err =tf.WritePkg(data)
	return
}
