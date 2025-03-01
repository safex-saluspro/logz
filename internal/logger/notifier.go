package logger

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/pebbe/zmq4"
	"io"
	"net/http"
	"strings"
)

type Notifier interface {
	Notify(entry *LogEntry) error
	Enable()
	Disable()
	Enabled() bool

	WebServer() *http.Server
	Websocket() *zmq4.Socket
	WebClient() *http.Client
	DBusClient() *dbus.Conn

	ReturnURL(returnURL *string) string
	HttpMethod(httpMethod *string) string
	LogLevel(loglevel *string) string
	WsEndpoint(wsEndpoint *string) string
	WebhookURL(webhookURL *string) string
	AuthToken(authToken *string) string
	Whitelist(whiteList []string) []string
}
type NotifierImpl struct {
	VlNotifierManager NotifierManager
	VlEnabled         bool     `json:"enabled"`
	VlReturnURL       string   `json:"returnURL"`
	VlWebhookURL      string   `json:"webhookURL"`
	VlHttpMethod      string   `json:"httpMethod"`
	VlAuthToken       string   `json:"authToken"`
	VlLogLevel        string   `json:"logLevel"`
	VlWsEndpoint      string   `json:"wsEndpoint"`
	VlWhitelist       []string `json:"whitelist"`
}

func NewNotifier(manager NotifierManager, enabled bool, webhookURL, httpMethod, authToken, logLevel, wsEndpoint string, whitelist []string) Notifier {
	if whitelist == nil {
		whitelist = make([]string, 0)
	}
	notifier := NotifierImpl{
		VlNotifierManager: manager,
		VlEnabled:         enabled,
		VlWebhookURL:      webhookURL,
		VlHttpMethod:      httpMethod,
		VlAuthToken:       authToken,
		VlLogLevel:        logLevel,
		VlWsEndpoint:      wsEndpoint,
		VlWhitelist:       whitelist,
	}
	return &notifier
}

func (n *NotifierImpl) Notify(entry *LogEntry) error {
	if n.VlEnabled {
		if n.VlLogLevel != "" {
			if n.VlLogLevel != string(entry.Level) {
				return nil
			}
		}
		if n.VlWebhookURL != "" {
			if n.VlHttpMethod == "POST" {
				bodyObj, err := http.NewRequest("POST", n.VlWebhookURL, strings.NewReader(entry.Message))
				if err != nil {
					return err
				}
				if n.VlAuthToken != "" {
					bodyObj.Header.Set("Authorization", "Bearer "+n.VlAuthToken)
				}
				bodyObj.Header.Set("Content-Type", "application/json")
				resp, err := n.WebClient().Do(bodyObj)
				if err != nil {
					return err
				}
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.Body)

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("erro ao enviar log para %s: %s", n.VlWebhookURL, resp.Status)
				}
			} else {
				return fmt.Errorf("método HTTP inválido: %s", n.VlHttpMethod)
			}
		}
		if n.VlWsEndpoint != "" {
			if n.VlAuthToken != "" {
				message := n.VlAuthToken + "|" + entry.Message
				if _, err := n.Websocket().Send(message, 0); err != nil {
					return err
				}
			} else {
				if _, err := n.Websocket().Send(entry.Message, 0); err != nil {
					return err
				}
			}
		}
		if n.DBusClient() != nil {
			output := fmt.Sprintf("DBusNotifier: sending log [%s]", entry.Message)
			if n.VlAuthToken != "" {
				output = n.VlAuthToken + "|" + output
			}
			dbusObj := n.DBusClient().Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
			dbusObj.Call("org.freedesktop.Notifications.Notify", 0, "", uint32(0), "", output, []string{}, map[string]dbus.Variant{}, int32(5000))
		}
	}
	return nil
}
func (n *NotifierImpl) Enable()                 { n.VlEnabled = true }
func (n *NotifierImpl) Disable()                { n.VlEnabled = false }
func (n *NotifierImpl) Enabled() bool           { return n.VlEnabled }
func (n *NotifierImpl) WebServer() *http.Server { return n.VlNotifierManager.WebServer(nil) }
func (n *NotifierImpl) Websocket() *zmq4.Socket { return n.VlNotifierManager.Websocket(nil) }
func (n *NotifierImpl) WebClient() *http.Client { return n.VlNotifierManager.WebClient(nil) }
func (n *NotifierImpl) DBusClient() *dbus.Conn  { return n.VlNotifierManager.DBusClient(nil) }

func (n *NotifierImpl) ReturnURL(returnURL *string) string {
	if returnURL != nil {
		n.VlReturnURL = *returnURL
	}
	return n.VlReturnURL
}
func (n *NotifierImpl) HttpMethod(httpMethod *string) string {
	if httpMethod != nil {
		n.VlHttpMethod = *httpMethod
	}
	return n.VlHttpMethod
}
func (n *NotifierImpl) LogLevel(loglevel *string) string {
	if loglevel != nil {
		n.VlLogLevel = *loglevel
	}
	return n.VlLogLevel
}
func (n *NotifierImpl) WsEndpoint(wsEndpoint *string) string {
	if wsEndpoint != nil {
		n.VlWsEndpoint = *wsEndpoint
	}
	return n.VlWsEndpoint
}
func (n *NotifierImpl) WebhookURL(webhookURL *string) string {
	if webhookURL != nil {
		n.VlWebhookURL = *webhookURL
	}
	return n.VlWebhookURL
}
func (n *NotifierImpl) AuthToken(authToken *string) string {
	if authToken != nil {
		n.VlAuthToken = *authToken
	}
	return n.VlAuthToken
}
func (n *NotifierImpl) Whitelist(whiteList []string) []string {
	if whiteList != nil {
		n.VlWhitelist = whiteList
	}
	return n.VlWhitelist
}
