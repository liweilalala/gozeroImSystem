package logic

import (
	"context"

	"gozeroImSystem/imrpc/internal/svc"
	"gozeroImSystem/imrpc/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImrpcLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImrpcLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImrpcLogic {
	return &ImrpcLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImrpcLogic) Imrpc(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
