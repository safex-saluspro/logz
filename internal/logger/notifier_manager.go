package logger

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/pebbe/zmq4"
	"github.com/spf13/viper"
	"net/http"
)

type NotifierManager interface {
	WebServer() *http.Server
	Websocket() *zmq4.Socket
	WebClient() *http.Client
	DBusClient() *dbus.Conn

	AddNotifier(name string, notifier Notifier)
	RemoveNotifier(name string)
	GetNotifier(name string) (Notifier, bool)
	ListNotifiers() []string

	UpdateFromConfig(config Config) error // Atualiza os notifiers dinamicamente com base no Config
}

type NotifierManagerImpl struct {
	webServer  *http.Server
	websocket  *zmq4.Socket
	webClient  *http.Client
	dbusClient *dbus.Conn
	notifiers  map[string]Notifier
}

func NewNotifierManager(notifiers map[string]Notifier) NotifierManager {
	if notifiers == nil {
		notifiers = make(map[string]Notifier)
	}
	return &NotifierManagerImpl{
		notifiers: notifiers,
	}
}

func (nm *NotifierManagerImpl) AddNotifier(name string, notifier Notifier) {
	nm.notifiers[name] = notifier
	fmt.Printf("Notifier '%s' added/updated.\n", name)
}

func (nm *NotifierManagerImpl) RemoveNotifier(name string) {
	delete(nm.notifiers, name)
	fmt.Printf("Notifier '%s' removed.\n", name)
}

func (nm *NotifierManagerImpl) GetNotifier(name string) (Notifier, bool) {
	notifier, ok := nm.notifiers[name]
	return notifier, ok
}

func (nm *NotifierManagerImpl) ListNotifiers() []string {
	keys := make([]string, 0, len(nm.notifiers))
	for name := range nm.notifiers {
		keys = append(keys, name)
	}
	return keys
}

func (nm *NotifierManagerImpl) UpdateFromConfig(config Config) error {
	var configNotifiers map[string]map[string]interface{}
	if err := viper.UnmarshalKey("notifiers", &configNotifiers); err != nil {
		return fmt.Errorf("failed to parse notifiers config: %w", err)
	}

	// Atualiza ou recria os notifiers dinamicamente
	for name, conf := range configNotifiers {
		typ, ok := conf["type"].(string)
		if !ok {
			fmt.Printf("Notifier '%s' does not specify a type and will be ignored.\n", name)
			continue
		}

		switch typ {
		case "http":
			webhookURL, _ := conf["webhookURL"].(string)
			authToken, _ := conf["authToken"].(string)
			notifier := NewHTTPNotifier(webhookURL, authToken)
			nm.AddNotifier(name, notifier)
		case "zmq":
			endpoint, _ := conf["endpoint"].(string)
			notifier := NewZMQNotifier(endpoint)
			nm.AddNotifier(name, notifier)
		case "dbus":
			notifier := NewDBusNotifier()
			nm.AddNotifier(name, notifier)
		default:
			fmt.Printf("Unknown notifier type '%s' for notifier '%s'.\n", typ, name)
		}
	}
	return nil
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
