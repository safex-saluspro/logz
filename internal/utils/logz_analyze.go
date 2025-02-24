package utils

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/charts"
	"os"
	"strings"
)

func CollectLogMetrics(logFilePath string) (*LogMetrics, error) {
	if !checkLogExists() {
		return nil, fmt.Errorf("arquivo de log não encontrado")
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

func GenerateLogReport(metrics *LogMetrics) {
	fmt.Println("Relatório de Logs:")
	fmt.Printf("Info: %d\n", metrics.InfoCount)
	fmt.Printf("Warn: %d\n", metrics.WarnCount)
	fmt.Printf("Error: %d\n", metrics.ErrorCount)
	fmt.Printf("Debug: %d\n", metrics.DebugCount)
	fmt.Printf("Success: %d\n", metrics.SuccessCount)
}

func GenerateLogChart(metrics *LogMetrics) error {
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
