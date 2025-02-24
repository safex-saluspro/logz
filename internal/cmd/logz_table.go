package cmd

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

// Table representa uma tabela simples.
type Table struct {
	data [][]string
}

// NewTable cria uma nova tabela simples.
func NewTable(data [][]string) Table {
	return Table{data}
}

// PrintTable imprime a tabela simples no shell com bordas laterais e verticais.
func (t Table) PrintTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)  // Ativa bordas laterais e verticais
	table.SetRowLine(true) // Ativa linhas entre as linhas
	for _, row := range t.data {
		table.Append(row)
	}
	table.Render()
}

// FormattedTable representa uma tabela formatada.
type FormattedTable struct {
	data   [][]string
	header []string
}

// NewFormattedTable cria uma nova tabela formatada.
func NewFormattedTable(header []string, data [][]string) FormattedTable {
	return FormattedTable{header: header, data: data}
}

// SaveFormattedTable salva a tabela formatada em um arquivo.
func (ft FormattedTable) SaveFormattedTable(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	table := tablewriter.NewWriter(file)
	table.SetHeader(ft.header)
	for _, row := range ft.data {
		table.Append(row)
	}
	table.Render()
	return nil
}
