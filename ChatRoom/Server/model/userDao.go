package model

import (
	"ChatRoom/Commen/message"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//在服务启动后就初始化一个userDao实例
//把它做成全局变量，需要和redis操作时，直接使用
var (
	MyUserDao *UserDao
)
//定义UserDao结构体
//完成UserDao结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}
//工厂模式，创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao)  {
	userDao=&UserDao{
		pool:pool,
	}
	return
}
//UserDao提供的方法
//根据用户id返回一个User实例+err
func (this *UserDao) getUserById(conn redis.Conn,id int) (user *User,err error)  {
	//通过给定id去redis查询这个用户
	res,err :=redis.String(conn.Do("HGet","users",id))
	if err!=nil{
		//错误
		if err ==redis.ErrNil{
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	//把res反序列化成user实例
	err =json.Unmarshal([]byte(res),user)
	if err!=nil{
		fmt.Println("json.Unmarshal err=",err)
		return
	}
	return
}
//完成登入的校验
//1.login完成对用户的验证
//2.如果id和pwd都正确，则返回一个user实例
//3.如果id和pwd有误，则返回错误信息
func (this *UserDao) Login(userID int,userPWD string) (user *User,err error)  {
	//先从UserDao的连接池中取出一个链接
	conn :=this.pool.Get()
	defer conn.Close()
	user,err=this.getUserById(conn,userID)
	if err!=nil{
		return
	}
	//这时证明这个用户是获取到的
	if user.UserPWD !=userPWD{
		err=ERROR_USER_PWD
		return
	}
	return
}
//实现注册功能方法
func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao中的连接池中取出一个链接
	conn:=this.pool.Get()
	defer conn.Close()
	_,err=this.getUserById(conn,user.UserID)
	if err==nil{//没有错误表示数据库里存在了该账户，无法注册该用户，返回错误
		err =ERROR_USER_EXISTS
		return
	}
	//这是说明redis里面还没有，可以完成注册
	data,err:=json.Marshal(user)//序列化
	if err!=nil{
		return
	}
	//写入数据库
	_,err=conn.Do("HSet","users",user.UserID,string(data))
	if err!=nil{
		fmt.Println("保存注册用户出错 err=",err)
		return
	}
	return
}