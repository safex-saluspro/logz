package cmd

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/journal"
	"github.com/faelmori/kbx/mods/utils"
	lgzUtl "github.com/faelmori/logz/internal/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	ANSWER = "answer"
	STDOUT = "stdout"
	ALL    = "all"
	ERROR  = "error"
	QUIET  = "quiet"
	TRUE   = "true"
	FALSE  = "false"
)

type Logz interface {
	Log(logArgs ...string) error
	ShowLog(args ...string) ([]string, error)
	ClearLogs(whatToClear string) error
	AnswerLog(logArgs ...string) error
	InfoLog(logArgs ...string) error
	WarnLog(logArgs ...string) error
	ErrorLog(logArgs ...string) error
	DebugLog(logArgs ...string) error
	SuccessLog(logArgs ...string) error
	Panic(args ...interface{})
	Writer(module string) io.Writer
	AnalyzeLog(logFilePath string) error
}
type logzImpl struct{}

func (l *logzImpl) Log(logArgs ...string) error {
	if len(logArgs) < 2 {
		return fmt.Errorf("invalid log arguments")
	}

	logType := ""
	message := ""
	logModuleName := ""
	quietFlag := false

	logType = logArgs[0]
	message = logArgs[1]
	if len(logArgs) > 2 {
		logModuleName = logArgs[2]
	}
	if len(logArgs) > 3 {
		quietFlag = logArgs[3] == TRUE || logArgs[3] == QUIET
	} else {
		quietFlag, _ = lgzUtl.CheckQuietFlag()
	}

	lgzUtl.InitLogz(logModuleName, quietFlag)

	color, exists := lgzUtl.LogColor[logType]
	if !exists {
		color = lgzUtl.LogColor["default"]
	}

	var logMessage string
	var loggerMessage string
	var willLog bool

	if logType != ANSWER {
		if !quietFlag {
			logMessage = fmt.Sprintf("[%s] %s%s%s", lgzUtl.LogModule, color, message, lgzUtl.LogColor["default"])
		}
		loggerMessage = fmt.Sprintf("[%s] [%s] %s", lgzUtl.LogModule, logType, message)
		willLog = true
	} else {
		if !quietFlag {
			logMessage = fmt.Sprint(message)
		}
		loggerMessage = fmt.Sprintf("[%s] [%s] %s", lgzUtl.LogModule, logType, message)
		willLog = false
	}

	switch lgzUtl.LogOutput {
	case STDOUT:
		if !quietFlag {
			fmt.Println(logMessage)
		}
	case "stderr":
		_, fPrintlnErr := fmt.Fprintln(os.Stderr, logMessage)
		if fPrintlnErr != nil {
			fmt.Println("Erro ao enviar mensagem para o stderr:", fPrintlnErr)
			return fPrintlnErr
		}
		stdErrErr := os.Stderr.Sync()
		if stdErrErr != nil {
			fmt.Println("Erro ao sincronizar o stderr:", stdErrErr)
			return stdErrErr
		}
	case "journal":
		sendToJournalErr := journal.Send(logMessage, journal.PriInfo, nil)
		if sendToJournalErr != nil {
			fmt.Println("Erro ao enviar mensagem para o journal:", sendToJournalErr)
		}
		willLog = false
	default:
		if !quietFlag {
			fmt.Println(logMessage)
		}
	}

	err := lgzUtl.CheckLogSize()
	if err != nil {
		fmt.Println("Erro ao enviar mensagem para o journal:", err)
		return err
	}

	if !willLog {
		return nil
	}

	if logType == ERROR {
		logToFileErr := lgzUtl.LogToFile(loggerMessage)
		if logToFileErr != nil {
			return logToFileErr
		}
		return fmt.Errorf("%s", message)
	}
	return lgzUtl.LogToFile(loggerMessage)
}
func (l *logzImpl) ShowLog(args ...string) ([]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("invalid number of arguments")
	}

	module := args[0]
	tempDir, tempDirErr := utils.GetTempDir()
	if tempDirErr != nil {
		return nil, tempDirErr
	}
	var logFiles []string

	if module == ALL {
		files, filesErr := os.ReadDir(tempDir)
		if filesErr != nil {
			return nil, fmt.Errorf("erro ao ler o diretório de logs: %w", filesErr)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".log") {
				logFiles = append(logFiles, filepath.Join(tempDir, file.Name()))
			}
		}
	} else {
		logFile := filepath.Join(tempDir, module+".log")
		logFiles = append(logFiles, logFile)
	}

	if len(logFiles) == 0 {
		_ = l.ErrorLog("Nenhum arquivo de log encontrado", "kbx")
		return nil, nil
	}
	var logFileErr error
	var logFileMessages []string
	if follow, _ := lgzUtl.CheckFollowFlag(); follow {
		return nil, lgzUtl.FollowAllLogFiles(logFiles)
	} else {

		if len(logFiles) == 1 {
			logFileMessages, logFileErr = lgzUtl.PrintLogFile(logFiles[0])
			if logFileErr != nil {
				return nil, logFileErr
			}
		} else {
			for _, logFile := range logFiles {
				logFileMessages, logFileErr = lgzUtl.PrintLogFile(logFile)
				if logFileErr != nil {
					return nil, logFileErr
				}
			}
		}
	}
	return logFileMessages, nil
}
func (l *logzImpl) ClearLogs(whatToClear string) error {
	tempDir, tempDirErr := utils.GetTempDir()
	if tempDirErr != nil {
		return tempDirErr
	}

	var logsToClear []string

	files, filesErr := os.ReadDir(tempDir)
	if filesErr != nil {
		return fmt.Errorf("erro ao ler o diretório de logs: %w", filesErr)
	}

	force, _ := lgzUtl.CheckFollowFlag()
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") {
			willClear := false
			fileWithoutExtension := strings.ToLower(strings.TrimSuffix(file.Name(), ".log"))
			if whatToClear == fileWithoutExtension || whatToClear == file.Name() {
				willClear = true
			} else if whatToClear == ALL && force {
				willClear = true
			} else if whatToClear == ALL {
				if !strings.HasPrefix(file.Name(), "kbx") {
					willClear = true
				}
			}
			if willClear {
				logsToClear = append(logsToClear, filepath.Join(tempDir, file.Name()))
			}
		}
	}

	for _, logFile := range logsToClear {
		cmdRemoveLogUser := exec.Command("rm", "-f", logFile)
		cmdRemoveLogUserErr := cmdRemoveLogUser.Run()
		if cmdRemoveLogUserErr != nil {
			cmdRemoveLogRoot := exec.Command("sudo", "rm", "-f", logFile)
			cmdRemoveLogRootErr := cmdRemoveLogRoot.Run()
			if cmdRemoveLogRootErr != nil {
				return fmt.Errorf("erro ao remover o arquivo de log: %w", cmdRemoveLogRootErr)
			}
		}
	}

	if whatToClear == ALL {
		if force {
			fmt.Println("Todos os logs removidos com sucesso")
		} else {
			_ = l.AnswerLog("Todos os logs removidos com sucesso", "kbx")
		}
	} else {
		if whatToClear == "kbx" {
			fmt.Println("Log kbx removido com sucesso")
		} else {
			_ = l.AnswerLog("Log "+whatToClear+" removido com sucesso", "kbx")
		}
	}

	return nil
}
func (l *logzImpl) AnswerLog(logArgs ...string) error {
	module := ""
	if len(logArgs) > 1 {
		module = logArgs[1]
	}
	quiet := FALSE
	blIsQuiet := logArgs[len(logArgs)-1] != module && logArgs[len(logArgs)-1] == QUIET
	if blIsQuiet || logArgs[len(logArgs)-1] == TRUE || logArgs[len(logArgs)-1] == FALSE {
		lgzUtl.SetLogOutput(logArgs[len(logArgs)-1])
		quiet = fmt.Sprintf("%t", logArgs[len(logArgs)-1] == QUIET)
	}
	return l.Log(ANSWER, logArgs[0], module, quiet)
}
func (l *logzImpl) InfoLog(logArgs ...string) error {
	module := ""
	if len(logArgs) > 1 {
		module = logArgs[1]
	}
	quiet := FALSE
	blIsQuiet := logArgs[len(logArgs)-1] != module && logArgs[len(logArgs)-1] == QUIET
	if blIsQuiet || logArgs[len(logArgs)-1] == TRUE || logArgs[len(logArgs)-1] == FALSE {
		lgzUtl.SetLogOutput(logArgs[len(logArgs)-1])
		quiet = fmt.Sprintf("%t", logArgs[len(logArgs)-1] == QUIET)
	}
	return l.Log("info", logArgs[0], module, quiet)
}
func (l *logzImpl) WarnLog(logArgs ...string) error {
	module := ""
	if len(logArgs) > 1 {
		module = logArgs[1]
	}
	quiet := FALSE
	blIsQuiet := logArgs[len(logArgs)-1] != module && logArgs[len(logArgs)-1] == QUIET
	if blIsQuiet || logArgs[len(logArgs)-1] == TRUE || logArgs[len(logArgs)-1] == FALSE {
		lgzUtl.SetLogOutput(logArgs[len(logArgs)-1])
		quiet = fmt.Sprintf("%t", logArgs[len(logArgs)-1] == QUIET)
	}
	return l.Log("warn", logArgs[0], module, quiet)
}
func (l *logzImpl) ErrorLog(logArgs ...string) error {
	module := ""
	if len(logArgs) > 1 {
		module = logArgs[1]
	}
	quiet := FALSE
	blIsQuiet := logArgs[len(logArgs)-1] != module && logArgs[len(logArgs)-1] == QUIET
	if blIsQuiet || logArgs[len(logArgs)-1] == TRUE || logArgs[len(logArgs)-1] == FALSE {
		lgzUtl.SetLogOutput(logArgs[len(logArgs)-1])
		quiet = fmt.Sprintf("%t", logArgs[len(logArgs)-1] == QUIET)
	}
	return l.Log(ERROR, logArgs[0], module, quiet)
}
func (l *logzImpl) DebugLog(logArgs ...string) error {
	module := ""
	if len(logArgs) > 1 {
		module = logArgs[1]
	}
	quiet := FALSE
	blIsQuiet := logArgs[len(logArgs)-1] != module && logArgs[len(logArgs)-1] == QUIET
	if blIsQuiet || logArgs[len(logArgs)-1] == TRUE || logArgs[len(logArgs)-1] == FALSE {
		lgzUtl.SetLogOutput(logArgs[len(logArgs)-1])
		quiet = fmt.Sprintf("%t", logArgs[len(logArgs)-1] == QUIET)
	}
	return l.Log("debug", logArgs[0], module, quiet)
}
func (l *logzImpl) SuccessLog(logArgs ...string) error {
	module := ""
	if len(logArgs) > 1 {
		module = logArgs[1]
	}
	quiet := FALSE
	blIsQuiet := logArgs[len(logArgs)-1] != module && logArgs[len(logArgs)-1] == QUIET
	if blIsQuiet || logArgs[len(logArgs)-1] == TRUE || logArgs[len(logArgs)-1] == FALSE {
		lgzUtl.SetLogOutput(logArgs[len(logArgs)-1])
		quiet = fmt.Sprintf("%t", logArgs[len(logArgs)-1] == QUIET)
	}
	return l.Log("success", logArgs[0], module, quiet)
}
func (l *logzImpl) Panic(args ...interface{}) {
	_ = fmt.Errorf("panic: %s", fmt.Sprint(args...))
	panic(fmt.Sprint(args...))
}
func (l *logzImpl) Writer(module string) io.Writer {
	lgzUtl.SetModule(module)
	logFilePath, logFilePathErr := lgzUtl.GetLogFileByModule()
	if logFilePathErr != nil {
		fmt.Println(logFilePathErr)
		return os.Stdout
	}
	logFile, logFileErr := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if logFileErr != nil {
		fmt.Println(logFileErr)
		return os.Stdout
	}
	return logFile
}
func (l *logzImpl) AnalyzeLog(logFilePath string) error {
	metrics, err := lgzUtl.CollectLogMetrics(logFilePath)
	if err != nil {
		return err
	}

	lgzUtl.GenerateLogReport(metrics)
	err = lgzUtl.GenerateLogChart(metrics)
	if err != nil {
		return err
	}

	fmt.Println("Relatório de logs gerado com sucesso em: log_metrics.html")
	return nil
}

func NewLogz() Logz { return &logzImpl{} }
