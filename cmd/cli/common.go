package cli

import (
	"os"
	"strings"
)

//var (
//	outputType   string
//	outputTarget string
//)
//
//func RunWithLoader(task func(chan tea.Msg) error) error {
//	messages := make(chan tea.Msg)
//
//	go func() {
//		defer close(messages)
//		if err := task(messages); err != nil {
//			messages <- KbdzLoaderMsg{Message: "Error: " + err.Error()}
//		}
//		messages <- KbdzLoaderCloseMsg{}
//	}()
//
//	return StartLoader(messages)
//}

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	var description, banner string

	if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
		description = descriptionArg[0]
	} else {
		description = descriptionArg[1]
	}

	if !hideBanner {
		banner = ` _                    
 | |    ___   __ _ ____
 | |   / _ \ / _\ |_  /
 | |__| (_) | (_| |/ / 
 |_____\___/ \__, /___|
             |___/     
`
	} else {
		banner = ""
	}
	return map[string]string{"banner": banner, "description": description}
}
