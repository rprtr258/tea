package autocomplete

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rprtr258/fun"

	"github.com/rprtr258/tea"
	"github.com/rprtr258/tea/components/textinput"
	"github.com/rprtr258/tea/styles"
)

type (
	gotReposSuccessMsg []repo
	gotReposErrMsg     error
)

type repo struct {
	Name string `json:"name"`
}

const reposURL = "https://api.github.com/orgs/charmbracelet/repos"

func getRepos() tea.Msg {
	req, err := http.NewRequest(http.MethodGet, reposURL, http.NoBody) //nolint:noctx
	if err != nil {
		return gotReposErrMsg(err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gotReposErrMsg(err)
	}
	defer resp.Body.Close() // nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return gotReposErrMsg(err)
	}

	var repos []repo
	if err := json.Unmarshal(data, &repos); err != nil {
		return gotReposErrMsg(err)
	}

	return gotReposSuccessMsg(repos)
}

type model struct {
	textInput textinput.Model
}

func (m model) Init(yield func(...tea.Cmd)) {
	yield(getRepos, textinput.Blink)
}

func (m model) Update(msg tea.Msg, yield func(...tea.Cmd)) {
	switch msg := msg.(type) {
	case tea.MsgKey:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			yield(tea.Quit)
			return
		}
	case gotReposSuccessMsg:
		m.textInput.SetSuggestions(fun.Map[string](func(r repo) string {
			return r.Name
		}, msg...))
	}

	m.textInput.Update(msg, yield)
}

func (m model) View(vb tea.Viewbox) {
	vb.WriteLine(`What’s your favorite Charm repository?`)
	m.textInput.View(vb.PaddingTop(2))
	vb.WriteText(4, 0, "(tab to complete, ctrl+n/ctrl+p to cycle through suggestions, esc to quit)")
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = "charmbracelet/"
	ti.Placeholder = "repo..."
	ti.Cursor.Style = styles.Style{}.Foreground(styles.FgColor("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.ShowSuggestions = true
	return model{textInput: ti}
}

func Main(ctx context.Context) error {
	_, err := tea.NewProgram(ctx, initialModel()).Run()
	return err
}
