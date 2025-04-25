package dialog

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kujtimiihoxha/opencode/internal/config"
	"github.com/kujtimiihoxha/opencode/internal/llm/models"
	"github.com/kujtimiihoxha/opencode/internal/tui/layout"
	"github.com/kujtimiihoxha/opencode/internal/tui/styles"
	"github.com/kujtimiihoxha/opencode/internal/tui/util"
)

// ModelSelectedMsg is sent when a model is selected
type ModelSelectedMsg struct {
	Model models.Model
}

// CloseModelDialogMsg is sent when the model dialog is closed
type CloseModelDialogMsg struct{}

// ModelDialog interface for the model selection dialog
type ModelDialog interface {
	tea.Model
	layout.Bindings
	SetModels(availableModels []models.Model)
	SetSelectedModel(modelID models.ModelID)
}

type modelDialogCmp struct {
	availableModels  []models.Model
	selectedIdx      int
	width            int
	height           int
	selectedModelID  models.ModelID
}

type modelKeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Enter  key.Binding
	Escape key.Binding
	J      key.Binding
	K      key.Binding
}

var modelKeys = modelKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous model"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next model"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select model"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
	J: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "next model"),
	),
	K: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "previous model"),
	),
}

func (m *modelDialogCmp) Init() tea.Cmd {
	return nil
}

func (m *modelDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, modelKeys.Up) || key.Matches(msg, modelKeys.K):
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
			return m, nil
		case key.Matches(msg, modelKeys.Down) || key.Matches(msg, modelKeys.J):
			if m.selectedIdx < len(m.availableModels)-1 {
				m.selectedIdx++
			}
			return m, nil
		case key.Matches(msg, modelKeys.Enter):
			if len(m.availableModels) > 0 {
				return m, util.CmdHandler(ModelSelectedMsg{
					Model: m.availableModels[m.selectedIdx],
				})
			}
		case key.Matches(msg, modelKeys.Escape):
			return m, util.CmdHandler(CloseModelDialogMsg{})
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m *modelDialogCmp) View() string {
	if len(m.availableModels) == 0 {
		return styles.BaseStyle.Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderBackground(styles.Background).
			BorderForeground(styles.ForgroundDim).
			Width(40).
			Render("No models available")
	}

	// Calculate max width needed for model names
	maxWidth := 40 // Minimum width
	for _, model := range m.availableModels {
		if len(model.Name) > maxWidth-4 { // Account for padding
			maxWidth = len(model.Name) + 4
		}
		if len(string(model.Provider)) > maxWidth-4 {
			maxWidth = len(string(model.Provider)) + 4
		}
	}

	// Limit height to avoid taking up too much screen space
	maxVisibleModels := min(10, len(m.availableModels))

	// Build the model list
	modelItems := make([]string, 0, maxVisibleModels)
	startIdx := 0

	// If we have more models than can be displayed, adjust the start index
	if len(m.availableModels) > maxVisibleModels {
		// Center the selected item when possible
		halfVisible := maxVisibleModels / 2
		if m.selectedIdx >= halfVisible && m.selectedIdx < len(m.availableModels)-halfVisible {
			startIdx = m.selectedIdx - halfVisible
		} else if m.selectedIdx >= len(m.availableModels)-halfVisible {
			startIdx = len(m.availableModels) - maxVisibleModels
		}
	}

	endIdx := min(startIdx+maxVisibleModels, len(m.availableModels))

	for i := startIdx; i < endIdx; i++ {
		model := m.availableModels[i]
		itemStyle := styles.BaseStyle.Width(maxWidth)
		providerStyle := styles.BaseStyle.Width(maxWidth).Foreground(styles.ForgroundDim)

		if i == m.selectedIdx {
			itemStyle = itemStyle.
				Background(styles.PrimaryColor).
				Foreground(styles.Background).
				Bold(true)
			providerStyle = providerStyle.
				Background(styles.PrimaryColor).
				Foreground(styles.Background)
		}

		name := itemStyle.Padding(0, 1).Render(model.Name)
		provider := providerStyle.Padding(0, 1).Render(string(model.Provider))
		modelItems = append(modelItems, lipgloss.JoinVertical(lipgloss.Left, name, provider))
	}

	title := styles.BaseStyle.
		Foreground(styles.PrimaryColor).
		Bold(true).
		Width(maxWidth).
		Padding(0, 1).
		Render("Select Model")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		styles.BaseStyle.Width(maxWidth).Render(""),
		styles.BaseStyle.Width(maxWidth).Render(lipgloss.JoinVertical(lipgloss.Left, modelItems...)),
		styles.BaseStyle.Width(maxWidth).Render(""),
		styles.BaseStyle.Width(maxWidth).Padding(0, 1).Foreground(styles.ForgroundDim).Render("↑/k: up  ↓/j: down  enter: select  esc: cancel"),
	)

	return styles.BaseStyle.Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(styles.Background).
		BorderForeground(styles.ForgroundDim).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

func (m *modelDialogCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(modelKeys)
}

func (m *modelDialogCmp) SetModels(availableModels []models.Model) {
	m.availableModels = availableModels

	// If we have a selected model ID, find its index
	if m.selectedModelID != "" {
		for i, model := range availableModels {
			if model.ID == m.selectedModelID {
				m.selectedIdx = i
				return
			}
		}
	}

	// Default to first model if selected not found
	m.selectedIdx = 0
}

func (m *modelDialogCmp) SetSelectedModel(modelID models.ModelID) {
	m.selectedModelID = modelID

	// Update the selected index if models are already loaded
	if len(m.availableModels) > 0 {
		for i, model := range m.availableModels {
			if model.ID == modelID {
				m.selectedIdx = i
				return
			}
		}
	}
}

// NewModelDialogCmp creates a new model selection dialog
func NewModelDialogCmp() ModelDialog {
	return &modelDialogCmp{
		availableModels: []models.Model{},
		selectedIdx:     0,
	}
}

// GetAvailableModels returns a list of models that are available based on the user's configuration
func GetAvailableModels() []models.Model {
	cfg := config.Get()
	if cfg == nil {
		return []models.Model{}
	}

	var availableModels []models.Model
	for _, model := range models.SupportedModels {
		// Check if the provider is configured and not disabled
		if provider, exists := cfg.Providers[model.Provider]; exists && !provider.Disabled {
			availableModels = append(availableModels, model)
		}
	}

	return availableModels
}