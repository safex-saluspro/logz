package logger

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/pebbe/zmq4"
	"net/http"
)

var (
	lServer *http.Server
	lSocket *zmq4.Socket
	lClient *http.Client
	lDBus   *dbus.Conn
)

type NotifierManager interface {
	WebServer(webServer *http.Server) *http.Server
	Websocket(websocket *zmq4.Socket) *zmq4.Socket
	WebClient(webClient *http.Client) *http.Client
	DBusClient(dBusClient *dbus.Conn) *dbus.Conn

	AddNotifier(name string, notifier Notifier)
	RemoveNotifier(name string)
	GetNotifier(name string) (*Notifier, bool)
	ListNotifiers() []string
}
type NotifierManagerImpl struct {
	webServer  *http.Server
	websocket  *zmq4.Socket
	webClient  *http.Client
	dbusClient *dbus.Conn
	notifiers  map[string]Notifier
}

func NewNotifierManager(notifiers map[string]Notifier) *NotifierManager {
	if notifiers == nil {
		notifiers = make(map[string]Notifier)
	}
	ntfMgr := &NotifierManagerImpl{notifiers: notifiers}

	var ntfM NotifierManager
	ntfM = ntfMgr

	return &ntfM
}

func (nm *NotifierManagerImpl) AddNotifier(name string, notifier Notifier) {
	nm.notifiers[name] = notifier
	fmt.Printf("Notifier '%s' added/updated.\n", name)
}
func (nm *NotifierManagerImpl) RemoveNotifier(name string) {
	delete(nm.notifiers, name)
	fmt.Printf("Notifier '%s' removed.\n", name)
}
func (nm *NotifierManagerImpl) GetNotifier(name string) (*Notifier, bool) {
	if notifier, ok := nm.notifiers[name]; ok {
		return &notifier, true
	} else {
		return nil, false
	}
}
func (nm *NotifierManagerImpl) ListNotifiers() []string {
	keys := make([]string, 0, len(nm.notifiers))
	for name := range nm.notifiers {
		keys = append(keys, name)
	}
	return keys
}

func (nm *NotifierManagerImpl) WebServer(webServer *http.Server) *http.Server {
	if webServer != nil {
		lServer = webServer
		nm.webServer = webServer
	} else if nm.webServer == nil {
		nm.webServer = lServer
	}
	return nm.webServer
}
func (nm *NotifierManagerImpl) Websocket(websocket *zmq4.Socket) *zmq4.Socket {
	if websocket != nil {
		lSocket = websocket
		nm.websocket = websocket
	} else if nm.websocket == nil {
		nm.websocket = lSocket
	}
	return nm.websocket
}
func (nm *NotifierManagerImpl) WebClient(webClient *http.Client) *http.Client {
	if webClient != nil {
		lClient = webClient
		nm.webClient = webClient
	} else if nm.webClient == nil {
		nm.webClient = lClient
	}
	return nm.webClient
}
func (nm *NotifierManagerImpl) DBusClient(dBusClient *dbus.Conn) *dbus.Conn {
	if dBusClient != nil {
		lDBus = dBusClient
		nm.dbusClient = dBusClient
	} else if nm.dbusClient == nil {
		nm.dbusClient = lDBus
	}
	return nm.dbusClient
}
