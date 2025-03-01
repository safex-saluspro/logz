package logger

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/pebbe/zmq4"
	"net/http"
)

type NotifierManager interface {
	WebServer() *http.Server
	Websocket() *zmq4.Socket
	WebClient() *http.Client
	DBusClient() *dbus.Conn

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

func (nm *NotifierManagerImpl) WebServer() *http.Server {
	if nm.webServer == nil {
		nm.webServer = Server()
	}
	return nm.webServer
}
func (nm *NotifierManagerImpl) Websocket() *zmq4.Socket {
	if nm.websocket == nil {
		nm.websocket = Socket()
	}
	return nm.websocket
}
func (nm *NotifierManagerImpl) WebClient() *http.Client {
	if nm.webClient == nil {
		nm.webClient = Client()
	}
	return nm.webClient
}
func (nm *NotifierManagerImpl) DBusClient() *dbus.Conn {
	if nm.dbusClient == nil {
		nm.dbusClient = DBus()
	}
	return nm.dbusClient
}
