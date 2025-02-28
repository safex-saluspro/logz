package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

// LogReader define o contrato para leitura de logs.
type LogReader interface {
	// Tail lê o arquivo de log em tempo real e imprime as novas linhas
	// no terminal. A operação pode ser interrompida enviando um sinal
	// através do canal stopChan.
	Tail(filePath string, stopChan <-chan struct{}) error
}

// FileLogReader implementa a interface LogReader lendo de um arquivo.
type FileLogReader struct{}

// NewFileLogReader cria uma nova instância de FileLogReader.
func NewFileLogReader() *FileLogReader {
	return &FileLogReader{}
}

// Tail segue o arquivo de log a partir do final e imprime as linhas novas conforme elas forem adicionadas.
// O canal stopChan permite interromper a operação (por exemplo, via Ctrl+C em uma sessão interativa).
func (fr *FileLogReader) Tail(filePath string, stopChan <-chan struct{}) error {
	// Abre o arquivo de log
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	// Posiciona o ponteiro no fim do arquivo para começar a ler apenas as novas linhas
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek to the end of the file: %w", err)
	}

	reader := bufio.NewReader(f)
	// Loop principal para ler linhas novas
	for {
		select {
		case <-stopChan:
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// Se não há novas linhas, aguarda um pouco e tenta novamente
					time.Sleep(500 * time.Millisecond)
					continue
				}
				return fmt.Errorf("error reading log file: %w", err)
			}
			fmt.Print(line)
		}
	}
}
