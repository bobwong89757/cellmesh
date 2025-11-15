package service

import (
	"errors"
	"fmt"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/codec"
	_ "github.com/bobwong89757/cellnet/codec/binary"
	"github.com/bobwong89757/cellnet/relay"
	"github.com/bobwong89757/cellnet/util"
	"reflect"
)

// ServiceIdentifyACK 是服务身份确认消息
// 当服务间建立连接时，用于交换服务身份信息，让双方知道对方的服务名称和ID
type ServiceIdentifyACK struct {
	SvcName string // 服务名称，如"game"、"login"等
	SvcID   string // 服务的唯一标识ID
}

func (self *ServiceIdentifyACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*ServiceIdentifyACK)(nil)).Elem(),
		ID:    int(util.StringHash("service.ServiceIdentifyACK")),
	})
}

var (
	// ErrInvalidRelayMessage 表示事件不是有效的relay消息
	ErrInvalidRelayMessage = errors.New("invalid relay message")
	// ErrInvalidRelayPassthroughType 表示透传数据类型不支持
	ErrInvalidRelayPassthroughType = errors.New("invalid relay passthrough type")
)

// GetPassThrough 从relay事件中提取透传数据
// 透传数据用于在消息转发过程中携带额外的上下文信息，如用户ID等
// 支持的类型：*int64, *[]int64, *string
// 参数:
//   - ev: cellnet事件，必须是relay.RecvMsgEvent类型
//   - ptrList: 指向目标变量的指针列表，用于接收透传数据
// 返回:
//   - error: 提取失败时返回错误信息
func GetPassThrough(ev cellnet.Event, ptrList ...interface{}) error {
	if relayEvent, ok := ev.(*relay.RecvMsgEvent); ok {

		for _, ptr := range ptrList {

			switch valuePtr := ptr.(type) {
			case *int64:
				*valuePtr = relayEvent.PassThroughAsInt64()
			case *[]int64:
				*valuePtr = relayEvent.PassThroughAsInt64Slice()
			case *string:
				*valuePtr = relayEvent.PassThroughAsString()
			default:
				return ErrInvalidRelayPassthroughType
			}
		}

		return nil
	} else {
		return ErrInvalidRelayMessage
	}

}

// Reply 向事件来源回复一个消息
// 这是一个便捷函数，用于在事件处理中快速回复消息
// 参数:
//   - ev: cellnet事件，必须实现replyEvent接口
//   - msg: 要回复的消息对象
// 注意: 如果事件不支持Reply方法，会触发panic
func Reply(ev cellnet.Event, msg interface{}) {

	type replyEvent interface {
		Reply(msg interface{})
	}

	if replyEv, ok := ev.(replyEvent); ok {
		replyEv.Reply(msg)
	} else {
		panic("Require 'ReplyEvent' to reply event")
	}
}
