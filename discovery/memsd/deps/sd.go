package deps

import (
	"github.com/bobwong89757/cellmesh/discovery"
	memsd "github.com/bobwong89757/cellmesh/discovery/memsd/api"
)

type DiscoveryExtend interface {
	discovery.Discovery

	QueryAll() (ret []*discovery.ServiceDesc)

	ClearKey()

	ClearService()

	GetRawValueList(prefix string) (ret []discovery.ValueMeta)
}

func InitSD(arg *string) DiscoveryExtend {
	config := memsd.DefaultConfig()
	if *arg != "" {
		config.Address = *arg
	}

	return memsd.NewDiscovery(config).(DiscoveryExtend)
}
