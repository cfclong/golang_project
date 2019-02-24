package model
//用户结构体
type User struct {
	//确定字段信息
	//为序列化和反序列化成功，我们必须保证
	//用户信息的json字符串的key和结构体字段对应的字段的tag名一致
	UserID int `json:"userId"`
	UserPWD string `json:"userPWD"`
	UserName string `json:"userName"`
}

