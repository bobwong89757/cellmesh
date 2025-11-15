package deps

import (
	"github.com/bobwong89757/cellmesh/discovery"
	memsd "github.com/bobwong89757/cellmesh/discovery/memsd/api"
)

// DiscoveryExtend 是扩展的服务发现接口
// 在discovery.Discovery基础上增加了额外的功能
type DiscoveryExtend interface {
	discovery.Discovery // 基础服务发现接口

	// QueryAll 查询所有已注册的服务
	// 返回:
	//   - ret: 所有服务的描述列表
	QueryAll() (ret []*discovery.ServiceDesc)

	// ClearKey 清空所有KV配置
	ClearKey()

	// ClearService 清空所有服务注册
	ClearService()

	// GetRawValueList 获取指定前缀的所有KV配置
	// 参数:
	//   - prefix: 键前缀
	// 返回:
	//   - ret: 匹配的ValueMeta列表
	GetRawValueList(prefix string) (ret []discovery.ValueMeta)
}

// InitSD 初始化服务发现客户端
// 参数:
//   - arg: 服务发现服务器地址，如果为空则使用默认地址
// 返回:
//   - DiscoveryExtend: 服务发现实例
func InitSD(arg *string) DiscoveryExtend {
	config := memsd.DefaultConfig()
	if *arg != "" {
		config.Address = *arg
	}

	return memsd.NewDiscovery(config).(DiscoveryExtend)
}
