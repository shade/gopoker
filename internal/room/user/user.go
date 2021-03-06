package user

import (
	"gopoker/internal/room/user/sock"
	"gopoker/pkg/pausabletimer"

	. "gopoker/internal/interfaces"

	"github.com/golang/protobuf/proto"
)

type User struct {
	sock ISock
	id   string

	timer *pausabletimer.PausableTimer
}

func NewUser(id string) IUser {
	return &User{
		sock: sock.NewSock(),
		id:   id,
	}
}

func (u User) GetID() string {
	return u.id
}

func (u *User) Send(msg proto.Message) {
	u.sock.Write(msg)
}

func (u *User) RegisterObserver(msgType proto.GeneratedEnum, cb func(IUser, proto.Message)) {
	u.sock.RegisterObserver(msgType, func(msg proto.Message) {
		cb(u, msg)
	})
}

func (u *User) IgnoreUser(msgType proto.GeneratedEnum) {
	u.sock.DeregisterObservers(msgType)
}
