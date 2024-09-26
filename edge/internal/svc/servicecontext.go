package svc

import (
	"gozeroImSystem/edge/internal/config"
	"gozeroImSystem/imrpc/imrpc_client"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	IMRpc  imrpc_client.Imrpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.IMRpc)
	return &ServiceContext{
		Config: c,
		IMRpc:  imrpc_client.NewImrpc(client),
	}
}
