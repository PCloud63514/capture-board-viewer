package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"capture-board-selector/internal/captureboard/app"
	"capture-board-selector/internal/captureboard/domain"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type phase int

const (
	phaseLoading phase = iota
	phaseSelectVideo
	phaseSelectAudio
	phaseDone
)

const refreshInterval = 5 * time.Second

type devicesLoadedMsg struct {
	devices domain.DeviceCatalog
	err     error
}

type refreshTickMsg time.Time

type model struct {
	ctx             context.Context
	useCase         app.SelectorUseCase
	spinner         spinner.Model
	phase           phase
	devices         domain.DeviceCatalog
	videoIdx        int
	audioIdx        int
	err             error
	choice          domain.Selection
	loading         bool
	lastRefreshedAt time.Time
	inputEnabled    bool
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
		ctx:          ctx,
		useCase:      useCase,
		spinner:      spin,
		phase:        phaseLoading,
		loading:      true,
		inputEnabled: term.IsTerminal(int(os.Stdin.Fd())),
	}

	result, err := tea.NewProgram(
		m,
		tea.WithInput(programInput(m.inputEnabled)),
		tea.WithOutput(os.Stdout),
	).Run()
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
	return tea.Batch(m.spinner.Tick, m.loadDevicesCmd(), refreshCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.loading {
			return m, cmd
		}
		return m, nil
	case refreshTickMsg:
		return m, tea.Batch(m.loadDevicesCmd(), refreshCmd())
	case devicesLoadedMsg:
		if msg.err != nil {
			m.loading = false
			m.err = msg.err
			m.lastRefreshedAt = time.Now()
			return m, nil
		}
		m.loading = false
		m.err = nil
		m.lastRefreshedAt = time.Now()
		m.applyDevices(msg.devices)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.err = fmt.Errorf("선택이 취소되었습니다")
			return m, tea.Quit
		case "r":
			m.err = nil
			m.loading = true
			return m, m.loadDevicesCmd()
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

	if m.loading {
		body.WriteString(fmt.Sprintf("%s 장치를 검색하는 중입니다...", m.spinner.View()))
		body.WriteString("\n")
		body.WriteString(mutedStyle.Render("ffmpeg dshow 장치 목록을 읽고 있습니다. 5초마다 자동 새로고침합니다."))
		return body.String()
	}

	if m.err != nil {
		body.WriteString(errorStyle.Render(m.err.Error()))
		body.WriteString("\n")
		body.WriteString(m.statusLine())
		body.WriteString("\n")
		body.WriteString(m.helpText())
		return body.String()
	}

	body.WriteString(sectionStyle.Render(m.currentTitle()))
	body.WriteString("\n")
	body.WriteString(m.statusLine())
	body.WriteString("\n")
	if m.choice.Video != "" {
		body.WriteString(mutedStyle.Render("선택된 비디오: " + m.choice.Video))
		body.WriteString("\n")
	}
	body.WriteString("\n")
	body.WriteString(m.renderItems())
	body.WriteString("\n\n")
	body.WriteString(m.helpText())

	return body.String()
}

func (m model) updateVideoSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.devices.Videos) == 0 {
		return m, nil
	}

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
	if len(m.devices.Audios) == 0 {
		return m, nil
	}

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

	if len(items) == 0 {
		return mutedStyle.Render("선택 가능한 장치가 없습니다.")
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

func (m model) helpText() string {
	if !m.inputEnabled {
		return mutedStyle.Render("입력 TTY가 없어 읽기 전용 모드입니다. 5초마다 자동 새로고침, q 종료")
	}
	if (m.phase == phaseSelectVideo && len(m.devices.Videos) == 0) ||
		(m.phase == phaseSelectAudio && len(m.devices.Audios) == 0) {
		return mutedStyle.Render("r 새로고침, q 종료")
	}
	return mutedStyle.Render("↑/↓ 이동, Enter 선택, r 새로고침, q 종료")
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

func (m *model) applyDevices(devices domain.DeviceCatalog) {
	m.devices = devices
	m.choice = reconcileChoice(m.choice, devices)
	m.videoIdx = reconcileIndex(devices.Videos, m.videoIdx, m.choice.Video)
	m.audioIdx = reconcileIndex(devices.Audios, m.audioIdx, m.choice.Audio)

	if m.choice.Video == "" {
		m.phase = phaseSelectVideo
		return
	}
	if m.choice.Audio == "" {
		m.phase = phaseSelectAudio
		return
	}
	m.phase = phaseDone
}

func (m model) statusLine() string {
	parts := []string{
		fmt.Sprintf("비디오 %d개", len(m.devices.Videos)),
		fmt.Sprintf("오디오 %d개", len(m.devices.Audios)),
	}

	if !m.lastRefreshedAt.IsZero() {
		parts = append(parts, "마지막 갱신 "+m.lastRefreshedAt.Format("15:04:05"))
	}

	if !m.inputEnabled {
		parts = append(parts, "읽기 전용 모드")
	}

	return mutedStyle.Render(strings.Join(parts, "  |  "))
}

func refreshCmd() tea.Cmd {
	return tea.Tick(refreshInterval, func(t time.Time) tea.Msg {
		return refreshTickMsg(t)
	})
}

func programInput(enabled bool) io.Reader {
	if enabled {
		return os.Stdin
	}
	return strings.NewReader("")
}

func reconcileChoice(choice domain.Selection, devices domain.DeviceCatalog) domain.Selection {
	if !containsDevice(devices.Videos, choice.Video) {
		choice.Video = ""
		choice.Audio = ""
		return choice
	}
	if !containsDevice(devices.Audios, choice.Audio) {
		choice.Audio = ""
	}
	return choice
}

func reconcileIndex(devices []domain.Device, current int, selected string) int {
	if len(devices) == 0 {
		return 0
	}
	if selected != "" {
		for i, device := range devices {
			if device.Name == selected {
				return i
			}
		}
	}
	if current < 0 {
		return 0
	}
	if current >= len(devices) {
		return len(devices) - 1
	}
	return current
}

func containsDevice(devices []domain.Device, name string) bool {
	if name == "" {
		return false
	}
	for _, device := range devices {
		if device.Name == name {
			return true
		}
	}
	return false
}
