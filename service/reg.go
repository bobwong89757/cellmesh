package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/util"
)

// peerListener 是获取监听端口的接口
// 用于从Peer中获取实际监听的端口号
type peerListener interface {
	Port() int // 返回监听端口号
}

// ServiceMeta 是服务元数据的类型定义
// 用于在注册服务时传递额外的元数据信息
type ServiceMeta map[string]string

// Register 将Acceptor注册到服务发现系统
// 会自动获取本地IP和监听端口，并设置服务的基本元数据
// 参数:
//   - p: 要注册的Peer实例，必须是Acceptor类型
//   - options: 可选的配置选项，支持ServiceMeta类型用于设置额外元数据
//
// 返回:
//   - *discovery.ServiceDesc: 注册的服务描述信息
func Register(p cellnet.Peer, options ...interface{}) *discovery.ServiceDesc {
	host := util.GetLocalIP()

	property := p.(cellnet.PeerProperty)

	sd := &discovery.ServiceDesc{
		ID:   MakeLocalSvcID(property.Name()),
		Name: property.Name(),
		Host: host,
		Port: p.(peerListener).Port(),
	}

	sd.SetMeta("SvcGroup", GetSvcGroup())
	sd.SetMeta("SvcIndex", GetSvcIndex())

	for _, opt := range options {

		switch optValue := opt.(type) {
		case ServiceMeta:
			for metaKey, metaValue := range optValue {
				sd.SetMeta(metaKey, metaValue)
			}
		}
	}

	if GetWANIP() != "" {
		sd.SetMeta("WANAddress", util.JoinAddress(GetWANIP(), sd.Port))
	}

	log.GetLog().Debugf("service '%s' listen at port: %d", sd.ID, sd.Port)

	p.(cellnet.ContextSet).SetContext("sd", sd)

	// 有同名的要先解除注册，再注册，防止watch不触发
	discovery.Default.Deregister(sd.ID)
	err := discovery.Default.Register(sd)
	if err != nil {
		log.GetLog().Errorf("service register failed, %s %s", sd.String(), err.Error())
	}

	return sd
}

// Unregister 从服务发现系统中注销Peer
// 参数:
//   - p: 要注销的Peer实例
func Unregister(p cellnet.Peer) {
	property := p.(cellnet.PeerProperty)
	discovery.Default.Deregister(MakeLocalSvcID(property.Name()))
}
