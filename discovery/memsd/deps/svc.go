package deps

import (
	"github.com/bobwong89757/cellmesh/discovery/memsd/api"
	"github.com/bobwong89757/cellmesh/discovery/memsd/model"
	"github.com/bobwong89757/cellmesh/discovery/memsd/proto"
	"github.com/bobwong89757/cellmesh/service"
	"github.com/bobwong89757/cellnet"
	"github.com/bobwong89757/cellnet/log"
	"github.com/bobwong89757/cellnet/peer"
	"github.com/bobwong89757/cellnet/proc"
	"strings"
)

func StartSvc(arg *string) {

	config := memsd.DefaultConfig()
	if *arg != "" {
		config.Address = *arg
	}

	model.Queue = cellnet.NewEventQueue()
	model.Queue.EnableCapturePanic(true)
	model.Queue.StartLoop()

	p := peer.NewGenericPeer("tcp.Acceptor", "memsd", config.Address, model.Queue)
	p.(cellnet.PeerCaptureIOPanic).EnableCaptureIOPanic(true)

	model.Listener = p
	msgFunc := proto.GetMessageHandler("memsd")

	proc.BindProcessorHandler(p, "memsd.svc", func(ev cellnet.Event) {

		if msgFunc != nil {
			msgFunc(ev)
		}
	})

	// 100M封包大小
	p.(cellnet.TCPSocketOption).SetMaxPacketSize(1024 * 1024 * 100)
	p.(cellnet.TCPSocketOption).SetSocketBuffer(1024*1024, 1024*1024, true)
	p.(cellnet.PeerCaptureIOPanic).EnableCaptureIOPanic(true)
	p.Start()
	service.WaitExitSignal()
}

func DeleteValueRecurse(key, reason string) {

	var keyToDelete []string
	model.VisitValue(func(meta *model.ValueMeta) bool {

		if strings.HasPrefix(meta.Key, key) {
			keyToDelete = append(keyToDelete, meta.Key)
		}

		return true
	})

	for _, key := range keyToDelete {
		DeleteNotify(key, reason)
	}
}

func DeleteNotify(key, reason string) {
	valueMeta := model.DeleteValue(key)

	var ack proto.ValueDeleteNotifyACK
	ack.Key = key

	if valueMeta != nil {
		ack.SvcName = valueMeta.SvcName
	}

	if valueMeta != nil {

		if valueMeta.SvcName == "" {
			log.GetLog().Info("DeleteValue '%s'  reason: %s", key, reason)
		} else {
			log.GetLog().Info("DeregisterService '%s'  reason: %s", model.GetSvcIDByServiceKey(key), reason)
		}
	}

	model.Broadcast(&ack)

}

func CheckAuth(ses cellnet.Session) bool {

	return model.GetSessionToken(ses) != ""
}
