package process2

import (
	"fmt"
)

//定义全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//完成对UserMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//完成对onlineUsers添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserID] = up
}

//删除
func (this *UserMgr) DelOnlineUser(userID int) {
	delete(this.onlineUsers, userID)
}

//返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

//根据id返回对应值
func (this *UserMgr) GetOnlineById(userID int) (up *UserProcess, err error) {
	//从map取出一个值，自带检测方法
	up, ok := this.onlineUsers[userID]
	if !ok { //说明你要找的用户当前不在线
		err = fmt.Errorf("用户%d,不在线", userID)
		return
	}
	return
}
