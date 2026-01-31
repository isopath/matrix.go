package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:embed assets/*.txt
var assetsFS embed.FS

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title    string
	filename string
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type column struct {
	x      int
	height int
	offset int
}

type tickMsg time.Time

type model struct {
	list     list.Model
	choice   string
	content  string
	quitting bool
	viewing  bool
	width    int
	height   int
	columns  []column
}

func (m model) Init() tea.Cmd {
	if m.viewing {
		return tick()
	}
	return nil
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.list.Items() != nil {
			m.list.SetWidth(msg.Width)
		}
		if m.viewing {
			m.initColumns()
		}
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			if m.viewing {
				// If there's no list (direct mode), quit entirely
				if m.list.Items() == nil {
					m.quitting = true
					return m, tea.Quit
				}
				// Otherwise return to list
				m.viewing = false
				m.choice = ""
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if !m.viewing && m.list.Items() != nil {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = i.filename
					content, err := readAsset(m.choice)
					if err != nil {
						m.content = fmt.Sprintf("Error: %v", err)
					} else {
						m.content = content
					}
					m.viewing = true
					m.initColumns()
					return m, tick()
				}
			}
		}

	case tickMsg:
		if m.viewing {
			m.updateColumns()
			return m, tick()
		}
	}

	var cmd tea.Cmd
	if m.list.Items() != nil {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m *model) initColumns() {
	if m.width == 0 || len(m.content) == 0 {
		return
	}

	numCols := m.width / 2
	if numCols < 1 {
		numCols = 1
	}

	m.columns = make([]column, numCols)
	runes := []rune(m.content)

	for i := range m.columns {
		m.columns[i] = column{
			x:      i * 2,
			height: rand.Intn(m.height) + 5,
			offset: rand.Intn(len(runes)),
		}
	}
}

func (m *model) updateColumns() {
	runes := []rune(m.content)
	if len(runes) == 0 {
		return
	}

	for i := range m.columns {
		m.columns[i].offset++
		if m.columns[i].offset >= len(runes) {
			m.columns[i].offset = 0
		}

		if rand.Float32() < 0.02 {
			m.columns[i].height = rand.Intn(m.height) + 5
		}
	}
}

func (m model) View() string {
	if m.viewing {
		return m.matrixView()
	}
	if m.quitting {
		return quitTextStyle.Render("Quitting is for losers.")
	}
	if m.list.Items() != nil {
		return "\n" + m.list.View()
	}
	return ""
}

func (m model) matrixView() string {
	if len(m.columns) == 0 || len(m.content) == 0 {
		return ""
	}

	grid := make([][]rune, m.height)
	for i := range grid {
		grid[i] = make([]rune, m.width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	runes := []rune(m.content)

	for _, col := range m.columns {
		if col.x >= m.width {
			continue
		}

		for row := 0; row < col.height && row < m.height; row++ {
			charIdx := (col.offset + row) % len(runes)
			if charIdx < 0 {
				charIdx += len(runes)
			}

			if row < m.height && col.x < m.width {
				grid[row][col.x] = runes[charIdx]
			}
		}
	}

	var result strings.Builder
	for rowIdx, row := range grid {
		for _, char := range row {
			if char != ' ' {
				color := getRandomColor()
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				result.WriteString(style.Render(string(char)))
			} else {
				result.WriteRune(' ')
			}
		}
		if rowIdx < len(grid)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func getRandomColor() string {
	colors := []string{
		"196", "197", "198", "199", "200", "201", // pinks/magentas
		"160", "161", "162", "163", "164", "165", // reds
		"202", "203", "204", "205", "206", "207", // oranges
		"208", "209", "210", "211", "212", "213", // light oranges
		"214", "215", "216", "217", "218", "219", // yellows
		"220", "221", "222", "223", "224", "225", // light yellows
		"226", "227", "228", "229", "230", "231", // whites
		"82", "83", "84", "85", "86", "87", // greens
		"28", "29", "30", "31", "32", "33", // dark greens
		"40", "41", "42", "43", "44", "45", // cyans
		"46", "47", "48", "49", "50", "51", // teals
		"75", "76", "77", "78", "79", "80", // blues
		"63", "64", "65", "66", "67", "68", // dark blues
		"90", "91", "92", "93", "94", "95", // purples
		"129", "130", "131", "132", "133", "134", // violets
		"141", "142", "143", "144", "145", "146", // light purples
	}
	return colors[rand.Intn(len(colors))]
}

func readAsset(filename string) (string, error) {
	data, err := assetsFS.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getAvailableFiles() []string {
	files := []string{}
	entries, err := assetsFS.ReadDir("assets")
	if err != nil {
		return files
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
			files = append(files, "assets/"+entry.Name())
		}
	}
	return files
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var showOptions bool
	var filePath string
	flag.BoolVar(&showOptions, "options", false, "Show list of available files to choose from")
	flag.StringVar(&filePath, "file", "", "Path to a specific .txt file to display")
	flag.Parse()

	// Default: show matrix directly with greatwork.txt
	targetFile := "assets/greatwork.txt"

	if filePath != "" {
		// User specified a file
		targetFile = filePath
		// If it's not a path with assets/, assume it's in assets/
		if !strings.Contains(targetFile, "/") && !strings.HasPrefix(targetFile, "assets/") {
			targetFile = "assets/" + targetFile
		}
	}

	if showOptions {
		// Show list widget
		files := getAvailableFiles()
		if len(files) == 0 {
			fmt.Println("No .txt files found in assets")
			os.Exit(1)
		}

		var items []list.Item
		for _, f := range files {
			title := f
			// Convert filename to friendly title
			title = strings.TrimPrefix(title, "assets/")
			title = strings.TrimSuffix(title, ".txt")
			title = strings.ReplaceAll(title, "-", " ")
			// Capitalize words
			words := strings.Split(title, " ")
			for i, word := range words {
				if len(word) > 0 {
					words[i] = strings.ToUpper(word[:1]) + word[1:]
				}
			}
			title = strings.Join(words, " ")
			items = append(items, item{title: title, filename: f})
		}

		const defaultWidth = 20
		l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
		l.Title = "Choose text:"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		m := model{list: l}

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	} else {
		// Direct matrix view
		content, err := readAsset(targetFile)
		if err != nil {
			// Try reading from filesystem as fallback
			data, err2 := os.ReadFile(targetFile)
			if err2 != nil {
				fmt.Printf("Error reading file %s: %v\n", targetFile, err)
				os.Exit(1)
			}
			content = string(data)
		}

		m := model{
			content: content,
			viewing: true,
		}

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}
}
