package utils

import (
	"bufio"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

func checkLogExists() bool {
	logFilePath, logFilePathErr := GetLogFileByModule()
	if logFilePathErr != nil {
		return false
	}

	// Verifica a existência do arquivo de log sem usar sudo
	if _, err := os.Stat(logFilePath); err == nil {
		return true
	} else {
		return false
	}
}

func PrintLogFile(logFilePath string) ([]string, error) {
	if !checkLogExists() {
		return nil, fmt.Errorf("arquivo de log não encontrado")
	}
	var logFileMessages = make([]string, 0)

	file, fileErr := os.Open(logFilePath)
	if fileErr != nil {
		return nil, fmt.Errorf("160: erro ao abrir o arquivo de log: %v", fileErr)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lang := DetectSystemLanguage()
	tag, err := language.Parse(lang)
	if err != nil {
		fmt.Println("Erro ao detectar a linguagem do sistema, usando en-US.")
		tag = language.AmericanEnglish
	}
	p := message.NewPrinter(tag)

	var previousLine string
	for {
		line, lineErr := reader.ReadString('\n')
		if lineErr != nil && lineErr != io.EOF {
			return nil, fmt.Errorf("162: erro ao ler a linha do arquivo de log: %v", lineErr)
		}
		if lineErr == io.EOF {
			break
		}

		parts := strings.SplitN(line, " ", 5)
		if len(parts) < 5 {
			// Linha sem timestamp, adiciona ao log anterior
			previousLine += "                    " + line
			continue
		}

		if previousLine != "" {
			logFileMessages = append(logFileMessages, previousLine)
			fmt.Print(previousLine)
		}

		var timestamp string
		rawTimestamp := strings.TrimSuffix(strings.TrimPrefix(strings.Join([]string{parts[0], parts[1]}, " "), "["), "]") // Pega o timestamp original
		dateParsed, dateParsedErr := time.ParseInLocation("2006-01-02 15:04:05", rawTimestamp, time.Local)
		if dateParsedErr != nil {
			timestamp = p.Sprintf("%s", rawTimestamp)
		} else {
			timestamp = p.Sprintf("%s", dateParsed)
		}

		logFileBaseName := filepath.Base(logFilePath)
		logFileBaseName = strings.TrimSuffix(logFileBaseName, ".log")
		coloredLine := replaceTypeWithColor(strings.Join(parts[3:], " "))
		previousLine = fmt.Sprintf("[%s] [%s] %s", logFileBaseName, timestamp, coloredLine)
	}

	if previousLine != "" {
		logFileMessages = append(logFileMessages, previousLine)
		fmt.Print(previousLine)
	}

	return logFileMessages, nil
}

func filterLogsByLevel(logFilePath string, level string) error {
	if !checkLogExists() {
		return fmt.Errorf("arquivo de log não encontrado")
	}
	file, fileErr := os.Open(logFilePath)
	if fileErr != nil {
		return fmt.Errorf("610: erro ao abrir o arquivo de log: %v", fileErr)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, lineErr := reader.ReadString('\n')
		if lineErr != nil {
			break
		}

		runTimeLogLevel := GetLogLevel(line)
		if LogLevels[runTimeLogLevel] >= LogLevels[runTimeLogLevel] {
			fmt.Println(line)
		}
	}

	return nil
}

func replaceTypeWithColor(logText string) string {
	logText = strings.ReplaceAll(logText, "[info]", "["+LogColor["info"]+"INFO"+LogColor["default"]+"]")
	logText = strings.ReplaceAll(logText, "[warn]", "["+LogColor["warn"]+"WARN"+LogColor["default"]+"]")
	logText = strings.ReplaceAll(logText, "[error]", "["+LogColor["error"]+"ERROR"+LogColor["default"]+"]")
	logText = strings.ReplaceAll(logText, "[debug]", "["+LogColor["debug"]+"DEBUG"+LogColor["default"]+"]")
	logText = strings.ReplaceAll(logText, "[success]", "["+LogColor["success"]+"SUCCESS"+LogColor["default"]+"]")
	logText = strings.ReplaceAll(logText, "[answer]", "["+LogColor["answer"]+"ANSWER"+LogColor["default"]+"]")
	return logText
}

func FollowAllLogFiles(logFiles []string) error {
	var wg sync.WaitGroup
	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	fileColors := make(map[string]string)
	for _, logFile := range logFiles {
		fileColors[logFile] = GetRandomColor()
	}

	for _, logFile := range logFiles {
		wg.Add(1)
		go func(logFile string) {
			defer wg.Done()
			err := followLogFileWithPrefix(logFile, fileColors[logFile])
			if err != nil {
				fmt.Println("Erro ao seguir o arquivo de log:", err)
			}
		}(logFile)
	}

	go func() {
		<-sigChan
		fmt.Println("Recebido sinal de interrupção, encerrando...")
		os.Exit(0)
	}()

	wg.Wait()
	return nil
}

func followLogFileWithPrefix(logFilePath string, fileColor string) error {
	if !checkLogExists() {
		return fmt.Errorf("arquivo de log não encontrado")
	}
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("161: erro ao abrir o arquivo de log: %v", err)
	}
	defer file.Close()

	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("erro ao buscar o final do arquivo de log: %v", err)
	}
	reader := bufio.NewReader(file)

	var previousLine string
	for {
		line, lineErr := reader.ReadString('\n')
		if lineErr != nil && lineErr != io.EOF {
			return fmt.Errorf("164: erro ao ler a linha do arquivo de log: %v", lineErr)
		}
		if lineErr == io.EOF {
			time.Sleep(time.Second / 3)
			continue
		}

		parts := strings.SplitN(line, " ", 5)
		if len(parts) < 5 {
			// Linha sem timestamp, adiciona ao log anterior
			previousLine += "                    " + line
			continue
		}

		if previousLine != "" {
			fmt.Print(previousLine)
		}

		timestamp := parts[0] + " " + parts[1] // Pega o timestamp original
		if len(timestamp) >= 21 {
			shortTimestamp := timestamp[12:20] // Exibe somente o horário
			logFileBaseName := filepath.Base(logFilePath)
			logFileBaseName = strings.TrimSuffix(logFileBaseName, ".log")
			coloredLine := replaceTypeWithColor(strings.Join(parts[3:], " "))
			coloredFile := fileColor + logFileBaseName + LogColor["default"]
			previousLine = fmt.Sprintf("[%s] [%s] %s", coloredFile, shortTimestamp, coloredLine)
		} else {
			previousLine = line
		}
	}
}

func LogToFile(message string) error {
	logFilePath, logFilePathErr := GetLogFileByModule()
	if logFilePathErr != nil {
		return logFilePathErr
	}

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Could not create file: %s", logFilePath))
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, writeStringErr := file.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
	if writeStringErr != nil {
		return writeStringErr
	}

	go func() {
		checkLogSizeErr := CheckLogSize()
		if checkLogSizeErr != nil {
			fmt.Println("Erro ao verificar o tamanho do log:", checkLogSizeErr)
		}
	}()

	return nil
}

func GetLogLevel(line string) string {
	for level := range LogLevels {
		if strings.Contains(line, level) {
			return level
		}
	}
	return "INFO"
}
