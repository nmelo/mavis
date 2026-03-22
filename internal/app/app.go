package app

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nmelo/typer/internal/drill"
	"github.com/nmelo/typer/internal/keyboard"
	"github.com/nmelo/typer/internal/level"
	"github.com/nmelo/typer/internal/progress"
	"github.com/nmelo/typer/internal/ui"
)

const drillsPerPhase = 10
const drillLength = 25

type state int

const (
	stateDrilling state = iota
	stateDrillComplete
	statePhaseComplete
	stateLevelComplete
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

	message string
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.saveProgress()
			return m, tea.Quit

		case tea.KeyBackspace:
			if m.state == stateDrilling {
				m.drillState.HandleBackspace()
				m.updateKeyboard()
			}

		case tea.KeyEnter:
			switch m.state {
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
					m.onDrillComplete()
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
					m.onDrillComplete()
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
	sb.WriteString("\n\n")

	sb.WriteString(ui.RenderPrompt(
		m.drillState.Prompt,
		m.drillState.Position,
		m.drillState.HasError,
	))
	sb.WriteString("\n\n")

	if m.message != "" {
		msgStyle := lipgloss.NewStyle().Bold(true)
		sb.WriteString(msgStyle.Render(m.message))
		sb.WriteString("\n\n")
	}

	sb.WriteString(m.kb.View())
	sb.WriteString("\n")

	sb.WriteString(ui.DimStyle.Render("  esc/ctrl+c to quit"))

	return sb.String()
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

func (m *Model) onDrillComplete() {
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

	if m.drillNum >= drillsPerPhase {
		phaseAcc := phaseAccuracy(m.phaseKeystrokes, m.phaseErrors)
		if level.PhasePassed(phaseAcc) {
			m.progress.CompletePhase(m.currentLevel.Number, m.currentPhase, phaseAcc, m.drillState.WPM())
			next := level.NextPhase(m.currentPhase, m.currentLevel)
			if next == "" {
				m.state = stateLevelComplete
				m.message = fmt.Sprintf("Level %d complete! Accuracy: %.1f%%. Press Enter to continue.",
					m.currentLevel.Number, phaseAcc)
			} else {
				m.state = statePhaseComplete
				m.message = fmt.Sprintf("Phase complete! Accuracy: %.1f%%. Press Enter for %s drills.",
					phaseAcc, next)
			}
		} else {
			m.state = statePhaseComplete
			m.message = fmt.Sprintf("Phase accuracy: %.1f%% (need 95%%). Press Enter to retry.",
				phaseAcc)
		}
	} else {
		m.state = stateDrillComplete
		m.message = fmt.Sprintf("Drill %d/%d done. Accuracy: %.1f%%. Press Enter to continue.",
			m.drillNum, drillsPerPhase, m.drillState.Accuracy())
	}

	m.saveProgress()
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
