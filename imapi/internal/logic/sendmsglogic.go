package logic

import (
	"context"
	"fmt"

	"gozeroImSystem/imapi/internal/svc"
	"gozeroImSystem/imapi/internal/types"
	"gozeroImSystem/imrpc/imrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMsgLogic) SendMsg(req *types.SendMsgRequest) (*types.SendMsgResponse, error) {
	// 调用这个方法发送消息
	_, err := l.svcCtx.IMRpc.PostMessage(l.ctx, &imrpc.PostMessageRequest{
		Token: fmt.Sprintf("%d", req.ToUserId),
		Body:  []byte(req.Content),
	})
	if err != nil {
		return nil, err
	}

	return &types.SendMsgResponse{}, nil
}
