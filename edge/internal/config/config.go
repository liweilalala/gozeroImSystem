package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	TCPListenOn  string
	SendChanSize int
	IMRpc        zrpc.RpcClientConf
	// KqConf       kq.KqConf
}
