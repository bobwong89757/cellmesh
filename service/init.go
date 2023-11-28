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

func Init(name string) {
	procName = name
	LinkRules = ParseMatchRule(getLinkRule())
}

func getLinkRule() string {
	if flagLinkRule == "" {
		return flagSvcGroup
	} else {
		return flagLinkRule
	}
}

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

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	log.GetLog().Debugf("Connecting to discovery '%s' ...", flagDiscoveryAddr)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = flagDiscoveryAddr
	discovery.Default = memsd.NewDiscovery(sdConfig)
}

func WaitExitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
