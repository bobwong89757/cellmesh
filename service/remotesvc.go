package service

import (
	"sync"

	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
)

// RemoteServiceContext 是远程服务的上下文信息
// 存储已连接远程服务的基本信息
type RemoteServiceContext struct {
	Name  string // 服务名称
	SvcID string // 服务唯一标识ID
}

// NotifyFunc 是远程服务通知回调函数类型
// 当远程服务状态发生变化时会被调用
// 参数:
//   - ctx: 远程服务上下文
//   - ses: 对应的会话对象
type NotifyFunc func(ctx *RemoteServiceContext, ses cellnet.Session)

var (
	connBySvcID        = map[string]cellnet.Session{}
	connBySvcNameGuard sync.RWMutex
	removeNotify       NotifyFunc
)

// AddRemoteService 添加一个远程服务到管理列表
// 当服务间建立连接时调用，用于记录已连接的远程服务
// 参数:
//   - ses: 会话对象
//   - svcid: 服务唯一标识ID
//   - name: 服务名称
func AddRemoteService(ses cellnet.Session, svcid, name string) {

	connBySvcNameGuard.Lock()
	ses.(cellnet.ContextSet).SetContext("ctx", &RemoteServiceContext{Name: name, SvcID: svcid})
	connBySvcID[svcid] = ses
	connBySvcNameGuard.Unlock()

	log.GetLog().Infof("remote service added: '%s' sid: %d", svcid, ses.ID())
}

// RemoveRemoteService 从管理列表中移除远程服务
// 当服务间连接断开时调用
// 参数:
//   - ses: 会话对象
func RemoveRemoteService(ses cellnet.Session) {

	if ses == nil {
		return
	}

	ctx := SessionToContext(ses)
	if ctx != nil {

		if removeNotify != nil {
			removeNotify(ctx, ses)
		}

		connBySvcNameGuard.Lock()
		delete(connBySvcID, ctx.SvcID)
		connBySvcNameGuard.Unlock()

		log.GetLog().Infof("remote service removed '%s' sid: %d", ctx.SvcID, ses.ID())
	} else {
		log.GetLog().Infof("remote service removed sid: %d, context lost", ses.ID())
	}
}

// SetRemoteServiceNotify 设置远程服务状态变化的通知回调
// 参数:
//   - mode: 通知模式，目前支持"remove"（服务移除）
//   - callback: 通知回调函数
func SetRemoteServiceNotify(mode string, callback NotifyFunc) {

	switch mode {
	case "remove":
		removeNotify = callback
	default:
		panic("unknown notify mode")
	}
}

// SessionToContext 从会话中获取远程服务上下文
// 参数:
//   - ses: 会话对象
//
// 返回:
//   - *RemoteServiceContext: 远程服务上下文，如果不存在则返回nil
func SessionToContext(ses cellnet.Session) *RemoteServiceContext {
	if ses == nil {
		return nil
	}

	if raw, ok := ses.(cellnet.ContextSet).GetContext("ctx"); ok {
		return raw.(*RemoteServiceContext)
	}

	return nil
}

// GetRemoteService 根据服务ID获取远程服务的会话
// 参数:
//   - svcid: 服务唯一标识ID
//
// 返回:
//   - cellnet.Session: 对应的会话对象，如果不存在则返回nil
func GetRemoteService(svcid string) cellnet.Session {
	connBySvcNameGuard.RLock()
	defer connBySvcNameGuard.RUnlock()

	if ses, ok := connBySvcID[svcid]; ok {

		return ses
	}

	return nil
}

// VisitRemoteService 遍历所有已连接的远程服务
// 参数:
//   - callback: 回调函数，参数为会话和上下文，返回false时停止遍历
func VisitRemoteService(callback func(ses cellnet.Session, ctx *RemoteServiceContext) bool) {
	connBySvcNameGuard.RLock()

	for _, ses := range connBySvcID {

		if !callback(ses, SessionToContext(ses)) {
			break
		}
	}

	connBySvcNameGuard.RUnlock()
}
