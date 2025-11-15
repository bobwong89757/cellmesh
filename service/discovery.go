package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
)

// DiscoveryOption 是服务发现的选项配置
type DiscoveryOption struct {
	Rules         []MatchRule // 匹配规则列表，用于过滤要连接的服务
	MaxCount      int         // 最大连接数，0表示不限制，默认发起多条连接
	MatchSvcGroup string      // 匹配的服务组，空字符串时匹配所有同类服务，否则只连接指定组的服务
}

// DiscoveryService 发现并连接到指定的服务
// 服务可能拥有多个实例，每个实例都会创建一个连接
// 函数会持续监听服务变化，自动处理服务的添加、更新和移除
// 参数:
//   - tgtSvcName: 目标服务名称
//   - opt: 发现选项配置
//   - peerCreator: Peer创建函数，当发现新服务时会调用此函数创建连接
// 返回:
//   - cellnet.Peer: MultiPeer实例，可以通过IsReady()判断所有连接是否已准备好
func DiscoveryService(tgtSvcName string, opt DiscoveryOption, peerCreator func(MultiPeer, *discovery.ServiceDesc)) cellnet.Peer {

	// 从发现到连接有一个过程，需要用Map防止还没连上，又创建一个新的连接
	multiPeer := newMultiPeer()

	go func() {

		notify := discovery.Default.RegisterNotify("add")
		for {

			QueryService(tgtSvcName,
				Filter_MatchRule(opt.Rules),
				Filter_MatchSvcGroup(opt.MatchSvcGroup),
				func(desc *discovery.ServiceDesc) interface{} {

					//log.Info("found '%s' address '%s' ", tgtSvcName, desc.Address())

					prePeer := multiPeer.GetPeer(desc.ID)

					// 如果svcid重复汇报, 可能svcid内容有变化
					if prePeer != nil {

						var preDesc *discovery.ServiceDesc
						if prePeer.(cellnet.ContextSet).FetchContext("sd", &preDesc) && !preDesc.Equals(desc) {

							log.GetLog().Infof("service '%s' change desc, %+v -> %+v...", desc.ID, preDesc, desc)

							// 移除之前的连接
							multiPeer.RemovePeer(desc.ID)

							// 停止重连
							prePeer.Stop()

						} else {
							return true
						}

					}

					// 达到最大连接
					if opt.MaxCount > 0 && len(multiPeer.GetPeers()) >= opt.MaxCount {
						return true
					}

					// 用户创建peer
					peerCreator(multiPeer, desc)

					return true
				})

			<-notify
		}

	}()

	return multiPeer
}
