package app

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/mavis/internal/drill"
	"github.com/nmelo/mavis/internal/keyboard"
	"github.com/nmelo/mavis/internal/level"
	"github.com/nmelo/mavis/internal/progress"
	"github.com/nmelo/mavis/internal/ui"
)

const drillsPerPhase = 10
const drillLength = 25

type state int

const (
	stateDrilling state = iota
	stateDrillComplete
	statePhaseComplete
	stateLevelComplete
	stateCelebrating
)

type Model struct {
	progress     *progress.Progress
	progressPath string
	words        []string

	currentLevel level.Level
	currentPhase string
	drillNum     int
	drillState   *drill.DrillState
	kb           *keyboard.Keyboard
	state        state

	phaseKeystrokes int
	phaseErrors     int

	message     string
	celebration *celebration

	width  int
	height int
}

func New(prog *progress.Progress, progressPath string, words []string) Model {
	lvl := level.Get(prog.CurrentLevel)
	unlocked := unlockedKeySet(prog.CurrentLevel)

	m := Model{
		progress:     prog,
		progressPath: progressPath,
		words:        words,
		currentLevel: lvl,
		currentPhase: prog.CurrentPhase,
		drillNum:     1,
		state:        stateDrilling,
	}

	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.firstPromptKey())

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.state == stateCelebrating && m.celebration != nil {
			m.celebration.tick()
			if !m.celebration.active() {
				return m, nil
			}
			return m, tickCmd()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.saveProgress()
			return m, tea.Quit

		case tea.KeyCtrlN:
			m.skipToLevel(m.currentLevel.Number + 1)

		case tea.KeyCtrlB:
			m.skipToLevel(m.currentLevel.Number - 1)

		case tea.KeyCtrlT:
			// Cheat: trigger pass celebration
			m.celebration = spawnCelebration(tierPass, m.state, "Phase complete! Accuracy: 97.3%.")
			m.state = stateCelebrating
			return m, nil

		case tea.KeyCtrlY:
			// Cheat: trigger perfect celebration
			m.celebration = spawnCelebration(tierPerfect, m.state, "PERFECT!")
			m.state = stateCelebrating
			return m, tickCmd()

		case tea.KeyCtrlU:
			// Cheat: trigger level-up celebration
			msg := fmt.Sprintf("Level %d complete!", m.currentLevel.Number)
			m.celebration = spawnCelebration(tierLevelUp, m.state, msg)
			m.state = stateCelebrating
			return m, tickCmd()

		case tea.KeyBackspace:
			if m.state == stateDrilling {
				m.drillState.HandleBackspace()
				m.updateKeyboard()
			}

		case tea.KeyEnter:
			switch m.state {
			case stateCelebrating:
				if m.celebration != nil {
					m.state = m.celebration.nextState
					m.celebration = nil
				}
			case stateDrillComplete:
				m.advanceDrill()
			case statePhaseComplete:
				m.advancePhase()
			case stateLevelComplete:
				m.advanceLevel()
			}

		case tea.KeyRunes:
			if m.state == stateDrilling && len(msg.Runes) > 0 {
				m.drillState.HandleKey(msg.Runes[0])
				m.updateKeyboard()

				if m.drillState.Complete {
					cmd := m.onDrillComplete()
					return m, cmd
				}
			}

		case tea.KeySpace:
			if m.state == stateDrilling {
				// Ignore space during char drills (space is for word/code drills only)
				if m.currentPhase == progress.PhaseChars {
					break
				}
				m.drillState.HandleKey(' ')
				m.updateKeyboard()

				if m.drillState.Complete {
					cmd := m.onDrillComplete()
					return m, cmd
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder

	wpm := m.drillState.WPM()
	acc := m.drillState.Accuracy()
	sb.WriteString(ui.RenderTopBar(m.currentLevel.Name, wpm, acc, m.drillNum, drillsPerPhase))
	sb.WriteString("\n\n\n")

	sb.WriteString(ui.RenderPrompt(
		m.drillState.Prompt,
		m.drillState.Position,
		m.drillState.HasError,
	))
	sb.WriteString("\n\n")

	if m.state == stateCelebrating && m.celebration != nil {
		sb.WriteString(m.celebration.render(m.width))
		sb.WriteString("\n\n")
	} else if m.message != "" {
		msgStyle := lipgloss.NewStyle().Bold(true)
		sb.WriteString(msgStyle.Render(m.message))
		sb.WriteString("\n\n")
	}

	sb.WriteString("\n")
	sb.WriteString(m.kb.View())
	sb.WriteString("\n\n")

	sb.WriteString(ui.DimStyle.Render("esc to quit    ctrl+n next level    ctrl+b prev level"))

	content := sb.String()

	// Center the content block in the terminal without
	// centering each line individually (which breaks keyboard stagger)
	if m.width > 0 && m.height > 0 {
		block := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Render(content)
		return block
	}

	return content
}

func (m *Model) generateDrill(unlocked map[rune]bool) *drill.DrillState {
	var prompt []rune

	switch m.currentPhase {
	case progress.PhaseChars:
		prompt = drill.GenerateCharDrill(
			level.UnlockedKeys(m.currentLevel.Number),
			m.currentLevel.NewKeys,
			drillLength,
		)
	case progress.PhaseWords:
		filtered := drill.FilterWords(m.words, unlocked)
		words := drill.GenerateWordDrill(filtered, 5)
		prompt = []rune(strings.Join(words, " "))
	case progress.PhaseCode:
		allSnippets := drill.LoadCodeSnippets()
		filtered := drill.FilterCodeSnippets(allSnippets, unlocked)
		if len(filtered) > 0 {
			snippet := filtered[rand.Intn(len(filtered))]
			prompt = []rune(snippet.Code)
		} else {
			prompt = drill.GenerateCharDrill(
				level.UnlockedKeys(m.currentLevel.Number),
				m.currentLevel.NewKeys,
				drillLength,
			)
		}
	}

	return drill.NewDrillState(prompt)
}

func (m *Model) onDrillComplete() tea.Cmd {
	m.phaseKeystrokes += m.drillState.TotalKeystrokes
	m.phaseErrors += m.drillState.Errors

	session := progress.Session{
		Date:     time.Now(),
		Duration: int(time.Since(m.drillState.StartTime).Seconds()),
		Level:    m.currentLevel.Number,
		Phase:    m.currentPhase,
		Accuracy: m.drillState.Accuracy(),
		WPM:      m.drillState.WPM(),
		KeyStats: convertKeyStats(m.drillState.KeyStats()),
	}
	m.progress.RecordSession(session)

	drillAcc := m.drillState.Accuracy()

	if m.drillNum >= drillsPerPhase {
		phaseAcc := phaseAccuracy(m.phaseKeystrokes, m.phaseErrors)
		if level.PhasePassed(phaseAcc) {
			m.progress.CompletePhase(m.currentLevel.Number, m.currentPhase, phaseAcc, m.drillState.WPM())
			next := level.NextPhase(m.currentPhase, m.currentLevel)
			if next == "" {
				// Level complete - highest tier celebration
				msg := fmt.Sprintf("Level %d complete! Accuracy: %.1f%%.", m.currentLevel.Number, phaseAcc)
				m.celebration = spawnCelebration(tierLevelUp, stateLevelComplete, msg)
				m.state = stateCelebrating
				m.message = ""
				m.saveProgress()
				return tickCmd()
			}
			// Phase passed - green styled message (no animation)
			msg := fmt.Sprintf("Phase complete! Accuracy: %.1f%%. Press Enter for %s drills.", phaseAcc, next)
			m.state = statePhaseComplete
			m.message = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ui.ColorCorrect)).Render("+ " + msg)
		} else {
			m.state = statePhaseComplete
			m.message = fmt.Sprintf("Phase accuracy: %.1f%% (need 95%%). Press Enter to retry.", phaseAcc)
		}
	} else {
		// Individual drill complete
		if drillAcc == 100.0 {
			msg := fmt.Sprintf("PERFECT! Drill %d/%d.", m.drillNum, drillsPerPhase)
			m.celebration = spawnCelebration(tierPerfect, stateDrillComplete, msg)
			m.state = stateCelebrating
			m.message = ""
			m.saveProgress()
			return tickCmd()
		}
		m.state = stateDrillComplete
		m.message = fmt.Sprintf("Drill %d/%d done. Accuracy: %.1f%%. Press Enter to continue.",
			m.drillNum, drillsPerPhase, drillAcc)
	}

	m.saveProgress()
	return nil
}

func (m *Model) advanceDrill() {
	m.drillNum++
	unlocked := unlockedKeySet(m.currentLevel.Number)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.firstPromptKey())
	m.state = stateDrilling
	m.message = ""
}

func (m *Model) advancePhase() {
	phaseAcc := phaseAccuracy(m.phaseKeystrokes, m.phaseErrors)
	if level.PhasePassed(phaseAcc) {
		next := level.NextPhase(m.currentPhase, m.currentLevel)
		m.currentPhase = next
	}
	// If phase failed, m.currentPhase stays the same (retry)

	m.phaseKeystrokes = 0
	m.phaseErrors = 0
	m.drillNum = 1
	m.progress.CurrentPhase = m.currentPhase

	unlocked := unlockedKeySet(m.currentLevel.Number)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.firstPromptKey())
	m.state = stateDrilling
	m.message = ""
}

func (m *Model) skipToLevel(n int) {
	total := len(level.All())
	if n < 1 {
		n = 1
	}
	if n > total {
		n = total
	}

	m.progress.CurrentLevel = n
	m.progress.CurrentPhase = progress.PhaseChars
	m.currentLevel = level.Get(n)
	m.currentPhase = progress.PhaseChars
	m.phaseKeystrokes = 0
	m.phaseErrors = 0
	m.drillNum = 1

	unlocked := unlockedKeySet(n)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.firstPromptKey())
	m.state = stateDrilling
	m.message = fmt.Sprintf("Jumped to Level %d: %s", n, m.currentLevel.Name)
	m.saveProgress()
}

func (m *Model) advanceLevel() {
	if m.currentLevel.Number < len(level.All()) {
		m.progress.CurrentLevel = m.currentLevel.Number + 1
		m.progress.CurrentPhase = progress.PhaseChars
		m.currentLevel = level.Get(m.progress.CurrentLevel)
		m.currentPhase = progress.PhaseChars
	}

	m.phaseKeystrokes = 0
	m.phaseErrors = 0
	m.drillNum = 1

	unlocked := unlockedKeySet(m.currentLevel.Number)
	m.drillState = m.generateDrill(unlocked)
	m.kb = keyboard.New(unlocked, m.firstPromptKey())
	m.state = stateDrilling
	m.message = ""
}

func (m *Model) firstPromptKey() rune {
	if len(m.drillState.Prompt) > 0 {
		return m.drillState.Prompt[0]
	}
	return 'f' // fallback, should never happen
}

func (m *Model) updateKeyboard() {
	if !m.drillState.Complete && m.drillState.Position < len(m.drillState.Prompt) {
		m.kb.Update(m.drillState.Prompt[m.drillState.Position])
	}
}

func (m *Model) saveProgress() {
	_ = progress.Save(m.progress, m.progressPath)
}

func unlockedKeySet(levelNum int) map[rune]bool {
	keys := level.UnlockedKeys(levelNum)
	set := make(map[rune]bool, len(keys))
	for _, k := range keys {
		set[k] = true
	}
	if levelNum >= 5 {
		set[' '] = true
	}
	return set
}

func phaseAccuracy(keystrokes, errors int) float64 {
	if keystrokes == 0 {
		return 100.0
	}
	correct := keystrokes - errors
	return float64(correct) / float64(keystrokes) * 100
}

func convertKeyStats(drillStats map[rune]drill.KeyStat) map[string]progress.KeyStat {
	result := make(map[string]progress.KeyStat)
	for k, v := range drillStats {
		result[string(k)] = progress.KeyStat{
			Correct:   v.Correct,
			Incorrect: v.Incorrect,
		}
	}
	return result
}
