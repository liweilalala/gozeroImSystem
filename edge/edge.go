package main

import (
	"context"
	"flag"
	"fmt"

	"gozeroImSystem/common/libnet"
	"gozeroImSystem/common/socket"
	"gozeroImSystem/edge/internal/config"
	"gozeroImSystem/edge/internal/logic"
	"gozeroImSystem/edge/internal/server"
	"gozeroImSystem/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	zeroservice "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
)

var configFile = flag.String("f", "etc/edge.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	srvCtx := svc.NewServiceContext(c)
	// 通常情况下，日志库可能会收集一些关于日志记录的统计数据，例如日志的总数量、日志级别的分布等。
	// 调用 logx.DisableStat() 后，logx 将停止收集或报告这些统计信息。
	logx.DisableStat()

	// 创建TCP server
	tcpServer := server.NewTCPServer(srvCtx)
	protobuf := libnet.NewIMProtocol()
	var err error
	tcpServer.Server, err = socket.NewServe(c.Name, c.TCPListenOn, protobuf, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting tcp server at %s...\n", c.TCPListenOn)
	// 通过go-zero框架提供的GoSafe方法创建一个安全的goroutine，因为出现异常的时候会被捕获。
	threading.GoSafe(func() {
		fmt.Println("handle tcp server request")
		tcpServer.HandleRequest()
	})
	threading.GoSafe(func() {
		// kq注册，把kafka地址注册到etcd中
		tcpServer.KqHeart()
	})

	// // 创建grpc服务实例
	// s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
	// 	edge.RegisterEdgeServer(grpcServer, server.NewEdgeServer(ctx))
	// 	if c.Mode == service.DevMode || c.Mode == service.TestMode {
	// 		reflection.Register(grpcServer)
	// 	}
	// })
	// defer s.Stop()
	// fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	// s.Start()
	serviceGroup := zeroservice.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range logic.Consumers(context.Background(), srvCtx, tcpServer.Server) {
		serviceGroup.Add(mq)
	}
	serviceGroup.Start()
}
