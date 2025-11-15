package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/peer"
	"sync"
)

// MultiPeer 表示管理多个Peer连接的接口
// 用于一类服务需要连接到多个不同地址的场景，例如login服务需要连接到多个game服务
// 实现了cellnet.Peer接口，可以像普通Peer一样使用
type MultiPeer interface {
	// GetPeers 获取所有管理的Peer列表
	// 返回:
	//   - []cellnet.Peer: Peer列表
	GetPeers() []cellnet.Peer

	// ContextSet 提供上下文设置功能，用于在Peer上存储和获取上下文数据
	cellnet.ContextSet

	// AddPeer 添加一个Peer到管理列表中
	// 参数:
	//   - sd: 服务描述信息
	//   - p: 要添加的Peer实例
	AddPeer(sd *discovery.ServiceDesc, p cellnet.Peer)
}

// multiPeer 是MultiPeer接口的实现
// 内部使用读写锁保护peers列表，支持并发访问
type multiPeer struct {
	peer.CoreContextSet // 提供上下文管理功能
	peers      []cellnet.Peer // 管理的Peer列表
	peersGuard sync.RWMutex   // 保护peers列表的读写锁
	context    interface{}     // 上下文数据
}

func (self *multiPeer) Start() cellnet.Peer {
	return self
}

func (self *multiPeer) Stop() {

}

func (self *multiPeer) TypeName() string {
	return ""
}

func (self *multiPeer) GetPeers() []cellnet.Peer {
	self.peersGuard.RLock()
	defer self.peersGuard.RUnlock()

	return self.peers
}

func (self *multiPeer) IsReady() bool {

	peers := self.GetPeers()

	if len(peers) == 0 {
		return false
	}

	for _, p := range peers {
		if !p.(cellnet.PeerReadyChecker).IsReady() {
			return false
		}
	}

	return true
}

// AddPeer 添加一个Peer到管理列表中
// 注意: 必须在Peer.Start()之前调用，否则连接建立时可能因为缺少服务描述信息而导致服务信息无法正确上报
// 参数:
//   - sd: 服务描述信息，会被设置到Peer的上下文中
//   - p: 要添加的Peer实例
func (self *multiPeer) AddPeer(sd *discovery.ServiceDesc, p cellnet.Peer) {

	contextSet := p.(cellnet.ContextSet)
	contextSet.SetContext("sd", sd)

	self.peersGuard.Lock()
	self.peers = append(self.peers, p)
	self.peersGuard.Unlock()
}

// GetPeer 根据服务ID获取对应的Peer
// 参数:
//   - svcid: 服务的唯一标识ID
// 返回:
//   - cellnet.Peer: 找到的Peer实例，如果不存在则返回nil
func (self *multiPeer) GetPeer(svcid string) cellnet.Peer {
	for _, p := range self.peers {

		if getSvcIDByPeer(p) == svcid {
			return p
		}
	}

	return nil
}

// RemovePeer 从管理列表中移除指定服务ID的Peer
// 参数:
//   - svcid: 要移除的服务的唯一标识ID
func (self *multiPeer) RemovePeer(svcid string) {
	self.peersGuard.Lock()
	defer self.peersGuard.Unlock()
	for index, p := range self.peers {

		if getSvcIDByPeer(p) == svcid {
			self.peers = append(self.peers[:index], self.peers[index+1:]...)
			break
		}
	}
}

// getSvcIDByPeer 从Peer的上下文中获取服务ID
// 这是一个内部辅助函数
// 参数:
//   - p: Peer实例
// 返回:
//   - string: 服务ID，如果不存在则返回空字符串
func getSvcIDByPeer(p cellnet.Peer) string {
	var sd *discovery.ServiceDesc
	if p.(cellnet.ContextSet).FetchContext("sd", &sd) {
		return sd.ID
	}

	return ""
}

// newMultiPeer 创建一个新的MultiPeer实例
// 返回:
//   - *multiPeer: MultiPeer实例
func newMultiPeer() *multiPeer {
	return &multiPeer{}
}
