package server

import (
	"context"
	"errors"
	"fmt"
	"gozeroImSystem/common/libnet"
	"gozeroImSystem/edge/internal/svc"
	"net"
	"time"

	"gozeroImSystem/imrpc/imrpc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type TCPServer struct {
	Listener net.Listener
	svcCtx   *svc.ServiceContext
	Protocol libnet.Protocol
}

// func NewTCPServer(svcCtx *svc.ServiceContext, address string) (*TCPServer, error) {
// 	addr, err := net.ResolveTCPAddr("tcp", address)
// 	if err != nil {
// 		return nil, err
// 	}
// 	listener, err := net.ListenTCP("tcp", addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	protocol := libnet.NewIMProtocol()
// 	// 接口转对象，用类型断言，返回对象和error
// 	// 使用类型断言将 Protocol 转换为 *IMProtocol
// 	imProtocol, _ := protocol.(*libnet.IMProtocol)

// 	return &TCPServer{
// 		svcCtx:   svcCtx,
// 		Listener: listener,
// 		Protocol: imProtocol,
// 	}, nil
// }

func NewTCPServer(svcCtx *svc.ServiceContext, address string) (*TCPServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	fmt.Printf("tcp server addr: %v\n", addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	// listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		svcCtx:   svcCtx,
		Listener: listener,
		Protocol: libnet.NewIMProtocol(),
	}, nil
}

func (srv *TCPServer) HandleRequest() {
	var tmpDelay time.Duration
	fmt.Println("Start handle tcp server request")
	for {
		conn, err := srv.Listener.Accept()
		fmt.Printf("Get new conn from: %v\n", conn.RemoteAddr())
		if err != nil {
			fmt.Printf("handle tcp server request error: %v", err)
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				if tmpDelay == 0 {
					// 第一次等待重试时长为10ms
					tmpDelay = 10 * time.Microsecond
				} else {
					// 之后每次等待重试时长x2
					tmpDelay *= 2
				}
				// 等待时长最多为1s
				if max := time.Second; tmpDelay > max {
					tmpDelay = max
				}
				time.Sleep(tmpDelay)
				continue
			}
		}
		// 每一个接收到的连接，用一个单独的goroutine去处理
		threading.GoSafe(func() {
			srv.handleRequest(conn)
		})
	}
}

// 每个客户端连接都会创建一个session用来存储本次会话。
// 刚创建连接的第一个消息，需要来进行登录校验，所以在14行会调用imrpc的Login方法进行登录校验。
// 如果登录校验异常本次会话就会被关闭，同时关闭连接。
// 登录校验通过的话，后面就是从会话连接中不断的通过Receive方法接收数据，接收到的数据通过HandlePacket方法进行处理。
func (srv *TCPServer) handleRequest(conn net.Conn) {
	codec := srv.Protocol.NewCodec(conn)
	session := NewSession(srv, codec, srv.svcCtx.Config.SendChanSize)
	msg, err := session.Receive()
	if err != nil {
		logx.Errorf("[HandleRequest] session.Receive error: %v", err)
		_ = session.Close()
		return
	}

	logx.Infof("[HandleRequest] session.Receive msg: %s", msg.Fromat())

	// 登录校验
	if err := srv.Login(session, msg); err != nil {
		logx.Errorf("[HandleRequest] srv.Login error: %v", err)
		_ = session.Close()
		return
	}

	for {
		message, err := session.Receive()
		if err != nil {
			logx.Errorf("[HandleRequest] session.Receive error: %v", err)
			_ = session.Close()
			break
		}

		logx.Infof("[HandleRequest] session.Receive message: %s", message.Fromat())

		if err = session.HandlePacket(message); err != nil {
			logx.Errorf("[HandleRequest] session.HandlePacket message: %v error: %v", message, err)
		}
	}
}

// 关闭监听
func (srv *TCPServer) Close() error {
	return srv.Listener.Close()
}

// 调用rpc方法进行登录校验
func (srv *TCPServer) Login(session *Session, msg *libnet.Message) error {
	_, err := srv.svcCtx.IMRpc.Login(context.Background(), &imrpc.LoginRequest{})
	if err != nil {
		return err
	}
	// 登录成功，将消息发送到conn中
	_ = session.Send(*msg)

	return nil
}

// 调用rpc方法进行登出
func (srv *TCPServer) Logout(session *Session, msg *libnet.Message) error {
	_, err := srv.svcCtx.IMRpc.Logout(context.Background(), &imrpc.LogoutRequest{})
	if err != nil {
		return err
	}
	_ = session.Close()

	return nil
}

// func (srv *TCPServer) KqHeart() {
// 	work := discovery.NewKqWorker(srv.svcCtx.Config.Etcd.Key, srv.svcCtx.Config.Etcd.Hosts, srv.svcCtx.Config.KqConf)
// 	work.HeartBeat()
// }