package logger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pebbe/zmq4"
)

// Notifier define o contrato para envio de log a destinos externos.
// Adicionalmente, inclui um método para definir um token de autenticação.
type Notifier interface {
	Notify(entry *LogEntry)
	SetAuthToken(token string)
}

// =============================
// ExternalNotifier
// =============================

// ExternalNotifier envia logs via HTTP e via socket ZMQ simples.
type ExternalNotifier struct {
	externalURL string
	zmqSocket   *zmq4.Socket
	authToken   string
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

// SetAuthToken define o token de autenticação para as requisições.
func (n *ExternalNotifier) SetAuthToken(token string) {
	n.authToken = token
}

// Notify envia o log via HTTP (caso externalURL esteja configurada) e via ZMQ.
func (n *ExternalNotifier) Notify(entry *LogEntry) {
	// Envio via HTTP.
	if n.externalURL != "" {
		data, _ := json.Marshal(entry)
		req, err := http.NewRequest("POST", n.externalURL, strings.NewReader(string(data)))
		if err == nil {
			// Se houver token, define o header de Autorização.
			if n.authToken != "" {
				req.Header.Set("Authorization", "Bearer "+n.authToken)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Erro ao enviar log para %s: %v\n", n.externalURL, err)
			} else {
				resp.Body.Close()
			}
		}
	}
	// Envio via ZMQ.
	if n.zmqSocket != nil {
		data, _ := json.Marshal(entry)
		message := string(data)
		if n.authToken != "" {
			// Exemplo: concatenar o token no início, separado por "|".
			message = n.authToken + "|" + message
		}
		if _, err := n.zmqSocket.Send(message, 0); err != nil {
			fmt.Printf("Erro ao enviar log para ZMQ: %v\n", err)
			return
		}
	}
}

// =============================
// DBusNotifier
// =============================

// DBusNotifier envia logs de forma passiva via DBus.
type DBusNotifier struct {
	enabled   bool
	authToken string
	// Aqui você pode incluir campos para gerenciar uma conexão DBus real.
}

// NewDBusNotifier cria uma nova instância de DBusNotifier.
func NewDBusNotifier() *DBusNotifier {
	return &DBusNotifier{
		enabled: false,
	}
}

// SetAuthToken define o token de autenticação (caso seja necessário).
func (d *DBusNotifier) SetAuthToken(token string) {
	d.authToken = token
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

// Notify envia o log via DBus, se estiver habilitado.
// A implementação real da API DBus deve ser inserida aqui.
func (d *DBusNotifier) Notify(entry *LogEntry) {
	if !d.enabled {
		return
	}
	// Exemplo: formata a mensagem com token, se presente.
	output := fmt.Sprintf("DBusNotifier: sending log [%s]", entry.Message)
	if d.authToken != "" {
		output = d.authToken + "|" + output
	}
	fmt.Println(output)
}

// =============================
// ZMQSecNotifier
// =============================

// ZMQSecNotifier envia logs de forma segura via ZeroMQ, utilizando autenticação por JWT.
// Este notifier é dedicado à conexão com o broker gkbxsrv e não pode ser desativado se o broker estiver presente.
type ZMQSecNotifier struct {
	enabled        bool
	zmqSocket      *zmq4.Socket
	authToken      string // Gerenciado por métodos exclusivos; não exposto para edição externa.
	privateKeyPath string
	certPath       string
	configPath     string
}

// NewZMQSecNotifier cria uma nova instância de ZMQSecNotifier.
// A conexão é estabelecida automaticamente com o broker (por exemplo, "tcp://localhost:5555").
func NewZMQSecNotifier(zmqEndpoint, privateKeyPath, certPath, configPath string) *ZMQSecNotifier {
	socket, err := zmq4.NewSocket(zmq4.PUSH)
	if err != nil {
		fmt.Printf("Erro ao criar socket seguro ZMQ: %v\n", err)
		return nil
	}
	// Neste exemplo, a conexão é forçada ao broker local.
	if err := socket.Connect(zmqEndpoint); err != nil {
		fmt.Printf("Erro ao conectar socket seguro ZMQ: %v\n", err)
		return nil
	}
	return &ZMQSecNotifier{
		enabled:        true, // Conexão incondicional se o gkbxsrv estiver presente.
		zmqSocket:      socket,
		privateKeyPath: privateKeyPath,
		certPath:       certPath,
		configPath:     configPath,
	}
}

// SetAuthToken armazena o token de autenticação para uso interno.
func (z *ZMQSecNotifier) SetAuthToken(token string) {
	z.authToken = token
}

// Enable registra que o notificator está ativado. Para o ZMQSecNotifier, a ativação é incondicional.
func (z *ZMQSecNotifier) Enable() {
	// Essa conexão não pode ser desativada se o gkbxsrv estiver presente.
	z.enabled = true
	fmt.Println("ZMQSecNotifier enabled (non-negotiable).")
}

// Disable não deve permitir desativar o notificator se o broker local estiver presente.
func (z *ZMQSecNotifier) Disable() {
	fmt.Println("ZMQSecNotifier cannot be disabled because gkbxsrv is present on this host.")
	// Opcionalmente, mantenha enabled = true mesmo se esse método for chamado.
	z.enabled = true
}

// Notify envia o log de forma segura via ZMQ, utilizando o token se disponível.
// Aqui, a geração/validação do JWT já deve ser feita externamente; este método apenas gerencia o token.
func (z *ZMQSecNotifier) Notify(entry *LogEntry) {
	if !z.enabled {
		return
	}
	data, _ := json.Marshal(entry)
	message := string(data)
	if z.authToken != "" {
		message = z.authToken + "|" + message
	}
	if _, err := z.zmqSocket.Send(message, 0); err != nil {
		fmt.Printf("Erro ao enviar log via ZMQ seguro: %v\n", err)
		return
	}
}
