package cmd

import (
	"bufio"
	"fmt"
	"github.com/faelmori/logz/internal/extras"
	"github.com/go-echarts/go-echarts/charts"
	"os"
	"strings"
)

func collectLogMetrics(logFilePath string) (*LogMetrics, error) {
	if !checkLogExists() {
		_ = extras.ErrorLog("Arquivo de log não encontrado")
		return nil, nil
	}
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, fmt.Errorf("312: erro ao abrir o arquivo de log: %v", err)
	}
	defer file.Close()

	metrics := &LogMetrics{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Info") {
			metrics.InfoCount++
		} else if strings.Contains(line, "Warn") {
			metrics.WarnCount++
		} else if strings.Contains(line, "Error") {
			metrics.ErrorCount++
		} else if strings.Contains(line, "Debug") {
			metrics.DebugCount++
		} else if strings.Contains(line, "Success") {
			metrics.SuccessCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("300: erro ao ler o arquivo de log: %v", err)
	}

	return metrics, nil
}

func generateLogReport(metrics *LogMetrics) {
	fmt.Println("Relatório de Logs:")
	fmt.Printf("Info: %d\n", metrics.InfoCount)
	fmt.Printf("Warn: %d\n", metrics.WarnCount)
	fmt.Printf("Error: %d\n", metrics.ErrorCount)
	fmt.Printf("Debug: %d\n", metrics.DebugCount)
	fmt.Printf("Success: %d\n", metrics.SuccessCount)
}

func generateLogChart(metrics *LogMetrics) error {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Log Metrics"})

	bar.AddXAxis([]string{"Info", "Warn", "Error", "Debug", "Success"}).
		AddYAxis("Count", []int{metrics.InfoCount, metrics.WarnCount, metrics.ErrorCount, metrics.DebugCount, metrics.SuccessCount})

	f, err := os.Create("log_metrics.html")
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo de gráfico: %v", err)
	}
	defer f.Close()

	return bar.Render(f)
}

func analyzeLog(logFilePath string) error {
	metrics, err := collectLogMetrics(logFilePath)
	if err != nil {
		return err
	}

	generateLogReport(metrics)
	err = generateLogChart(metrics)
	if err != nil {
		return err
	}

	fmt.Println("Relatório de logs gerado com sucesso em: log_metrics.html")
	return nil
}
