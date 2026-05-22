package tui

import (
	"context"
	"fmt"
	"strings"

	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type phase int

const (
	phaseLoading phase = iota
	phaseSelectVideo
	phaseSelectAudio
	phaseDone
)

type devicesLoadedMsg struct {
	devices domain.DeviceCatalog
	err     error
}

type model struct {
	ctx      context.Context
	useCase  app.SelectorUseCase
	spinner  spinner.Model
	phase    phase
	devices  domain.DeviceCatalog
	videoIdx int
	audioIdx int
	err      error
	choice   domain.Selection
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("62")).
			Padding(0, 1)
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86"))
	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))
	mutedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("203"))
)

func Run(ctx context.Context, useCase app.SelectorUseCase) (domain.Selection, error) {
	spin := spinner.New()
	spin.Spinner = spinner.Line
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))

	m := model{
		ctx:     ctx,
		useCase: useCase,
		spinner: spin,
		phase:   phaseLoading,
	}

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return domain.Selection{}, err
	}

	finalModel := result.(model)
	if finalModel.err != nil {
		return domain.Selection{}, finalModel.err
	}

	return finalModel.choice, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadDevicesCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.phase == phaseLoading {
			return m, cmd
		}
		return m, nil
	case devicesLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.devices = msg.devices
		m.phase = phaseSelectVideo
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.err = fmt.Errorf("선택이 취소되었습니다")
			return m, tea.Quit
		}

		switch m.phase {
		case phaseSelectVideo:
			return m.updateVideoSelection(msg)
		case phaseSelectAudio:
			return m.updateAudioSelection(msg)
		}
	}

	return m, nil
}

func (m model) View() string {
	var body strings.Builder

	body.WriteString(titleStyle.Render(" Capture Board Selector "))
	body.WriteString("\n\n")

	if m.phase == phaseLoading {
		body.WriteString(fmt.Sprintf("%s 장치를 검색하는 중입니다...", m.spinner.View()))
		body.WriteString("\n")
		body.WriteString(mutedStyle.Render("ffmpeg dshow 장치 목록을 읽고 있습니다."))
		return body.String()
	}

	if m.err != nil {
		body.WriteString(errorStyle.Render(m.err.Error()))
		body.WriteString("\n")
		body.WriteString(mutedStyle.Render("q 또는 Ctrl+C로 종료하세요."))
		return body.String()
	}

	body.WriteString(sectionStyle.Render(m.currentTitle()))
	body.WriteString("\n")
	if m.choice.Video != "" {
		body.WriteString(mutedStyle.Render("선택된 비디오: " + m.choice.Video))
		body.WriteString("\n")
	}
	body.WriteString("\n")
	body.WriteString(m.renderItems())
	body.WriteString("\n\n")
	body.WriteString(mutedStyle.Render("↑/↓ 이동, Enter 선택, q 종료"))

	return body.String()
}

func (m model) updateVideoSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.videoIdx > 0 {
			m.videoIdx--
		}
	case "down":
		if m.videoIdx < len(m.devices.Videos)-1 {
			m.videoIdx++
		}
	case "enter":
		m.choice.Video = m.devices.Videos[m.videoIdx].Name
		m.phase = phaseSelectAudio
	}
	return m, nil
}

func (m model) updateAudioSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up":
		if m.audioIdx > 0 {
			m.audioIdx--
		}
	case "down":
		if m.audioIdx < len(m.devices.Audios)-1 {
			m.audioIdx++
		}
	case "enter":
		m.choice.Audio = m.devices.Audios[m.audioIdx].Name
		m.phase = phaseDone
		return m, tea.Quit
	}
	return m, nil
}

func (m model) renderItems() string {
	items := m.devices.Videos
	cursor := m.videoIdx
	if m.phase == phaseSelectAudio {
		items = m.devices.Audios
		cursor = m.audioIdx
	}

	lines := make([]string, 0, len(items))
	for i, item := range items {
		prefix := "  "
		style := lipgloss.NewStyle()
		if i == cursor {
			prefix = "› "
			style = selectedStyle
		}
		lines = append(lines, style.Render(prefix+item.Name))
	}

	return strings.Join(lines, "\n")
}

func (m model) currentTitle() string {
	if m.phase == phaseSelectAudio {
		return "오디오 장치 선택"
	}
	return "비디오 장치 선택"
}

func (m model) loadDevicesCmd() tea.Cmd {
	return func() tea.Msg {
		devices, err := m.useCase.LoadDevices(m.ctx)
		return devicesLoadedMsg{
			devices: devices,
			err:     err,
		}
	}
}
