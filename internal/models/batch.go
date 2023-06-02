package models

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

const maxLines = 20

type result struct {
	item string
	err  error
}

func (r result) String() string {
	var status string
	if r.err != nil {
		status = "NOK"
	} else {
		status = "OK"
	}
	return fmt.Sprintf("%-50s%s", r.item, status)
}

type batchModel struct {
	message     string
	spinner     spinner.Model
	progress    progress.Model
	items       []string
	currentItem int
	results     []result
	done        bool
	outputFile  io.StringWriter
	behavior    func(string) tea.Cmd
}

func (b batchModel) Init() tea.Cmd {
	return tea.Batch(
		b.behavior(b.items[b.currentItem]),
		b.spinner.Tick,
	)
}

func (b batchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progress.FrameMsg:
		progressModel, cmd := b.progress.Update(msg)
		b.progress = progressModel.(progress.Model)
		return b, cmd
	case spinner.TickMsg:
		var cmd tea.Cmd
		b.spinner, cmd = b.spinner.Update(msg)
		return b, cmd
	case message:
		err := msg.write(b.outputFile)
		if err != nil {
			return b, tea.Quit
		}

		if len(b.results) > maxLines {
			b.results = append(b.results[1:], msg.toResult())
		} else {
			b.results = append(b.results, msg.toResult())
		}
		b.currentItem = b.currentItem + 1

		if b.currentItem > len(b.items)-1 || b.progress.Percent() == 1.0 {
			b.done = true
			return b, tea.Quit
		}

		cmd := b.progress.SetPercent(float64(b.currentItem) / float64(len(b.items)))

		return b, tea.Batch(b.behavior(b.items[b.currentItem]), cmd)
	default:
		return b, nil
	}
}

func (b batchModel) View() string {
	var s string
	if b.done {
		s = "✔"
	} else {
		s = b.spinner.View()
	}

	s += fmt.Sprintf(" %s\n\n", b.message)

	for _, res := range b.results {
		s += fmt.Sprintln(res)
	}

	s += "\n" + b.progress.View()

	return s
}

func NewBatchModel(message string, items []string, outputFile *os.File, behavior func(string) tea.Cmd) batchModel {
	b := batchModel{
		message:    message,
		items:      items,
		outputFile: outputFile,
		behavior:   behavior,
	}

	b.spinner = spinner.New(spinner.WithSpinner(spinner.MiniDot))
	b.progress = progress.New(progress.WithDefaultGradient())

	b.currentItem = 0
	b.results = make([]result, 0)
	return b
}
