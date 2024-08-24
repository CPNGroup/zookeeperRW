package service

import (
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

// Set函数在指定路径的叶子节点写入数据。
func Set(conn *zk.Conn, path string, data string) error {
	defer conn.Close()
	var flags int32 = 0
	var acls = zk.WorldACL(zk.PermAll)

	// 先检查路径是否存在。如果不存在，则调用 createPath 函数创建路径。
	exists, _, err := conn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = createPath(conn, path, "", flags, acls)
		if err != nil {
			return err
		}
	}
	// 路径创建完成后，使用 conn.Set 方法在指定路径上设置数据。
	_, err = conn.Set(path, []byte(data), -1)
	if err != nil {
		return err
	}
	return nil
}

// createPath 函数递归地创建路径中的每个节点。
func createPath(conn *zk.Conn, path string, createdPath string, flags int32, acls []zk.ACL) error {
	// 将路径按照斜杠（/）分割成各个部分，并依次创建每个节点。
	parts := strings.Split(path, "/")

	// 每次创建节点前，检查节点是否已存在。
	// 如果不存在，则使用 conn.Create 方法创建节点，并打印创建的路径。
	// 递归更新已创建的路径，确保创建的是完整路径。
	// 创建路径中的每个节点
	for i := range parts {
		if parts[i] == "" {
			// 跳过空节点
			continue
		}
		// 构建当前节点路径
		currentPath := createdPath + "/" + parts[i]

		// 检查节点是否存在
		exists, _, err := conn.Exists(currentPath)
		if err != nil {
			return err
		}

		// 如果节点不存在，则创建节点
		if !exists {
			_, err := conn.Create(currentPath, []byte{}, flags, acls)
			if err != nil {
				return err
			}
		}

		// 更新已创建的路径
		createdPath = currentPath
	}

	return nil
}
