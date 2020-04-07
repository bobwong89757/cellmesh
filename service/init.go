package service

import (
	"github.com/bobwong89757/cellmesh/discovery"
	"github.com/bobwong89757/cellmesh/discovery/memsd/api"
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
	log.Info("Execuable: %s", os.Args[0])
	log.Info("WorkDir: %s", workdir)
	log.Info("ProcName: '%s'", GetProcName())
	log.Info("PID: %d", os.Getpid())
	log.Info("Discovery: '%s'", flagDiscoveryAddr)
	log.Info("LinkRule: '%s'", getLinkRule())
	log.Info("SvcGroup: '%s'", GetSvcGroup())
	log.Info("SvcIndex: %d", GetSvcIndex())
	log.Info("LANIP: '%s'", util.GetLocalIP())
	log.Info("WANIP: '%s'", GetWANIP())
}

// 连接到服务发现, 建议在service.Init后, 以及服务器逻辑开始前调用
func ConnectDiscovery() {
	log.Debug("Connecting to discovery '%s' ...", flagDiscoveryAddr)
	sdConfig := memsd.DefaultConfig()
	sdConfig.Address = flagDiscoveryAddr
	discovery.Default = memsd.NewDiscovery(sdConfig)
}

func WaitExitSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch
}
