package service

import (
	"github.com/samuel/go-zookeeper/zk"
)

// Get函数从指定路径的叶子节点读取数据。
// 返回数据和错误
func Get(conn *zk.Conn, path string) (data string, err error) {
	defer conn.Close()
	b, _, err := conn.Get(path) // 获取指定路径的数据
	if err != nil {
		return "", err
	}
	return string(b), nil
}
