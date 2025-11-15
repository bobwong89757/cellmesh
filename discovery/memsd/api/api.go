package memsd

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellmesh/discovery/memsd/model"
	"github.com/bobwong89757/cellnet"
	"sync"
)

// notifyContext 是通知上下文的内部结构
// 用于跟踪通知通道的注册信息
type notifyContext struct {
	stack string // 注册时的调用栈信息，用于调试
	mode  string // 通知模式
}

// memDiscovery 是memsd服务发现的实现
// 实现了discovery.Discovery接口，提供基于内存的服务发现功能
type memDiscovery struct {
	config *Config // 配置信息

	ses      cellnet.Session // 与服务发现服务器的会话
	sesGuard sync.RWMutex    // 保护会话的读写锁

	kvCache      map[string][]byte // KV配置缓存
	kvCacheGuard sync.RWMutex      // 保护KV缓存的读写锁

	svcCache      map[string][]*discovery.ServiceDesc // 服务缓存，键为服务名，值为服务描述列表
	svcCacheGuard sync.RWMutex                        // 保护服务缓存的读写锁

	notifyMap sync.Map // 通知通道映射，key为channel，value为notifyContext

	initWg *sync.WaitGroup // 初始化等待组，用于等待初始数据拉取完成

	token string // 认证令牌
}

// NewDiscovery 创建一个新的memsd服务发现实例
// 参数:
//   - config: 配置对象，如果为nil则使用默认配置
// 返回:
//   - discovery.Discovery: 服务发现实例
func NewDiscovery(config interface{}) discovery.Discovery {

	if config == nil {
		config = DefaultConfig()
	}

	self := &memDiscovery{
		config:   config.(*Config),
		kvCache:  make(map[string][]byte),
		svcCache: make(map[string][]*discovery.ServiceDesc),
	}

	model.Queue = cellnet.NewEventQueue()
	model.Queue.EnableCapturePanic(true)
	model.Queue.StartLoop()

	self.initWg = new(sync.WaitGroup)
	self.initWg.Add(1)

	self.connect(self.config.Address)

	// 等待拉取初始值
	self.initWg.Wait()
	self.initWg = nil

	return self
}
