package service

var (
	flagDiscoveryAddr string // 服务发现服务器地址
	flagLinkRule      string // 服务互联规则
	flagSvcGroup      string // 服务分组
	flagSvcIndex      string // 服务索引
	flagWANIP         string // 外网IP地址
	flagCommType      string // 通信类型
)

// InitServerConfig 初始化服务器配置
// 从配置映射中读取并设置各种服务参数
// 参数:
//   - serviceConf: 配置映射，键名包括: sdaddr, linkrule, svcgroup, svcindex, wanip, commtype
func InitServerConfig(serviceConf map[string]string) {
	// 服务发现地址
	flagDiscoveryAddr = serviceConf["sdaddr"]

	// 服务发现规则
	flagLinkRule = serviceConf["linkrule"]

	// 服务所在组
	flagSvcGroup = serviceConf["svcgroup"]

	// 服务索引
	flagSvcIndex = serviceConf["svcindex"]

	// 设置外网IP
	flagWANIP = serviceConf["wanip"]

	// 通讯类型
	flagCommType = serviceConf["commtype"]
}