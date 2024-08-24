package conn

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// 接受一个 Zookeeper 服务器的地址+端口（socket）
// 并返回一个 Zookeeper 连接（*zk.Conn）。
// 如果连接成功，返回连接对象，并打印 "连接成功！"
// 如果失败，打印错误信息，并返回 nil。
func Conn(socket string) (*zk.Conn, error) {
	var hosts = []string{socket}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println("连接失败", err)
		return nil, err
	} else {
		fmt.Println("连接成功！")
		return conn, nil
	}
}
