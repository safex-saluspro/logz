package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pebbe/zmq4"
)

// Notifier define o contrato para envio de log a destinos externos.
type Notifier interface {
	Notify(entry *LogEntry)
}

// ExternalNotifier envia o log para uma URL externa (via HTTP)
// e opcionalmente para um endpoint ZMQ.
type ExternalNotifier struct {
	externalURL string
	zmqSocket   *zmq4.Socket
}

// NewExternalNotifier cria uma nova instância de ExternalNotifier.
func NewExternalNotifier(url string, zmqEndpoint string) *ExternalNotifier {
	var socket *zmq4.Socket
	if zmqEndpoint != "" {
		var err error
		socket, err = zmq4.NewSocket(zmq4.PUSH)
		if err != nil {
			fmt.Printf("Erro ao criar socket ZMQ: %v\n", err)
		} else {
			if err := socket.Connect(zmqEndpoint); err != nil {
				return nil
			}
		}
	}
	return &ExternalNotifier{
		externalURL: url,
		zmqSocket:   socket,
	}
}

// Notify envia o log via HTTP (se externalURL estiver configurada)
// e via ZMQ (se o socket estiver ativo).
func (n *ExternalNotifier) Notify(entry *LogEntry) {
	// HTTP
	if n.externalURL != "" {
		data, _ := json.Marshal(entry)
		_, err := http.Post(n.externalURL, "application/json", strings.NewReader(string(data)))
		if err != nil {
			fmt.Printf("Erro ao enviar log para %s: %v\n", n.externalURL, err)
		}
	}
	// ZMQ
	if n.zmqSocket != nil {
		data, _ := json.Marshal(entry)
		if _, err := n.zmqSocket.Send(string(data), 0); err != nil {
			fmt.Printf("Erro ao enviar log para ZMQ: %v\n", err)
			return
		}
	}
}

// DBusNotifier envia logs para o DBus.
type DBusNotifier struct {
	enabled bool
	// Aqui pode ser armazenada uma conexão DBus, se necessário.
}

// NewDBusNotifier cria uma nova instância de DBusNotifier.
func NewDBusNotifier() *DBusNotifier {
	return &DBusNotifier{
		enabled: false,
	}
}

// Enable ativa o envio de logs via DBus.
func (d *DBusNotifier) Enable() {
	d.enabled = true
	fmt.Println("DBusNotifier enabled.")
}

// Disable desativa o envio de logs via DBus.
func (d *DBusNotifier) Disable() {
	d.enabled = false
	fmt.Println("DBusNotifier disabled.")
}

// Notify envia o log via DBus, se habilitado.
func (d *DBusNotifier) Notify(entry *LogEntry) {
	if !d.enabled {
		return
	}
	// Aqui, implemente a lógica de envio usando a API DBus.
	// Por enquanto, simulamos a operação.
	fmt.Printf("DBusNotifier: sending log [%s] via DBus\n", entry.Message)
}

// ZMQSecNotifier envia logs de forma segura via ZMQ, utilizando autenticação JWT.
type ZMQSecNotifier struct {
	enabled   bool
	zmqSocket *zmq4.Socket
	// Possivelmente, campos para armazenar chaves ou caminhos dos arquivos.
	// Exemplo:
	privateKeyPath string
	certPath       string
	configPath     string
}

// NewZMQSecNotifier cria uma nova instância de ZMQSecNotifier.
func NewZMQSecNotifier(zmqEndpoint, privateKeyPath, certPath, configPath string) *ZMQSecNotifier {
	socket, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		fmt.Printf("Error creating secure ZMQ socket: %v\n", err)
		return nil
	}
	if err := socket.Connect(zmqEndpoint); err != nil {
		fmt.Printf("Error connecting secure ZMQ socket: %v\n", err)
		return nil
	}
	return &ZMQSecNotifier{
		enabled:        false,
		zmqSocket:      socket,
		privateKeyPath: privateKeyPath,
		certPath:       certPath,
		configPath:     configPath,
	}
}

// Enable ativa o ZMQSecNotifier.
func (z *ZMQSecNotifier) Enable() {
	z.enabled = true
	fmt.Println("ZMQSecNotifier enabled.")
}

// Disable desativa o ZMQSecNotifier.
func (z *ZMQSecNotifier) Disable() {
	z.enabled = false
	fmt.Println("ZMQSecNotifier disabled.")
}

// Notify envia o log de forma segura via ZMQ, com autenticação JWT.
// Aqui, você pode usar funções do seu módulo gkbxsrv para gerar token JWT.
func (z *ZMQSecNotifier) Notify(entry *LogEntry) {
	if !z.enabled {
		return
	}
	// Gerar token JWT apropriado (use as funções de autenticação do gkbxsrv, por exemplo).
	// token := GenerateJWTToken(entry) // Implemente conforme necessário.
	// Adicione o token à mensagem ou no cabeçalho dos dados.

	data, _ := json.Marshal(entry)
	// Para exemplo, suponha que concatenamos o token à mensagem:
	// message := token + "|" + string(data)
	// Simulamos a operação:
	message := string(data) // Aqui, você incluiria a lógica de autenticação.

	if _, err := z.zmqSocket.Send(message, 0); err != nil {
		fmt.Printf("Error sending secure log via ZMQ: %v\n", err)
		return
	}
}
