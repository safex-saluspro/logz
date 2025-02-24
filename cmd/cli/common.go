package cli

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/faelmori/kbx/mods/ui/wrappers"
	"os"
	"strings"
)

var (
	outputType   string
	outputTarget string
)

func RunWithLoader(task func(chan tea.Msg) error) error {
	messages := make(chan tea.Msg)

	go func() {
		defer close(messages)
		if err := task(messages); err != nil {
			messages <- KbdzLoaderMsg{Message: "Error: " + err.Error()}
		}
		messages <- KbdzLoaderCloseMsg{}
	}()

	return StartLoader(messages)
}

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	var description, banner string

	if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
		description = descriptionArg[0]
	} else {
		description = descriptionArg[1]
	}

	if !hideBanner {
		banner = ` ____       ____  ____                 
 / ___|     |  _ \| __ )  __ _ ___  ___ 
| |  _ _____| | | |  _ \ / _| / __|/ _ \
| |_| |_____| |_| | |_) | (_| \__ \  __/
 \____|     |____/|____/ \__,_|___/\___|
`
	} else {
		banner = ""
	}
	return map[string]string{"banner": banner, "description": description}
}
