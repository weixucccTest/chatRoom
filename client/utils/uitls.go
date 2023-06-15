package utils

import (
	"common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 这里讲这些方法关联到结构体中
type Transfer struct {
	// 分析它应该有哪些字段
	Conn net.Conn
	// 这是传输时使用的缓冲
	Buf [8064]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {

	// buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据。。。")
	// conn.read 在conn没有被关闭的情况下，才会组赛
	// 如果客户端关闭了 conn，则不会阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		// err = errors.New("read pkg header error")
		return
	}

	// 根据buf[:4] 转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	// 根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		// err = errors.New("read pkg body error")
		return
	}

	// 把pkglen反序列化成 -》 message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {

	// 先发送一个长度给对方，
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.write(butes) fail, err = ", err)
		return
	}

	// 发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) fail, err = ", err)
		return
	}
	return
}
