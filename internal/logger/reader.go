package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// LogReader define o contrato para leitura de logs.
type LogReader interface {
	// Tail lê o arquivo de log em tempo real e envia as novas linhas para o
	// writer passado ou imprime-as no terminal. A operação pode ser interrompida
	// enviando um sinal através do canal stopChan.
	Tail(filePath string, stopChan <-chan struct{}) error
}

// FileLogReader implementa a interface LogReader lendo de um arquivo.
type FileLogReader struct {
	// pollInterval é o intervalo de polling para verificar novas linhas.
	pollInterval time.Duration
}

// NewFileLogReader cria uma nova instância de FileLogReader.
// O intervalo de polling é lido da variável de ambiente LOGZ_TAIL_POLL_INTERVAL (em milissegundos),
// ou usa 500ms por padrão.
func NewFileLogReader() *FileLogReader {
	intervalMs := 500 // padrão em milissegundos
	if val := os.Getenv("LOGZ_TAIL_POLL_INTERVAL"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			intervalMs = parsed
		} else {
			log.Printf("Invalid LOGZ_TAIL_POLL_INTERVAL value, using default 500ms: %v", err)
		}
	}
	return &FileLogReader{
		pollInterval: time.Duration(intervalMs) * time.Millisecond,
	}
}

// Tail segue o arquivo de log a partir do fim e imprime as novas linhas conforme elas são adicionadas.
// O canal stopChan permite interromper a operação (por exemplo, via Ctrl+C).
func (fr *FileLogReader) Tail(filePath string, stopChan <-chan struct{}) error {
	// Abre o arquivo de log
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	// Posiciona o ponteiro no final do arquivo para ler apenas novas linhas
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek to the end of the file: %w", err)
	}

	reader := bufio.NewReader(f)

	// Loop principal para ler linhas novas
	for {
		select {
		case <-stopChan:
			log.Println("Tail operation interrupted by stop signal")
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					time.Sleep(fr.pollInterval)
					continue
				}
				return fmt.Errorf("error reading log file: %w", err)
			}
			// Imprime a linha imediatamente; pode ser adaptado para enviar para outro canal se necessário.
			fmt.Print(line)
		}
	}
}
