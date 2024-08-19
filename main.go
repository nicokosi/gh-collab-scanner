package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	// config := parseFlags()
	// if config.version {
	// 	version := getVersion()
	// 	dirty := ""
	// 	if version.dirty {
	// 		dirty = "(dirty)"
	// 	}
	// 	fmt.Printf("Commit %s (%s) %s\n", version.commit, version.date, dirty)
	// } else if len(config.org) > 0 || len(config.user) > 0 {
	// 	repos, error := getRepos(config)
	// 	if error != nil {
	// 		fmt.Print(error)
	// 		os.Exit(2)
	// 	}
	// 	for _, repo := range repos {
	// 		repoMessage, repo, validRepo := scanRepo(config, repo.Full_name)
	// 		if validRepo {
	// 			fmt.Print(repoMessage)
	// 			collaboratorsMessage := scanCollaborators(config, repo.Full_name)
	// 			fmt.Print(collaboratorsMessage)
	// 			if strings.Compare(repo.Visibility, "public") == 0 {
	// 				communityScoreMessage := scanCommunityScore(config, repo.Full_name)
	// 				fmt.Print(communityScoreMessage)
	// 			}
	// 		}
	// 		fmt.Println()
	// 	}
	// } else {
	// 	repoWithOrg, error := getRepo(config)
	// 	if error != nil {
	// 		fmt.Print(error)
	// 		if strings.Contains(error.Error(), "none of the git remotes configured for this repository point to a known GitHub host") {
	// 			print("If current folder is related to a GitHub repository, please check 'gh auth status' and 'gh config list'.")
	// 		}
	// 		os.Exit(1)
	// 	}
	// 	repoMessage, repo, validRepo := scanRepo(config, repoWithOrg)
	// 	if validRepo {
	// 		fmt.Print(repoMessage)
	// 		collaboratorsMessage := scanCollaborators(config, repoWithOrg)
	// 		fmt.Print(collaboratorsMessage)
	// 		if !repo.Fork && strings.Compare(repo.Visibility, "public") == 0 {
	// 			communityScoreMessage := scanCommunityScore(config, repoWithOrg)
	// 			fmt.Print(communityScoreMessage)
	// 		}
	// 		fmt.Println()
	// 	}
	// }

	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "City", Width: 10},
		{Title: "Country", Width: 10},
		{Title: "Population", Width: 10},
	}

	rows := []table.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
		{"7", "Cairo", "Egypt", "21,750,020"},
		{"8", "Beijing", "China", "21,333,332"},
		{"9", "Mumbai", "India", "20,961,472"},
		{"10", "Osaka", "Japan", "19,059,856"},
		{"11", "Chongqing", "China", "16,874,740"},
		{"12", "Karachi", "Pakistan", "16,839,950"},
		{"13", "Istanbul", "Turkey", "15,636,243"},
		{"14", "Kinshasa", "DR Congo", "15,628,085"},
		{"15", "Lagos", "Nigeria", "15,387,639"},
		{"16", "Buenos Aires", "Argentina", "15,369,919"},
		{"17", "Kolkata", "India", "15,133,888"},
		{"18", "Manila", "Philippines", "14,406,059"},
		{"19", "Tianjin", "China", "14,011,828"},
		{"20", "Guangzhou", "China", "13,964,637"},
		{"21", "Rio De Janeiro", "Brazil", "13,634,274"},
		{"22", "Lahore", "Pakistan", "13,541,764"},
		{"23", "Bangalore", "India", "13,193,035"},
		{"24", "Shenzhen", "China", "12,831,330"},
		{"25", "Moscow", "Russia", "12,640,818"},
		{"26", "Chennai", "India", "11,503,293"},
		{"27", "Bogota", "Colombia", "11,344,312"},
		{"28", "Paris", "France", "11,142,303"},
		{"29", "Jakarta", "Indonesia", "11,074,811"},
		{"30", "Lima", "Peru", "11,044,607"},
		{"31", "Bangkok", "Thailand", "10,899,698"},
		{"32", "Hyderabad", "India", "10,534,418"},
		{"33", "Seoul", "South Korea", "9,975,709"},
		{"34", "Nagoya", "Japan", "9,571,596"},
		{"35", "London", "United Kingdom", "9,540,576"},
		{"36", "Chengdu", "China", "9,478,521"},
		{"37", "Nanjing", "China", "9,429,381"},
		{"38", "Tehran", "Iran", "9,381,546"},
		{"39", "Ho Chi Minh City", "Vietnam", "9,077,158"},
		{"40", "Luanda", "Angola", "8,952,496"},
		{"41", "Wuhan", "China", "8,591,611"},
		{"42", "Xi An Shaanxi", "China", "8,537,646"},
		{"43", "Ahmedabad", "India", "8,450,228"},
		{"44", "Kuala Lumpur", "Malaysia", "8,419,566"},
		{"45", "New York City", "United States", "8,177,020"},
		{"46", "Hangzhou", "China", "8,044,878"},
		{"47", "Surat", "India", "7,784,276"},
		{"48", "Suzhou", "China", "7,764,499"},
		{"49", "Hong Kong", "Hong Kong", "7,643,256"},
		{"50", "Riyadh", "Saudi Arabia", "7,538,200"},
		{"51", "Shenyang", "China", "7,527,975"},
		{"52", "Baghdad", "Iraq", "7,511,920"},
		{"53", "Dongguan", "China", "7,511,851"},
		{"54", "Foshan", "China", "7,497,263"},
		{"55", "Dar Es Salaam", "Tanzania", "7,404,689"},
		{"56", "Pune", "India", "6,987,077"},
		{"57", "Santiago", "Chile", "6,856,939"},
		{"58", "Madrid", "Spain", "6,713,557"},
		{"59", "Haerbin", "China", "6,665,951"},
		{"60", "Toronto", "Canada", "6,312,974"},
		{"61", "Belo Horizonte", "Brazil", "6,194,292"},
		{"62", "Khartoum", "Sudan", "6,160,327"},
		{"63", "Johannesburg", "South Africa", "6,065,354"},
		{"64", "Singapore", "Singapore", "6,039,577"},
		{"65", "Dalian", "China", "5,930,140"},
		{"66", "Qingdao", "China", "5,865,232"},
		{"67", "Zhengzhou", "China", "5,690,312"},
		{"68", "Ji Nan Shandong", "China", "5,663,015"},
		{"69", "Barcelona", "Spain", "5,658,472"},
		{"70", "Saint Petersburg", "Russia", "5,535,556"},
		{"71", "Abidjan", "Ivory Coast", "5,515,790"},
		{"72", "Yangon", "Myanmar", "5,514,454"},
		{"73", "Fukuoka", "Japan", "5,502,591"},
		{"74", "Alexandria", "Egypt", "5,483,605"},
		{"75", "Guadalajara", "Mexico", "5,339,583"},
		{"76", "Ankara", "Turkey", "5,309,690"},
		{"77", "Chittagong", "Bangladesh", "5,252,842"},
		{"78", "Addis Ababa", "Ethiopia", "5,227,794"},
		{"79", "Melbourne", "Australia", "5,150,766"},
		{"80", "Nairobi", "Kenya", "5,118,844"},
		{"81", "Hanoi", "Vietnam", "5,067,352"},
		{"82", "Sydney", "Australia", "5,056,571"},
		{"83", "Monterrey", "Mexico", "5,036,535"},
		{"84", "Changsha", "China", "4,809,887"},
		{"85", "Brasilia", "Brazil", "4,803,877"},
		{"86", "Cape Town", "South Africa", "4,800,954"},
		{"87", "Jiddah", "Saudi Arabia", "4,780,740"},
		{"88", "Urumqi", "China", "4,710,203"},
		{"89", "Kunming", "China", "4,657,381"},
		{"90", "Changchun", "China", "4,616,002"},
		{"91", "Hefei", "China", "4,496,456"},
		{"92", "Shantou", "China", "4,490,411"},
		{"93", "Xinbei", "Taiwan", "4,470,672"},
		{"94", "Kabul", "Afghanistan", "4,457,882"},
		{"95", "Ningbo", "China", "4,405,292"},
		{"96", "Tel Aviv", "Israel", "4,343,584"},
		{"97", "Yaounde", "Cameroon", "4,336,670"},
		{"98", "Rome", "Italy", "4,297,877"},
		{"99", "Shijiazhuang", "China", "4,285,135"},
		{"100", "Montreal", "Canada", "4,276,526"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
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
	client, err := api.DefaultRESTClient()
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
	currentRepo, error := repository.Current()
	if error != nil {
		return "", error
	}
	return currentRepo.Owner + "/" + currentRepo.Name, nil
}

func scanRepo(config config, repoWithOrg string) (message string, repository repo, validRepo bool) {
	// https://docs.github.com/en/rest/reference/repos#get-a-repository-readme
	readme := struct {
		Name string
	}{}
	client, err := api.DefaultRESTClient()
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
	client, err := api.DefaultRESTClient()
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
	client, err := api.DefaultRESTClient()
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
