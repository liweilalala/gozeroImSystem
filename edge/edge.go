package main

import (
	"flag"
	"fmt"

	"gozeroImSystem/edge/edge"
	"gozeroImSystem/edge/internal/config"
	"gozeroImSystem/edge/internal/server"
	"gozeroImSystem/edge/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/edge.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	// 通常情况下，日志库可能会收集一些关于日志记录的统计数据，例如日志的总数量、日志级别的分布等。
	// 调用 logx.DisableStat() 后，logx 将停止收集或报告这些统计信息。
	logx.DisableStat()
	// 创建grpc服务实例
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		edge.RegisterEdgeServer(grpcServer, server.NewEdgeServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 创建TCP server
	tcpServer, err := server.NewTCPServer(ctx, c.TCPListenOn)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = tcpServer.Close()
	}()
	fmt.Printf("Starting tcp server at %s...\n", c.TCPListenOn)

	// 通过go-zero框架提供的GoSafe方法创建一个安全的goroutine，因为出现异常的时候会被捕获。
	threading.GoSafe(func() {
		fmt.Println("handle tcp server request")
		tcpServer.HandleRequest()
	})

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
