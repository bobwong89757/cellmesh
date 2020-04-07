package service

var (
	flagDiscoveryAddr string
	flagLinkRule string
	flagSvcGroup string
	flagSvcIndex string
	flagWANIP string
	flagCommType string
)

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