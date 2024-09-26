package libnet

import "time"

type Codec interface {
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	Receive() (*Message, error)
	Send(Message) error
	Close() error
}
