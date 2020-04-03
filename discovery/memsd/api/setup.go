package memsd

import (
	"github.com/bobwong89757/cellmesh/discovery/memsd/model"
	"github.com/bobwong89757/cellnet"
	_ "github.com/bobwong89757/cellnet/peer/tcp"
	"github.com/bobwong89757/cellnet/proc"
	"github.com/bobwong89757/cellnet/proc/tcp"
)

func init() {
	// 仅供demo使用的
	proc.RegisterProcessor("memsd.cli", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(TCPMessageTransmitter))

		if model.Debug {
			bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker), new(typeRPCHooker)))
		} else {
			bundle.SetHooker(new(typeRPCHooker))
		}

		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})

	proc.RegisterProcessor("memsd.svc", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback, args ...interface{}) {

		bundle.SetTransmitter(new(TCPMessageTransmitter))
		if model.Debug {
			bundle.SetHooker(proc.NewMultiHooker(new(tcp.MsgHooker)))
		}

		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}