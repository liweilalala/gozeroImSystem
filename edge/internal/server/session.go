package server

import (
	"errors"
	"gozeroImSystem/common/libnet"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	LoginCmd  = uint16(1)
	LogoutCmd = uint16(2)
)

var (
	SessionClosedError  = errors.New("session closed")
	SessionBlockedError = errors.New("session blocked")
)

type Session struct {
	tcpServer *TCPServer
	codec     libnet.Codec
	closeChan chan int
	closeFlag int32
	sendChan  chan libnet.Message
}

func NewSession(tcpServer *TCPServer, codec libnet.Codec, sendChanSize int) *Session {
	sn := &Session{
		tcpServer: tcpServer,
		codec:     codec,
		closeChan: make(chan int),
	}
	if sendChanSize > 0 {
		sn.sendChan = make(chan libnet.Message, sendChanSize)
		go sn.sendLoop()
	}
	return sn
}

func (s *Session) HandlePacket(msg *libnet.Message) error {
	// 心跳处理
	// 登出
	if msg.Cmd == LogoutCmd {
		s.Close()
	}
	// Send方法往连接中写入消息，s.codec中的连接就是请求方创建的连接
	err := s.codec.Send(*msg)
	if err != nil {
		logx.Errorf("[HandlePacket] s.codec.Send msg: %v error: %v", msg, err)
	}
	return nil
}

func (s *Session) Send(msg libnet.Message) error {
	if s.IsClosed() {
		return SessionClosedError
	}
	if s.sendChan == nil {
		return s.codec.Send(msg)
	}
	select {
	// 使用channel缓冲消息
	case s.sendChan <- msg:
		return nil
	default:
		return SessionBlockedError
	}
}

func (s *Session) sendLoop() {
	defer s.Close()
	for {
		select {
		// 异步携程用于消费channel中的消息
		case msg := <-s.sendChan:
			err := s.codec.Send(msg)
			if err != nil {
				logx.Errorf("[sendLoop] s.codec.Send msg: %v error: %v", msg, err)
			}
		case <-s.closeChan:
			return
		}
	}
}

// atomic.CompareAndSwapInt32 是 Go 的标准库中的原子操作函数，
// 它的作用是对给定的整数值进行“比较并交换”操作。这个函数会对 closeFlag 变量进行以下操作：
// 比较：检查 closeFlag 当前的值是否为 0。
// 交换：如果当前值是 0，则将其设置为 1，并返回 true。
// 如果当前值不是 0，则不会进行交换，直接返回 false。
func (s *Session) Close() error {
	if atomic.CompareAndSwapInt32(&s.closeFlag, 0, 1) {
		err := s.codec.Close()
		close(s.closeChan)
		return err
	}
	return SessionClosedError
}

// 检查发送数据管道是否关闭
func (s *Session) IsClosed() bool {
	// 使用sync.atomic保证原子操作
	return atomic.LoadInt32(&s.closeFlag) == 1
}

func (s *Session) Receive() (*libnet.Message, error) {
	return s.codec.Receive()
}
