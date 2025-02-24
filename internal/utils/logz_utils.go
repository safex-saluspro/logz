package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/faelmori/kbx/mods/utils"
	lgzCmd "github.com/faelmori/logz/internal/cmd"
	"github.com/spf13/cobra"
	"golang.org/x/exp/rand"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var excludedModuleNames = []string{
	"follow", "show", "clear",
	"--follow", "-f", "--show", "-s", "--clear", "-c",
	"logType", "message", "whatToShow",
	"-type", "-t", "-message", "-m"}

var logColors = []string{
	"\033[30m", // Black
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
	"\033[37m", // White
}

func GetRandomColor() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	return logColors[rand.Intn(len(logColors))]
}

func GetLogLevel(line string) string {
	for level := range lgzCmd.LogLevels {
		if strings.Contains(line, level) {
			return level
		}
	}
	return "INFO"
}

func SetModule(module string) {
	if module == "" {
		if lgzCmd.LogModule == "" {
			lgzCmd.LogModule = "logz"
		}
	} else {
		lgzCmd.LogModule = module
	}
}

func GetLogFileByModule() (string, error) {
	if lgzCmd.LogModule == "" {
		lgzCmd.LogModule = "kbx"
	}
	tempDir, tempDirErr := utils.GetTempDir()
	if tempDirErr != nil {
		return "", tempDirErr
	}
	_logModule := strings.ToUpper(lgzCmd.LogModule)
	_logModuleEnv := fmt.Sprintf("KBX_%s_LOG_FILE", _logModule)
	logFilePath := os.Getenv(_logModuleEnv)
	if logFilePath == "" {
		_logModuleFile := fmt.Sprintf("%s.log", strings.ToLower(lgzCmd.LogModule))
		logFilePath = filepath.Join(tempDir, _logModuleFile)
	}
	ensureFileErr := utils.EnsureFile(logFilePath, 0777, []string{})
	if ensureFileErr != nil {
		fmt.Println("Erro ao garantir o arquivo de log:", ensureFileErr)
		return "", ensureFileErr
	}
	return logFilePath, nil
}

func SetlgzCmdLogLevel() {

	if lgzCmd.LogLevel == "" {
		lgzCmd.LogLevel = os.Getenv("KBX_LOG_LEVEL")
	}
	if lgzCmd.LogLevel == "" {
		lgzCmd.LogLevel = os.Getenv("log_level")
	}
	if lgzCmd.LogLevel == "" {
		lgzCmd.LogLevel = os.Getenv("LOG_LEVEL")
	}
	if lgzCmd.LogLevel == "" {
		lgzCmd.LogLevel = "info"
	}
}

func SetLogOutput(output string) {
	logOutput = output
	if logOutput == "" {
		logOutput = os.Getenv("kbx_log_output")
		if logOutput == "" {
			logOutput = os.Getenv("KBX_LOG_OUTPUT")
		}
		if logOutput == "" {
			logOutput = os.Getenv("log_output")
		}
		if logOutput == "" {
			logOutput = os.Getenv("LOG_OUTPUT")
		}
		if logOutput == "" {
			logOutput = "stdout"
		}
	}
}

func GetLogOutput(quiet bool, output string) string {
	if quiet {
		return "/dev/null"
	}
	if output == "" {
		return "stdout"
	}
	if logOutput == "" {
		return output
	}
	return logOutput
}

func InitLogz(logModuleName string, quietFlag bool) {
	SetModule(logModuleName)
	lgzCmd.setLogLevel()
	if quietFlag {
		SetLogOutput("/dev/null")
	} else {
		SetLogOutput("")
	}
	_, _ = GetLogFileByModule()
}

func CheckLogSize() error {
	tempDir, tempDirErr := utils.GetTempDir()
	if tempDirErr != nil {
		return tempDirErr
	}
	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("erro ao ler o diretório de logs: %v", err)
	}

	// Buscar tamanhos máximos dos logs
	maxLogSizeEnv := os.Getenv("KBX_MAX_LOG_SIZE")
	maxLogSizeVar, err := strconv.ParseInt(maxLogSizeEnv, 10, 64)
	if err != nil || maxLogSizeVar == 0 {
		maxLogSizeVar = 20 * 1024 * 1024 // 20 MB
	}

	moduleLogSizeEnv := os.Getenv(fmt.Sprintf("KBX_%s_LOG_SIZE", strings.ToUpper(lgzCmd.LogModule)))
	moduleLogSize, err := strconv.ParseInt(moduleLogSizeEnv, 10, 64)
	if err != nil || moduleLogSize == 0 {
		moduleLogSize = 5 * 1024 * 1024 // 5 MB
	}

	var totalSize int64
	filesToRotateList := make([]string, 0)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") && (strings.HasPrefix(file.Name(), "kb") || file.Name() == "logz.log") {
			fileInfo, fileInfoErr := file.Info()
			if fileInfoErr != nil {
				return fmt.Errorf("erro ao obter informações do arquivo de log: %v", fileInfoErr)
			}
			totalSize += fileInfo.Size()

			// Verificar tamanho individual do arquivo
			if fileInfo.Size() > moduleLogSize {
				filesToRotateList = append(filesToRotateList, file.Name())
			}
		}
	}

	// Verificar tamanho total dos arquivos
	if totalSize > maxLogSizeVar {
		allLogFilesList := make([]string, 0)
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".log") && (strings.HasPrefix(file.Name(), "kb") || file.Name() == "logz.log") {
				allLogFilesList = append(allLogFilesList, file.Name())
			}
		}
		rotateAllLogFilesErr := rotateAllLogFiles(allLogFilesList)
		if rotateAllLogFilesErr != nil {
			return rotateAllLogFilesErr
		}
	}

	if len(filesToRotateList) > 0 {
		rotateLogFilesErr := rotateAllLogFiles(filesToRotateList)
		if rotateLogFilesErr != nil {
			return rotateLogFilesErr
		}
	}

	return nil
}

func RotateLogFile(logFilePath string) error {
	archivePath := fmt.Sprintf("%s.tar.gz", logFilePath)
	err := createTarGz(archivePath, []string{logFilePath})
	if err != nil {
		return err
	}

	err = os.Remove(logFilePath)
	if err != nil {
		return fmt.Errorf("erro ao remover o arquivo de log: %v", err)
	}

	err = utils.EnsureFile(logFilePath, 0777, []string{})
	if err != nil {
		return fmt.Errorf("erro ao recriar o arquivo de log: %v", err)
	}

	return nil
}

func RotateAllLogFiles(filesList []string) error {
	hasDifferentDirs := false
	filesDirList := make(map[string][]string)
	for _, file := range filesList {
		fileDir := filepath.Dir(file)
		existsDirInList := utils.Contains(filesDirList, fileDir)
		if !existsDirInList {
			filesDirList[fileDir] = []string{file}
		} else {
			filesDirList[fileDir] = append(filesDirList[fileDir], file)
		}
	}
	if len(filesDirList) > 1 {
		hasDifferentDirs = true
	}
	filesEntries := make([]os.DirEntry, 0)
	for dir, files := range filesDirList {
		dirEntries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("erro ao ler o diretório de logs: %v", err)
		}
		for _, file := range files {
			for _, dirEntry := range dirEntries {
				if dirEntry.Name() == file || strings.HasPrefix(dirEntry.Name(), file) {
					filesEntries = append(filesEntries, dirEntry)
				}
			}
		}
	}
	var logFiles []string
	for _, file := range filesEntries {
		logFiles = append(logFiles, file.Name())
	}

	var archivePath string
	if !hasDifferentDirs {
		archivePath = fmt.Sprintf("%s.tar.gz", filesEntries[0].Name())
		createTarGzErr := createTarGz(archivePath, logFiles)
		if createTarGzErr != nil {
			return createTarGzErr
		}
	} else {
		for dir, files := range filesDirList {
			archivePath = fmt.Sprintf("%s.tar.gz", dir)
			createTarGzErr := createTarGz(archivePath, files)
			if createTarGzErr != nil {
				return createTarGzErr
			}
		}
	}

	for _, logFile := range logFiles {
		err := os.Remove(logFile)
		if err != nil {
			return fmt.Errorf("erro ao remover o arquivo de log: %v", err)
		}

		err = utils.EnsureFile(logFile, 0777, []string{})
		if err != nil {
			return fmt.Errorf("erro ao recriar o arquivo de log: %v", err)
		}
	}

	return nil
}

func CreateTarGz(archivePath string, files []string) error {
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo tar.gz: %v", err)
	}
	defer archiveFile.Close()

	gw := gzip.NewWriter(archiveFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		err := addFileToTar(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addFileToTar(tw *tar.Writer, filePath string) error {
	tempDir, tempDirErr := utils.GetTempDir()
	if tempDirErr != nil {
		return tempDirErr
	}

	// Validar e sanitizar o caminho do arquivo
	cleanPath := filepath.Clean(filePath)
	if !strings.HasPrefix(cleanPath, tempDir) {
		return fmt.Errorf("caminho do arquivo inválido: %s", filePath)
	}

	file, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo para o tar: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return fmt.Errorf("erro ao criar o cabeçalho do tar: %v", err)
	}
	header.Name = filepath.Base(cleanPath)

	err = tw.WriteHeader(header)
	if err != nil {
		return fmt.Errorf("erro ao escrever o cabeçalho do tar: %v", err)
	}

	_, err = io.Copy(tw, file)
	if err != nil {
		return fmt.Errorf("erro ao copiar o arquivo para o tar: %v", err)
	}

	return nil
}

func ArchiveOldLogs() error {
	tempDir, tempDirErr := utils.GetTempDir()
	if tempDirErr != nil {
		return tempDirErr
	}

	archiveName := fmt.Sprintf("logs_archive_%s.zip", time.Now().Format("20060102_150405"))
	archivePath := filepath.Join(tempDir, archiveName)

	zipFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo zip: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("erro ao ler o diretório de logs: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") {
			filePath := filepath.Join(tempDir, file.Name())
			err := addFileToZip(zipWriter, filePath)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Logs arquivados com sucesso em:", archivePath)
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("erro ao criar o cabeçalho do arquivo zip: %v", err)
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("erro ao criar o cabeçalho do arquivo zip: %v", err)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("erro ao copiar o arquivo para o zip: %v", err)
	}

	return nil
}

func DetectSystemLanguage() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		// Se não conseguir, usa um padrão
		lang = "en-US"
	}
	if strings.Contains(lang, ".") {
		lang = strings.Split(lang, ".")[0]
	}
	return lang
}

func GetMessageFlag(cmd cobra.Command, args []string) string {
	message, _ := cmd.Flags().GetString("message")
	if message != "" {
		return message
	}
	hasTypeFlag := false
	if tf, _ := cmd.Flags().GetString("type"); tf != "" {
		hasTypeFlag = true
	}
	messageFlagIndex := 1
	if hasTypeFlag {
		messageFlagIndex = 0
	}
	if len(args) > messageFlagIndex {
		return args[messageFlagIndex]
	}
	return "no message"
}

func GetModuleNameFlag(cmd cobra.Command, args []string) string {
	logModuleName, _ := cmd.Flags().GetString("name")
	if logModuleName != "" {
		return logModuleName
	}
	hasTypeFlag := false
	if tf, _ := cmd.Flags().GetString("type"); tf != "" {
		hasTypeFlag = true
	}
	hasMessageFlag := false
	if mf, _ := cmd.Flags().GetString("message"); mf != "" {
		hasMessageFlag = true
	}
	moduleFlagIndex := 0
	if hasTypeFlag {
		if hasMessageFlag {
			moduleFlagIndex = 2
		} else {
			moduleFlagIndex = 1
		}
	} else {
		if hasMessageFlag {
			moduleFlagIndex = 1
		} else {
			moduleFlagIndex = 2
		}
	}
	if len(args) > moduleFlagIndex {
		return args[moduleFlagIndex]
	} else {
		return ""
	}
}

func GetLogTypeFlag(cmd cobra.Command, args []string) string {
	logType, _ := cmd.Flags().GetString("type")
	if logType != "" {
		return logType
	}
	if len(args) > 0 {
		return args[0]
	}
	return "info"
}

func GetShowFlag(cmd cobra.Command, args []string) string {
	showFlag, _ := cmd.Flags().GetString("show")
	if showFlag != "" {
		return showFlag
	}
	if len(args) > 0 {
		containsShow := utils.Contains(args, "show")
		if containsShow {
			for i, arg := range args {
				if arg == "show" {
					if len(args) >= i+1 {
						return args[i+1]
					} else {
						return "all"
					}
				}
			}
		} else {
			return ""
		}
	} else {
		return ""
	}
	return ""
}

func GetClearFlag(cmd cobra.Command, args []string) string {
	clearFlag, _ := cmd.Flags().GetString("clear")
	if clearFlag != "" {
		return clearFlag
	}
	if len(args) > 0 {
		containsClear := utils.Contains(args, "clear")
		if containsClear {
			for i, arg := range args {
				if arg == "clear" {
					if len(args) >= i+1 {
						return args[i+1]
					} else {
						return "all"
					}
				}
			}
		} else {
			return ""
		}
	} else {
		return ""
	}
	return ""
}

func GetArchiveFlag(cmd cobra.Command, args []string) string {
	archiveFlag, _ := cmd.Flags().GetString("archive")
	if archiveFlag != "" {
		return archiveFlag
	}
	if len(args) > 0 {
		containsArchive := utils.Contains(args, "archive")
		if containsArchive {
			for i, arg := range args {
				if arg == "archive" {
					if len(args) >= i+1 {
						return args[i+1]
					} else {
						return "all"
					}
				}
			}
		} else {
			return ""
		}
	} else {
		return ""
	}
	return ""
}

func CheckFollowFlag() (bool, error) {
	if len(os.Args) < 2 {
		return false, fmt.Errorf("não foi informado nenhum argumento para o logger %v\n", logModule)
	}
	for _, arg := range os.Args[2:] {
		if arg == "-f" || arg == "--follow" || arg == "follow" {
			return true, nil
		}
	}
	return false, nil
}

func CheckQuietFlag() (bool, error) {
	if len(os.Args) < 2 {
		return false, fmt.Errorf("não foi informado nenhum argumento para o logger %v\n", logModule)
	}
	for _, arg := range os.Args[2:] {
		if arg == "-q" || arg == "--quiet" || arg == "quiet" {
			return true, nil
		}
	}
	return false, nil
}
