package main

import (
	"ChatRoom/Client/process"
	"fmt"
	"os"
)

//用户id和密码
var userID int
var userPWD string
var userName string
func main() {
	//接受用户选择
	var key int
	//var loop = true
	for  {
		fmt.Println("---------------欢迎登入多人聊天系统---------------")
		fmt.Println("\t\t\t 1  登入聊天室")
		fmt.Println("\t\t\t 2  注册用户")
		fmt.Println("\t\t\t 3  退出系统")
		fmt.Println("\t\t\t    请输入（1-3）")
		fmt.Scanf("%d\n",&key)
		switch key {
		case 1:
			fmt.Println("登入聊天室")
			fmt.Println("输入用户id：")
			fmt.Scanf("%d\n",&userID)
			fmt.Println("输入用户密码：")
			fmt.Scanf("%s\n",&userPWD)
			//完成登入
			//1.创建实例
			up:=process.UserProcess{}
			up.Login(userID,userPWD)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("输入用户ID")
			fmt.Scanf("%d\n",&userID)
			fmt.Println("输入用户密码")
			fmt.Scanf("%s\n",&userPWD)
			fmt.Println("输入用户名")
			fmt.Scanf("%s\n",&userName)
			//2.调用UserProcess完成注册
			up := &process.UserProcess{}
			up.Register(userID,userPWD,userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误，重新输入")
		}
	}
}
