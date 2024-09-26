package logic

import (
	"context"

	"gozeroImSystem/common/libnet"
	"gozeroImSystem/common/session"
	"gozeroImSystem/common/socket"
	"gozeroImSystem/edge/internal/svc"
	"gozeroImSystem/imrpc/imrpc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"google.golang.org/protobuf/proto"
)

type MqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	server *socket.Server
	logx.Logger
}

func NewMqLogic(ctx context.Context, svcCtx *svc.ServiceContext, server *socket.Server) *MqLogic {
	return &MqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		server: server,
		Logger: logx.WithContext(ctx),
	}
}

// 消费逻辑，edge服务
func (l *MqLogic) Consume(ctx context.Context, key, val string) error {
	var msg imrpc.PostMessageRequest
	// 先将消息反序列化得到结构体
	err := proto.Unmarshal([]byte(val), &msg)
	if err != nil {
		logx.Errorf("[Consume] proto.Unmarshal val: %s error: %v", val, err)
		return err
	}
	logx.Infof("[Consume] succ msg: %+v body: %s", msg, string(msg.Body))

	if len(msg.ToToken) > 0 {
		// 获取token对应的session
		sessions := l.server.Manager.GetTokenSessions(msg.ToToken)
		for i := range sessions {
			s := sessions[i]
			if s == nil {
				logx.Errorf("[Consume] session not found, msg: %v", msg)
				continue
			}
			// 调用session的Send方法发送消息
			err := s.Send(makeMessage(&msg))
			if err != nil {
				logx.Errorf("[Consume] session send error, msg: %v, err: %v", msg, err)
			}
		}
	} else {
		s := l.server.Manager.GetSession(session.FromString(msg.SessionId))
		if s == nil {
			logx.Errorf("[Consume] session not found, msg: %v", msg)
			return nil
		}
		return s.Send(makeMessage(&msg))
	}

	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext, server *socket.Server) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConf, NewMqLogic(ctx, svcCtx, server)),
	}
}

func makeMessage(msg *imrpc.PostMessageRequest) libnet.Message {
	var message libnet.Message
	message.Version = uint8(msg.Version)
	message.Status = uint8(msg.Status)
	message.ServiceId = uint16(msg.ServiceId)
	message.Cmd = uint16(msg.Cmd)
	message.Seq = msg.Seq
	message.Body = msg.Body
	return message
}
