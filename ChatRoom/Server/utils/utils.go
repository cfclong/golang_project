package utils

import (
	"ChatRoom/Commen/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//将这些方法关联到结构体
type Transfer struct {
	Conn net.Conn
	Buf  [4096]byte //传输时缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送的数据")
	//buf := make([]byte, 8096)
	//conn.Read在conn没有关闭的时候才会阻塞，如果客户端关闭了conn则不会
	n, err := this.Conn.Read(this.Buf[0:4])
	if n != 4 || err != nil {
		//fmt.Println("conn.Read err=",err)
		return
	}
	//fmt.Println("读取到buf=",buf[0:4])
	//根据buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	//根据pkLen读取消息内容
	n2, err := this.Conn.Read(this.Buf[:pkgLen])
	if n2 != int(pkgLen) || err != nil {
		fmt.Println("conn.Read err=", err)
		return
	}
	//反序列化
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.unmarshal err", err)
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送长度给客户端
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	return
}
