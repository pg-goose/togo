package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the taskfile path if it does not exist
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	taskfileDir := fmt.Sprintf("%s/.config/togo/", home)
	taskfilePath := fmt.Sprintf("%s/.config/togo/tasks.json", home)
	if _, err := os.Stat(taskfilePath); os.IsNotExist(err) {
		err := os.MkdirAll(taskfileDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(taskfilePath, []byte("[]"), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	p := tea.NewProgram(NewTogo(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

const taskfile = "%s/.config/togo/tasks.json"

type Task struct {
	Task     string `json:"task"`
	Complete bool   `json:"completed"`
}

func (t Task) String() string {
	complete := " "
	if t.Complete {
		complete = "x"
	}
	return fmt.Sprintf("[%s] %s", complete, t.Task)
}

type Tasks []Task
type loadTasksMsg Tasks

type Togo struct {
	taskIn textinput.Model
	tasks  Tasks
	cursor int
}

func NewTogo() Togo {
	ti := textinput.New()
	ti.Placeholder = "Task name"
	ti.Focus()
	ti.CharLimit = 128
	return Togo{
		taskIn: ti,
		tasks:  make(Tasks, 0),
		cursor: -1,
	}
}

func loadTasks() tea.Msg {
	// read file
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	taskFilePath := fmt.Sprintf(taskfile, home)
	data, err := os.ReadFile(taskFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var tasks Tasks
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatal(err)
	}
	return loadTasksMsg(tasks)
}

func (m Togo) saveTasks() tea.Msg {
	jsonStr, err := json.Marshal(m.tasks)
	if err != nil {
		log.Fatal(err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile(fmt.Sprintf(taskfile, home), jsonStr, fs.FileMode(os.O_TRUNC))
	return nil
}

func (m Togo) Init() tea.Cmd {
	return loadTasks
}

func (m Togo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.QuitMsg:
		return m, tea.Quit
	case loadTasksMsg:
		m.tasks = Tasks(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if len(m.tasks) == 0 {
				m.cursor = -1
				break
			}
			switch m.cursor {
			case -1:
				m.cursor = len(m.tasks) - 1
			case len(m.tasks) - 1:
				m.cursor -= 1
			default:
				m.cursor -= 1
			}
		case tea.KeyDown:
			if len(m.tasks) == 0 {
				m.cursor = -1
				break
			}
			switch m.cursor {
			case -1:
				m.cursor += 1
			case len(m.tasks) - 1:
				m.cursor = -1
			default:
				m.cursor += 1
			}
		case tea.KeyEnter, tea.KeySpace:
			if m.taskIn.Focused() {
				if msg.Type == tea.KeySpace {
					break
				}
				v := m.taskIn.Value()
				if v == "" {
					break
				}
				task := Task{
					Task:     v,
					Complete: false,
				}
				m.tasks = append(m.tasks, task)
				m.taskIn.Reset()
				break
			}
			c := m.tasks[m.cursor].Complete
			m.tasks[m.cursor].Complete = !c
		case tea.KeyDelete, tea.KeyBackspace:
			if m.taskIn.Focused() {
				break
			}
			index := m.cursor
			m.cursor -= 1
			m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	if m.cursor == -1 {
		m.taskIn.Focus()
	} else {
		m.taskIn.Blur()
	}
	m.saveTasks()
	m.taskIn, cmd = m.taskIn.Update(msg)
	return m, cmd
}

func (m Togo) View() string {
	ret := fmt.Sprintf("Tasks\t%d\n", len(m.tasks))
	ret += m.taskIn.View() + "\n"
	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		ret += fmt.Sprintf("%s%s\n", cursor, task.String())
	}
	return ret
}
