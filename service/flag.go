package service

import (
	flag2 "flag"
	"fmt"
	"github.com/bobwong89757/cellmesh/helpers"
)

var (
	flagDiscoveryAddr string
	flagLinkRule string
	flagSvcGroup string
	flagSvcIndex string
	flagWANIP string
	flagCommType string
)

func init() {
	config := flag2.String("config", "development", "runtime config type")
	flag2.Parse()
	helpers.MConfig.InitConfig(fmt.Sprintf("./cfg/%s.yml", *config))
	serviceConf := helpers.MConfig.Get("svc")

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