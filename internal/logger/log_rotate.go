package logger

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CheckLogSize verifica e gerencia o tamanho dos logs
func CheckLogSize(config Config) error {
	logDir := config.DefaultLogPath()
	files, err := os.ReadDir(logDir)
	if err != nil {
		globalLogger.Error("Erro ao ler o diretório de logs", map[string]interface{}{"error": err})
		return err
	}

	// Buscar tamanhos máximos dos logs da configuração
	maxLogSize := config.GetInt("maxLogSize", 20*1024*1024)      // Default 20 MB
	moduleLogSize := config.GetInt("moduleLogSize", 5*1024*1024) // Default 5 MB

	var totalSize int64
	filesToRotate := []string{}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") {
			fileInfo, err := file.Info()
			if err != nil {
				globalLogger.Error("Erro ao obter informações do arquivo", map[string]interface{}{"file": file.Name(), "error": err})
				continue
			}
			totalSize += fileInfo.Size()
			if fileInfo.Size() > int64(moduleLogSize) {
				filesToRotate = append(filesToRotate, filepath.Join(logDir, file.Name()))
			}
		}
	}

	// Rotação com base no tamanho total
	if totalSize > int64(maxLogSize) {
		globalLogger.Info("Tamanho total dos logs excedido. Arquivando logs antigos...", nil)
		if err := ArchiveLogs(filesToRotate); err != nil {
			globalLogger.Error("Erro ao arquivar logs", map[string]interface{}{"error": err})
			return err
		}
	}

	// Rotação individual de arquivos grandes
	if len(filesToRotate) > 0 {
		globalLogger.Info("Arquivando logs individuais devido ao tamanho excessivo...", nil)
		if err := RotateLogFiles(filesToRotate); err != nil {
			globalLogger.Error("Erro ao rotacionar logs", map[string]interface{}{"error": err})
			return err
		}
	}

	return nil
}

// RotateLogFiles compacta e recria os arquivos de log
func RotateLogFiles(files []string) error {
	for _, logFile := range files {
		if err := RotateLogFile(logFile); err != nil {
			globalLogger.Error("Erro ao rotacionar arquivo de log", map[string]interface{}{"file": logFile, "error": err})
			continue
		}
		globalLogger.Info("Arquivo de log rotacionado com sucesso", map[string]interface{}{"file": logFile})
	}
	return nil
}

// RotateLogFile compacta um único arquivo de log
func RotateLogFile(logFilePath string) error {
	archivePath := fmt.Sprintf("%s.tar.gz", logFilePath)
	if err := CreateTarGz(archivePath, []string{logFilePath}); err != nil {
		return err
	}

	if err := os.Remove(logFilePath); err != nil {
		return fmt.Errorf("erro ao remover o arquivo de log: %v", err)
	}

	if err := os.WriteFile(logFilePath, []byte{}, 0644); err != nil {
		return fmt.Errorf("erro ao recriar o arquivo de log: %v", err)
	}
	return nil
}

// CreateTarGz cria um arquivo tar.gz a partir dos logs
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
		if err := addFileToTar(tw, file); err != nil {
			return err
		}
	}

	globalLogger.Info("Arquivo tar.gz criado com sucesso", map[string]interface{}{"path": archivePath})
	return nil
}

// addFileToTar adiciona um arquivo ao arquivo tar
func addFileToTar(tw *tar.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("erro ao criar o cabeçalho do tar: %v", err)
	}
	header.Name = filepath.Base(filePath)

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("erro ao escrever o cabeçalho do tar: %v", err)
	}

	if _, err := io.Copy(tw, file); err != nil {
		return fmt.Errorf("erro ao copiar o conteúdo do arquivo para o tar: %v", err)
	}

	return nil
}

// ArchiveLogs arquiva os logs antigos em um arquivo zip
func ArchiveLogs(files []string) error {
	logDir := GetLogPath()
	if len(files) == 0 {
		err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".log") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("erro ao listar os arquivos de log: %v", err)
		}
	}
	tempDir := os.TempDir()
	archiveName := fmt.Sprintf("logs_archive_%s.zip", time.Now().Format("20060102_150405"))
	archivePath := filepath.Join(tempDir, archiveName)

	zipFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo zip: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return err
		}
	}

	globalLogger.Info("Logs arquivados com sucesso", map[string]interface{}{"archive": archivePath})
	return nil
}

// addFileToZip adiciona um arquivo ao zip
func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo para o zip: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("erro ao criar o cabeçalho do zip: %v", err)
	}
	header.Name = filepath.Base(filePath)

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo no zip: %v", err)
	}

	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("erro ao copiar o conteúdo do arquivo para o zip: %v", err)
	}

	return nil
}

func GetLogDirectorySize(directory string) (int64, error) {
	if directory == "" {
		directory = filepath.Dir(GetLogPath())
	}
	var totalSize int64

	// Percorre o diretório especificado
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("erro ao acessar o caminho %s: %v", path, err)
		}

		// Apenas arquivos são considerados no tamanho total
		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("erro ao calcular o tamanho do diretório: %v", err)
	}

	return totalSize, nil
}
