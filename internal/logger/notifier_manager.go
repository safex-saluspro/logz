package logger

import (
	"fmt"
	"github.com/faelmori/logz/internal/services"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

// NotifierManager centraliza o gerenciamento dos Notifiers.
// Ele mantém um mapa associando um nome (string) a instâncias de Notifier.
type NotifierManager struct {
	notifiers map[string]Notifier
}

// NewNotifierManager cria uma nova instância de NotifierManager.
func NewNotifierManager() *NotifierManager {
	return &NotifierManager{
		notifiers: make(map[string]Notifier),
	}
}

// AddNotifier adiciona ou atualiza um Notifier com o nome especificado.
func (nm *NotifierManager) AddNotifier(name string, notifier Notifier) {
	nm.notifiers[name] = notifier
	fmt.Printf("Notifier '%s' added/updated.\n", name)
}

// RemoveNotifier remove o Notifier associado ao nome.
func (nm *NotifierManager) RemoveNotifier(name string) {
	delete(nm.notifiers, name)
	fmt.Printf("Notifier '%s' removed.\n", name)
}

// GetNotifier retorna o Notifier associado ao nome.
func (nm *NotifierManager) GetNotifier(name string) (Notifier, bool) {
	notifier, ok := nm.notifiers[name]
	return notifier, ok
}

// ListNotifiers retorna uma lista com os nomes de todos os Notifiers registrados.
func (nm *NotifierManager) ListNotifiers() []string {
	keys := make([]string, 0, len(nm.notifiers))
	for name := range nm.notifiers {
		keys = append(keys, name)
	}
	return keys
}

// UpdateFromConfig atualiza o estado do NotifierManager a partir da configuração.
// Espera que a configuração contenha uma seção "notifiers" estruturada, por exemplo:
// notifiers:
//
//	external:
//	  type: "external"
//	  externalURL: "https://discord.com/api/webhooks/XYZ"
//	  zmqEndpoint: "tcp://localhost:5556"
//	dbus:
//	  type: "dbus"
//	  enabled: true
//	zmqsec:
//	  type: "zmqsec"
//	  enabled: true
//	  zmqEndpoint: "tcp://localhost:5555"
//	  privateKeyPath: "/path/to/kubex-key.pem"
//	  certPath: "/path/to/kubex-cert.pem"
//	  configPath: "/path/to/config.json"
func (nm *NotifierManager) UpdateFromConfig() error {
	var configNotifiers map[string]map[string]interface{}
	if err := viper.UnmarshalKey("notifiers", &configNotifiers); err != nil {
		return fmt.Errorf("failed to parse notifiers config: %w", err)
	}

	for name, conf := range configNotifiers {
		typ, ok := conf["type"].(string)
		if !ok {
			// Ignora se não houver um tipo definido.
			continue
		}

		switch typ {
		case "external":
			externalURL, _ := conf["externalURL"].(string)
			zmqEndpoint, _ := conf["zmqEndpoint"].(string)
			notifier := NewExternalNotifier(externalURL, zmqEndpoint)
			// Se houver token, podemos verificá-lo
			if token, ok := conf["authToken"].(string); ok {
				notifier.SetAuthToken(token)
			}
			nm.AddNotifier(name, notifier)

		case "dbus":
			enabled, _ := conf["enabled"].(bool)
			notifier := NewDBusNotifier()
			if enabled {
				notifier.Enable()
			} else {
				notifier.Disable()
			}
			// Se houver token, podemos setar
			if token, ok := conf["authToken"].(string); ok {
				notifier.SetAuthToken(token)
			}
			nm.AddNotifier(name, notifier)

		case "zmqsec":
			enabled, _ := conf["enabled"].(bool)
			// Assume que os demais campos são obrigatórios para essa configuração.
			zmqEndpoint, _ := conf["zmqEndpoint"].(string)
			privateKeyPath, _ := conf["privateKeyPath"].(string)
			certPath, _ := conf["certPath"].(string)
			configPath, _ := conf["configPath"].(string)
			notifier := NewZMQSecNotifier(zmqEndpoint, privateKeyPath, certPath, configPath)
			if enabled {
				notifier.Enable()
			}
			// Token opcional
			if token, ok := conf["authToken"].(string); ok {
				notifier.SetAuthToken(token)
			}
			nm.AddNotifier(name, notifier)

		default:
			// Se o tipo é desconhecido, podemos emitir um log e continuar.
			fmt.Printf("Unknown notifier type '%s' for notifier '%s'\n", typ, name)
			continue
		}
	}
	return nil
}

// Supondo que notifierManager seja uma variável global.
var notifierManager = NewNotifierManager()

func loadConfigAndWatch(port string) error {
	// Configuração de Viper conforme seu código existente.
	cfgPath, cfgType, err := services.GetConfig(port)
	if err != nil {
		return err
	}
	viper.SetConfigFile(cfgPath)
	viper.SetConfigType(cfgType)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("erro ao ler config: %w", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Configuração alterada: %s", e.Name)
		// Atualiza os notifiers a partir da nova configuração.
		if err := notifierManager.UpdateFromConfig(); err != nil {
			log.Printf("Erro ao atualizar notifiers: %v", err)
		} else {
			log.Println("Notifiers atualizados com sucesso.")
		}
	})
	return nil
}
