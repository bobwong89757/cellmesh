package model

import (
	"github.com/bobwong89757/cellmesh/util"
	"github.com/bobwong89757/cellnet"
	"strings"
)

const (
	ServiceKeyPrefix = "_svcdesc_" // 服务描述在KV存储中的键前缀
)

var (
	Queue cellnet.EventQueue              // 事件队列，用于处理所有网络事件
	IDGen = meshutil.NewUUID64Generator() // UUID生成器，用于生成唯一ID

	Listener cellnet.Peer // 监听器Peer，用于接收客户端连接
	Debug    bool         // 调试模式标志

	Version = "0.1.0" // 版本号
)

func IsServiceKey(rawkey string) bool {

	return strings.HasPrefix(rawkey, ServiceKeyPrefix)
}

func GetSvcIDByServiceKey(rawkey string) string {

	if IsServiceKey(rawkey) {
		return rawkey[len(ServiceKeyPrefix):]
	}

	return ""
}

func init() {
	IDGen.AddTimeComponent(8)
	IDGen.AddSeqComponent(8, 0)
}

func GetSessionToken(ses cellnet.Session) (token string) {
	ses.(cellnet.ContextSet).FetchContext("token", &token)

	return
}

func Broadcast(msg interface{}) {
	Listener.(cellnet.TCPAcceptor).VisitSession(func(ses cellnet.Session) bool {
		ses.Send(msg)
		return true
	})
}

func TokenExists(token string) (ret bool) {
	Listener.(cellnet.TCPAcceptor).VisitSession(func(ses cellnet.Session) bool {

		if GetSessionToken(ses) == token {
			ret = true
			return false
		}

		return true
	})

	return
}
