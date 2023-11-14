package manager

import (
	"sync"
	"time"

	"github.com/brunoga/unitybridge"
	"github.com/brunoga/unitybridge/support/finder"
	"github.com/brunoga/unitybridge/support/logger"
	"github.com/brunoga/unitybridge/support/token"
	"github.com/brunoga/unitybridge/unity/event"
	"github.com/brunoga/unitybridge/unity/key"
	"github.com/brunoga/unitybridge/unity/result"
)

const (
	subTypeConnectionOpen = iota
	subTypeConnectionClose
	subTypeConnectionSetIP
	subTypeConnectionSetPort
)

type Connection struct {
	ub    unitybridge.UnityBridge
	l     *logger.Logger
	appID uint64

	f                     *finder.Finder
	connectionStatusToken token.Token

	m         sync.Mutex
	cond      *sync.Cond
	connected bool
}

func NewConnection(ub unitybridge.UnityBridge,
	l *logger.Logger, appID uint64) (*Connection, error) {

	cm := &Connection{
		ub:    ub,
		l:     l,
		f:     finder.New(appID, l),
		appID: appID,
	}

	cm.cond = sync.NewCond(&cm.m)

	return cm, nil
}

func (cm *Connection) Start() error {
	cm.m.Lock()

	token, err := cm.ub.AddKeyListener(key.KeyAirLinkConnection,
		cm.connectionStatusChanged, false)
	if err != nil {
		cm.m.Unlock()
		return err
	}
	cm.connectionStatusToken = token

	cm.m.Unlock()

	b, err := cm.f.Find(30 * time.Second)
	if err != nil {
		return err
	}

	cm.f.SendACK(b.SourceIp(), b.AppId())

	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err = cm.ub.SendEvent(e)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetIP)
	err = cm.ub.SendEventWithString(e, b.SourceIp().String())
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionSetPort)
	err = cm.ub.SendEventWithUint64(e, 10607)
	if err != nil {
		return err
	}

	e.ResetSubType(subTypeConnectionOpen)
	err = cm.ub.SendEvent(e)
	if err != nil {
		return err
	}

	cm.waitForConnectionStatusChange()

	return nil
}

func (cm *Connection) Stop() error {
	e := event.NewFromType(event.TypeConnection)

	e.ResetSubType(subTypeConnectionClose)
	err := cm.ub.SendEvent(e)
	if err != nil {
		return err
	}

	cm.waitForConnectionStatusChange()

	err = cm.ub.RemoveKeyListener(key.KeyAirLinkConnection,
		cm.connectionStatusToken)
	if err != nil {
		return err
	}

	return nil
}

func (cm *Connection) Connected() bool {
	cm.m.Lock()
	defer cm.m.Unlock()

	return cm.connected
}

func (cm *Connection) waitForConnectionStatusChange() {
	cm.m.Lock()
	defer cm.m.Unlock()

	current := cm.connected
	for cm.connected == current {
		cm.cond.Wait()
	}
}

func (cm *Connection) connectionStatusChanged(r *result.Result) {
	connected, ok := r.Value().(bool)
	if !ok {
		cm.l.Error("Unexpected connection status value", "value", r.Value())
		return
	} else if cm.connected != connected {
		cm.l.Debug("Connection status changed", "connected", connected)

		cm.m.Lock()

		cm.connected = connected
		cm.cond.Broadcast()

		cm.m.Unlock()

		return
	}
}
