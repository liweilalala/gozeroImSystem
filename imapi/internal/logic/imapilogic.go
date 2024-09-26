package logic

import (
	"context"

	"gozeroImSystem/imapi/internal/svc"
	"gozeroImSystem/imapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImapiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImapiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImapiLogic {
	return &ImapiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImapiLogic) Imapi(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
