// Code generated by goctl. DO NOT EDIT.
// Source: imrpc.proto

package imrpc_client

import (
	"context"

	"gozeroImSystem/imrpc/imrpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	LoginRequest   = imrpc.LoginRequest
	LoginResponse  = imrpc.LoginResponse
	LogoutRequest  = imrpc.LogoutRequest
	LogoutResponse = imrpc.LogoutResponse

	Imrpc interface {
		Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
		Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
	}

	defaultImrpc struct {
		cli zrpc.Client
	}
)

func NewImrpc(cli zrpc.Client) Imrpc {
	return &defaultImrpc{
		cli: cli,
	}
}

func (m *defaultImrpc) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	client := imrpc.NewImrpcClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultImrpc) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	client := imrpc.NewImrpcClient(m.cli.Conn())
	return client.Logout(ctx, in, opts...)
}