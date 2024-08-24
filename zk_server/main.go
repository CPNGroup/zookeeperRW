package main

import (
	"context"
	"log"
	"net"
	"pro/conn"

	pb "pro/pb"
	"pro/service"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var port, socket = ReadConfig()

type server struct {
	pb.UnimplementedZkServiceServer
} //服务对象

// 实现服务的接口Get  在proto中定义的所有服务都是接口
func (s *server) Get(ctx context.Context, path *pb.Path) (*pb.Message, error) {
	var conn, err = conn.Conn(socket)
	if err != nil {
		return &pb.Message{Err: err.Error()}, err
	} else {
		data, err := service.Get(conn, path.Path)
		if err != nil {
			return &pb.Message{Err: err.Error()}, err
		} else {
			return &pb.Message{Data: data}, nil
		}
	}

}

// 实现服务的接口Set
func (s *server) Set(ctx context.Context, pathandname *pb.PathAndData) (*pb.Message, error) {
	var conn, err = conn.Conn(socket)
	if err != nil {
		return &pb.Message{Err: err.Error()}, err
	} else {
		err := service.Set(conn, pathandname.Path, pathandname.Data)
		if err != nil {
			return &pb.Message{Err: err.Error()}, err
		} else {
			return &pb.Message{}, nil
		}
	}

}

func ReadConfig() (port, socket string) {
	v := viper.New()
	v.AddConfigPath(".")      // 路径
	v.SetConfigName("config") // 名称
	v.SetConfigType("yaml")   // 类型

	err := v.ReadInConfig() // 读配置
	if err != nil {
		panic(err)
	}

	port = v.GetString("config.port")
	socket = v.GetString("config.socket")
	return

}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("zkserver的grpc服务监听失败: %v", err)
	}
	s := grpc.NewServer() //起一个服务

	pb.RegisterZkServiceServer(s, &server{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
