package main

import (
	"flag"
	"fmt"
	"io"

	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type config struct {
	repo    string
	org     string
	user    string
	topic   string
	page    int
	verbose bool
	version bool
}

func parseFlags() config {
	org := flag.String("org", "", "an optional GitHub organization (i.e. 'python') to scan the repositories from (100 max) ; use repository for current folder if omitted and no '-repo' nor '-user' flag")
	page := flag.Int("page", 1, "page number for '-repo' and '-user' flags, 100 repositories per page")
	repo := flag.String("repo", "", "an optional GitHub repository (i.e. 'python/peps') ; use repository for current folder if omitted and no '-org' nor '-user' flag")
	topic := flag.String("topic", "", "an optional GitHub topic (i.e. 'testing') to filter the repositories ; ignored if no '-user' nor '-org' flag")
	user := flag.String("user", "", "an optional GitHub user (i.e. 'torvalds') to scan the repositories from (100 max) ; use repository for current folder if omitted and no '-repo' nor '-org' flag")
	verbose := flag.Bool("verbose", false, "verbose mode outputs several lines per repository ; non-verbose mode outputs a one-liner per repository ; default: false")
	version := flag.Bool("version", false, "outputs version-related information")
	flag.Parse()
	return config{*repo, *org, *user, *topic, *page, *verbose, *version}
}

type model struct {
	list    list.Model
	spinner spinner.Model
	repos   []repo
}

func initialModel() model {
	s := spinner.New()
	items := []list.Item{
		item("Ramen"),
		item("Tomato Soup"),
		item("Hamburgers"),
	}
	const defaultWidth = 20
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return model{
		list:    l,
		spinner: s,
		repos:   []repo{},
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		default:
			return m, nil
		}

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	// The header
	s := fmt.Sprintf("\n\n   %s Loading from GitHub... press q to quit.\n\n", m.spinner.View())

	// Iterate over our choices
	for _, repo := range m.repos {

		// Render the row
		s += fmt.Sprintf("Repo %s \n", repo.Name)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

type owner struct{ Login string }

type repo struct {
	Name        string
	Full_name   string
	Owner       owner
	Description string
	Topics      []string
	Visibility  string
	Fork        bool
}

type collaborator struct{}

type version struct {
	commit string
	date   time.Time
	dirty  bool
}

func main() {
	// config := parseFlags()
	const defaultWidth = 20

	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Print(err)
	}

	/*
	if config.version {
		version := getVersion()
		dirty := ""
		if version.dirty {
			dirty = "(dirty)"
		}
		fmt.Printf("Commit %s (%s) %s\n", version.commit, version.date, dirty)
	} else if len(config.org) > 0 || len(config.user) > 0 {
		repos, error := getRepos(config)
		if error != nil {
			fmt.Print(error)
			os.Exit(2)
		}
		for _, repo := range repos {
			repoMessage, repo, validRepo := scanRepo(config, repo.Full_name)
			if validRepo {
				fmt.Print(repoMessage)
				collaboratorsMessage := scanCollaborators(config, repo.Full_name)
				fmt.Print(collaboratorsMessage)
				if strings.Compare(repo.Visibility, "public") == 0 {
					communityScoreMessage := scanCommunityScore(config, repo.Full_name)
					fmt.Print(communityScoreMessage)
				}
			}
			fmt.Println()
		}
	} else {
		repoWithOrg, error := getRepo(config)
		if error != nil {
			fmt.Print(error)
			if strings.Contains(error.Error(), "none of the git remotes configured for this repository point to a known GitHub host") {
				print("If current folder is related to a GitHub repository, please check 'gh auth status' and 'gh config list'.")
			}
			os.Exit(1)
		}
		repoMessage, repo, validRepo := scanRepo(config, repoWithOrg)
		if validRepo {
			fmt.Print(repoMessage)
			collaboratorsMessage := scanCollaborators(config, repoWithOrg)
			fmt.Print(collaboratorsMessage)
			if !repo.Fork && strings.Compare(repo.Visibility, "public") == 0 {
				communityScoreMessage := scanCommunityScore(config, repoWithOrg)
				fmt.Print(communityScoreMessage)
			}
			fmt.Println()
		}
	}
	 */
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getRepos(config config) ([]repo, error) {
	if len(config.org) == 0 && len(config.user) == 0 {
		return []repo{}, nil
	}
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return []repo{}, err
	}
	if len(config.org) > 0 {
		// https://docs.github.com/en/rest/reference/repos#list-organization-repositories
		repos := []repo{}
		err = client.Get(
			"orgs/"+config.org+"/repos?sort=full_name&per_page=100&page="+strconv.Itoa(config.page),
			&repos)
		return reposWithTopic(repos, config.topic), err
	} else {
		// https://docs.github.com/en/rest/reference/repos#list-repositories-for-a-user
		repos := []repo{}
		err = client.Get(
			"users/"+config.user+"/repos?sort=full_name&per_page=100&page="+strconv.Itoa(config.page),
			&repos)
		return reposWithTopic(repos, config.topic), err
	}
}

func reposWithTopic(repos []repo, topic string) []repo {
	if len(topic) > 0 {
		filtered := []repo{}
		for _, repo := range repos {
			if contains(repo.Topics, topic) {
				filtered = append(filtered, repo)
			}
		}
		return filtered
	}
	return repos
}

func getRepo(config config) (string, error) {
	if len(config.repo) > 1 {
		return config.repo, nil
	}
	if config.verbose {
		fmt.Printf("(current repo)\n")
	}
	currentRepo, error := gh.CurrentRepository()
	if error != nil {
		return "", error
	}
	return currentRepo.Owner() + "/" + currentRepo.Name(), nil
}

func scanRepo(config config, repoWithOrg string) (message string, repository repo, validRepo bool) {
	// https://docs.github.com/en/rest/reference/repos#get-a-repository-readme
	readme := struct {
		Name string
	}{}
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	err = client.Get(
		"repos/"+repoWithOrg+"/readme",
		&readme)
	if config.verbose {
		message += repoWithOrg + " has: "
	}
	if !config.verbose && (len(config.repo) > 1 || len(config.user) > 1 || len(config.org) > 1) {
		message += repoWithOrg + ": "
	}
	if len(readme.Name) > 0 {
		if config.verbose {
			message += "\n  - a README ☑️"
		} else {
			message += "README ☑️, "
		}
	} else if strings.HasPrefix(err.Error(), "HTTP 404: Not Found") {
		if config.verbose {
			message += "\n  - no README 😇"
		} else {
			message += "no README 😇, "
		}
	} else {
		fmt.Print(err)
	}

	repo := struct {
		Name        string
		Full_name   string
		Owner       owner
		Description string
		Topics      []string
		Visibility  string
		Fork        bool
	}{}
	errRepo := client.Get(
		"repos/"+repoWithOrg,
		&repo)
	if errRepo != nil {
		fmt.Print(errRepo)
		return
	}
	if len(repo.Description) > 0 {
		if config.verbose {
			message += "\n  - a description ☑️"
		} else {
			message += "description ☑️, "
		}
	} else {
		if config.verbose {
			message += "\n  - no description 😇"
		} else {
			message += "no description 😇, "
		}
	}
	if len(repo.Topics) > 0 {
		if config.verbose {
			message += "\n  - topics ☑️"
		} else {
			message += "topics ☑️, "
		}
	} else {
		if config.verbose {
			message += "\n  - no topics 😇"
		} else {
			message += "no topics 😇, "
		}
	}
	return message, repo, true
}

func scanCollaborators(config config, repoWithOrg string) string {
	// https://docs.github.com/en/rest/reference/collaborators#list-repository-collaborators
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	collaborators := []collaborator{}
	err = client.Get(
		"repos/"+repoWithOrg+"/collaborators",
		&collaborators)
	message := ""
	if err != nil && len(err.Error()) > 0 {
		if strings.HasPrefix(err.Error(), "HTTP 403") {
			// 🤫
		} else {
			fmt.Print(err)
		}
	} else if len(collaborators) <= 1 {
		if config.verbose {
			message += fmt.Sprintf("\n  - %d collaborator 👤", len(collaborators))
		} else {
			message += fmt.Sprintf("%d collaborator 👤, ", len(collaborators))
		}
	} else {
		if config.verbose {
			message += fmt.Sprintf("\n  - %d collaborators 👥", len(collaborators))
		} else {
			message += fmt.Sprintf("%d collaborators 👥, ", len(collaborators))
		}
	}
	return message
}

func scanCommunityScore(config config, repoWithOrg string) string {
	// https://docs.github.com/en/rest/reference/metrics#get-community-profile-metrics
	communityProfile := struct {
		Health_percentage int64
	}{}
	client, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	err = client.Get(
		"repos/"+repoWithOrg+"/community/profile",
		&communityProfile)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	message := ""
	if config.verbose {
		message += fmt.Sprintf("\n  - a community profile score of %d 💯", communityProfile.Health_percentage)
	} else {
		message += fmt.Sprintf("community profile score: %d 💯", communityProfile.Health_percentage)
	}
	return message
}

func getVersion() version {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Cannot read build info")
	}
	revision := "?"
	dirtyBuild := false
	date := time.Now()
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			revision = kv.Value
		case "vcs.time":
			date, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			dirtyBuild = kv.Value == "true"
		}
	}
	return version{revision, date, dirtyBuild}
}
