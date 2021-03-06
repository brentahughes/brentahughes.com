package webserver

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/bah2830/brentahughes.com/repo"
	"github.com/spf13/viper"
)

var (
	templatePath  = "templates/"
	indexTemplate *template.Template
)

type Webserver struct {
	repoClient *repo.RepoClient
}

type Page struct {
	Title                string
	Repos                []*repo.Repo
	Contributions        []*repo.Repo
	Name                 string
	Email                string
	PhoneNumber          string
	SocialIcons          []SocialIcon
	ProjectSource        string
	ProjectLocation      string
	ProjectLocationLower string
}

type SocialIcon struct {
	Site string
	URL  string
}

func init() {
	indexTemplate, _ = template.ParseFiles(templatePath + "index.html")
}

func GetWebserver(c *repo.RepoClient) *Webserver {
	return &Webserver{
		repoClient: c,
	}
}

func (w *Webserver) Start() {
	http.HandleFunc("/favicon.ico", w.faviconHandler)
	http.HandleFunc("/", w.indexHandler)

	// Setup file server for html resources
	fs := http.FileServer(http.Dir("content"))
	http.Handle("/content/", http.StripPrefix("/content/", fs))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setUserDetails(p *Page) {
	p.Title = viper.GetString("site_title")
	p.Name = viper.GetString("name")
	p.Email = viper.GetString("email")
	p.PhoneNumber = viper.GetString("phone")
	p.ProjectSource = viper.GetString("project_source")

	urlParts, _ := url.Parse(p.ProjectSource)
	projectIcon := strings.TrimPrefix(urlParts.Host, "www.")
	parts := strings.Split(projectIcon, ".")
	projectIcon = parts[0]
	p.ProjectLocation = strings.Title(projectIcon)
	p.ProjectLocationLower = strings.ToLower(projectIcon)
}

func getSocialIcons() []SocialIcon {
	icons := []SocialIcon{}
	links := viper.GetStringSlice("social_links")
	for _, link := range links {
		parts, err := url.Parse(link)
		if err != nil {
			fmt.Println(err)
		}

		host := strings.Replace(parts.Host, "www.", "", -1)
		hostParts := strings.Split(host, ".")

		icon := SocialIcon{
			Site: hostParts[0],
			URL:  link,
		}

		icons = append(icons, icon)
	}

	return icons
}

func (s *Webserver) faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "content/favicon.ico")
}

func (s *Webserver) indexHandler(w http.ResponseWriter, r *http.Request) {
	originalRepos := make([]*repo.Repo, 0)
	contributions := make([]*repo.Repo, 0)

	repos, _ := s.repoClient.GetRepos(false)
	for _, r := range repos {
		if r.Contribution {
			contributions = append(contributions, r)
		} else {
			originalRepos = append(originalRepos, r)
		}
	}

	p := Page{
		Repos:         originalRepos,
		Contributions: contributions,
	}

	p.SocialIcons = getSocialIcons()

	setUserDetails(&p)

	indexTemplate.Execute(w, p)
}
