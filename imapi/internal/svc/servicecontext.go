package svc

import (
	"gozeroImSystem/imapi/internal/config"
	"gozeroImSystem/imrpc/imrpc"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMRpc  imrpc.ImrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient := zrpc.MustNewClient(c.ImRPC)
	return &ServiceContext{
		Config: c,
		IMRpc:  imrpc.NewImrpcClient(userClient.Conn()),
	}
}
