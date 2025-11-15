package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellmesh/discovery/memsd/api"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/util"
	"os"
	"os/signal"
	"syscall"
)

// Init 初始化服务框架
// 设置进程名称并解析服务互联规则
// 参数:
//   - name: 进程名称，用于标识当前服务
func Init(name string) {
	procName = name
	LinkRules = ParseMatchRule(getLinkRule())
}

// getLinkRule 获取服务互联规则
// 如果未设置linkrule，则使用svcgroup作为默认规则
// 返回:
//   - string: 互联规则字符串
func getLinkRule() string {
	if flagLinkRule == "" {
		return flagSvcGroup
	} else {
		return flagLinkRule
	}
}

// LogParameter 打印当前服务的所有参数信息
// 包括可执行文件路径、工作目录、进程名、PID、服务发现地址等
func LogParameter() {
	workdir, _ := os.Getwd()
	log.GetLog().Infof("Execuable: %s", os.Args[0])
	log.GetLog().Infof("WorkDir: %s", workdir)
	log.GetLog().Infof("ProcName: '%s'", GetProcName())
	log.GetLog().Infof("PID: %d", os.Getpid())
	log.GetLog().Infof("Discovery: '%s'", flagDiscoveryAddr)
	log.GetLog().Infof("LinkRule: '%s'", getLinkRule())
	log.GetLog().Infof("SvcGroup: '%s'", GetSvcGroup())
	log.GetLog().Infof("SvcIndex: %d", GetSvcIndex())
	log.GetLog().Infof("LANIP: '%s'", util.GetLocalIP())
	log.GetLog().Infof("WANIP: '%s'", GetWANIP())
}

// ConnectDiscovery 连接到服务发现服务器
// 建议在service.Init()之后、服务器逻辑开始之前调用
// 函数会阻塞直到连接建立并完成初始化
func ConnectDiscovery() {
	log.GetLog().Debugf("Connecting to discovery '%s' ...", flagDiscoveryAddr)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = flagDiscoveryAddr
	discovery.Default = memsd.NewDiscovery(sdConfig)
}

// WaitExitSignal 等待退出信号
// 阻塞当前goroutine，直到收到SIGTERM、SIGINT或SIGQUIT信号
// 通常用于主函数中等待程序退出
func WaitExitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
