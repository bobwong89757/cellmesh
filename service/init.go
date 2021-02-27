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
	log.GetLog().Info("Execuable: %s", os.Args[0])
	log.GetLog().Info("WorkDir: %s", workdir)
	log.GetLog().Info("ProcName: '%s'", GetProcName())
	log.GetLog().Info("PID: %d", os.Getpid())
	log.GetLog().Info("Discovery: '%s'", flagDiscoveryAddr)
	log.GetLog().Info("LinkRule: '%s'", getLinkRule())
	log.GetLog().Info("SvcGroup: '%s'", GetSvcGroup())
	log.GetLog().Info("SvcIndex: %d", GetSvcIndex())
	log.GetLog().Info("LANIP: '%s'", util.GetLocalIP())
	log.GetLog().Info("WANIP: '%s'", GetWANIP())
}

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	log.GetLog().Debug("Connecting to discovery '%s' ...", flagDiscoveryAddr)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = flagDiscoveryAddr
	discovery.Default = memsd.NewDiscovery(sdConfig)
}

func WaitExitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
