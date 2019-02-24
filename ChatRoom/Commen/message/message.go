package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

//用户状态常量
const (
	UserOnline = iota
	UserOfline
	UserBusyStatus
)

type Message struct {
	Tyep string `json:"type"` //消息类型
	Data string `json:"data"` //消息内容
}
type LoginMes struct {
	UserID   int    `json:"userID"`
	UserPWD  string `json:"userPWD"`
	UserName string `json:"userName"`
}
type LoginResMes struct {
	Code    int `json:"code"`     //返回状态码
	UsersID []int                 //增加字段，保持用户id切片
	Error   string `json:"error"` //返回错误信息
}
type RegisterMes struct {
	User User `json:"user"`
}
type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码400表示该用户已经占有，200表示注册成功
	Error string `json:"error"` //返回错误信息
}

//为了配合服务器端推送用户状态的消息
type NotifyUserStatusMes struct {
	UserID int `json:"userID"`
	Status int `json:"status"`
}

//SmsMes机构体，发送消息
type SmsMes struct {
	Content string `json:"content"`
	User //继承
}
