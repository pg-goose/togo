package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Task represents an individual task.
type Task struct {
	Task     string `json:"task"`
	Complete bool   `json:"completed"`
}

func (t Task) String() string {
	status := " "
	if t.Complete {
		status = "x"
	}
	return fmt.Sprintf("[%s] %s", status, t.Task)
}

// Tasks is a slice of Task.
type Tasks []Task

// loadTasksMsg is a message containing loaded tasks.
type loadTasksMsg Tasks

// Togo is the Bubble Tea model.
type Togo struct {
	taskFilePath string
	taskIn       textinput.Model
	tasks        Tasks
	cursor       int
}

// NewTogo initializes the Togo model.
func NewTogo(path string) *Togo {
	ti := textinput.New()
	ti.Placeholder = "Task name"
	ti.Focus()
	ti.CharLimit = 128

	return &Togo{
		taskFilePath: path,
		taskIn:       ti,
		tasks:        make(Tasks, 0),
		cursor:       -1, // no task selected
	}
}

// loadTasks reads tasks from the file and returns a message.
func (m *Togo) loadTasks() tea.Msg {
	data, err := os.ReadFile(m.taskFilePath)
	if err != nil {
		log.Println("error reading tasks file:", err)
		return nil
	}
	var tasks Tasks
	if err := json.Unmarshal(data, &tasks); err != nil {
		log.Println("error unmarshalling tasks:", err)
		return nil
	}
	return loadTasksMsg(tasks)
}

// saveTasks writes the current tasks to the file.
func (m *Togo) saveTasks() {
	jsonData, err := json.Marshal(m.tasks)
	if err != nil {
		log.Println("error marshalling tasks:", err)
		return
	}
	if err := os.WriteFile(m.taskFilePath, jsonData, 0644); err != nil {
		log.Println("error writing tasks file:", err)
	}
}

// Init returns an initial command to load tasks.
func (m *Togo) Init() tea.Cmd {
	// Wrap loadTasks in a command function.
	return func() tea.Msg {
		return m.loadTasks()
	}
}

// Update handles incoming messages and key events.
func (m *Togo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if len(m.tasks) > 0 {
				if m.cursor <= -1 {
					m.cursor = len(m.tasks) - 1
				} else {
					m.cursor--
				}
			}
		case tea.KeyDown:
			if len(m.tasks) > 0 {
				if m.cursor >= len(m.tasks)-1 {
					m.cursor = -1
				} else {
					m.cursor++
				}
			}
		case tea.KeyEnter:
			if m.taskIn.Focused() {
				v := m.taskIn.Value()
				if v != "" {
					m.tasks = append(m.tasks, Task{Task: v})
					m.taskIn.Reset()
				}
			} else if m.cursor >= 0 && m.cursor < len(m.tasks) {
				m.tasks[m.cursor].Complete = !m.tasks[m.cursor].Complete
			}
		case tea.KeySpace:
			if !m.taskIn.Focused() && m.cursor >= 0 && m.cursor < len(m.tasks) {
				m.tasks[m.cursor].Complete = !m.tasks[m.cursor].Complete
			}
		case tea.KeyDelete, tea.KeyBackspace:
			if !m.taskIn.Focused() && len(m.tasks) > 0 && m.cursor >= 0 && m.cursor < len(m.tasks) {
				m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
				if m.cursor >= len(m.tasks) {
					m.cursor = len(m.tasks) - 1
				}
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case loadTasksMsg:
		m.tasks = Tasks(msg)
	}

	// If no task is selected, keep the text input focused.
	if m.cursor == -1 {
		m.taskIn.Focus()
	} else {
		m.taskIn.Blur()
	}

	m.taskIn, cmd = m.taskIn.Update(msg)
	m.saveTasks() // Persist tasks on every update.
	return m, cmd
}

func Count[T any](s []T, p func(T) bool) (c int) {
	for _, y := range s {
		if p(y) {
			c++
		}
	}
	return c
}

// View renders the UI.
func (m *Togo) View() string {
	s := fmt.Sprintf("Tasks (%d/%d):\n", Count(m.tasks, func(t Task) bool { return t.Complete }), len(m.tasks))
	s += m.taskIn.View() + "\n"
	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, task.String())
	}
	return s
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	taskfileDir := fmt.Sprintf("%s/.config/togo", home)
	taskfilePath := fmt.Sprintf("%s/tasks.json", taskfileDir)

	// Create the configuration directory if it doesn't exist.
	if _, err := os.Stat(taskfileDir); os.IsNotExist(err) {
		if err := os.MkdirAll(taskfileDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	// Create the tasks file if it doesn't exist.
	if _, err := os.Stat(taskfilePath); os.IsNotExist(err) {
		if err := os.WriteFile(taskfilePath, []byte("[]"), 0644); err != nil {
			log.Fatal(err)
		}
	}

	model := NewTogo(taskfilePath)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
