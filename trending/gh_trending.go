package trending

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

type GithubTrending struct {
	Host        string
	RequestPath string
	Language    string
	DataRange   string
	Proxy       *url.URL
}

type TrendingInfo struct {
	Title       string
	Description string
	UrlAddr     string
	Forks       string
	Stars       string
	StarsDay    string
}

const (
	DAILY   = "daily"
	WEEKLY  = "weekly"
	MONTHLY = "monthly"
)

var (
	cyan  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
	green = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#32CD32"))
	gray  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#696969"))
	gold  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#B8860B"))
)

var proxy, _ = url.Parse("http://localhost:7897")

func (gh *GithubTrending) getUrlPath() string {

	path := gh.Host + "/" + gh.RequestPath + "/" + gh.Language
	if gh.DataRange != "" {
		path += "?since=" + gh.DataRange
	}
	return path
}

func (gh *GithubTrending) Query() ([]*TrendingInfo, error) {
	trends := []*TrendingInfo{}

	requstPath := gh.getUrlPath()
	log.Infof("request path: %s", requstPath)

	cusTransport := &http.Transport{}

	if gh.Proxy != nil {
		cusTransport.Proxy = http.ProxyURL(gh.Proxy)
	}

	client := &http.Client{
		Transport: cusTransport,
		Timeout:   time.Second * 20,
	}

	response, err := client.Get(requstPath)

	if err != nil {
		log.Errorf("query body error. %v", err)
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Errorf("create document error %v", err)
	}
	pattern, _ := regexp.Compile("[\n ]")
	numberPattern, _ := regexp.Compile("[0-9]{1,}")
	doc.Find(".Box-row").Each(func(i int, s *goquery.Selection) {
		aLink := s.Find("h2 a")
		// repo的title
		title := aLink.Text()
		res := pattern.ReplaceAllString(title, "")

		// repo 链接
		aHref, _ := aLink.Attr("href")

		// repo的description
		description := s.Find("p:nth-child(3)").Text()
		description = strings.TrimFunc(description, func(r rune) bool {
			if r == '\n' || r == ' ' {
				return true
			}
			return false
		})
		// repo fork 数量
		forkNumber := s.Find("div:nth-child(4)").Find("a:nth-child(2)").Text()
		forkNumber = pattern.ReplaceAllString(forkNumber, "")
		// repo start 数量
		starNumber := s.Find("div:nth-child(4)").Find("a:first-of-type").Text()
		starNumber = pattern.ReplaceAllString(starNumber, "")

		// stars of today
		starsDay := s.Find("div:nth-child(4)").Find("span:nth-of-type(3)").Text()
		starsNum := numberPattern.FindString(starsDay)
		log.Debugf("search startDay content: %s, number: %s", starsDay, starsNum)

		info := TrendingInfo{
			Title:       res,
			Description: description,
			UrlAddr:     gh.Host + aHref,
			Forks:       forkNumber,
			Stars:       starNumber,
			StarsDay:    starsNum,
		}
		trends = append(trends, &info)
		log.Debugf("query title:  %s, href: %s, descr: %s, fork %s, starNumber: %s", res, aHref, description, forkNumber, starNumber)
	})

	return trends, nil
}

var JavaDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "java",
	DataRange:   DAILY,
	Proxy:       proxy,
}

var GoDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "go",
	DataRange:   DAILY,
	Proxy:       proxy,
}

var PythonDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "python",
	DataRange:   DAILY,
	Proxy:       proxy,
}

var NodeJsDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "javascript",
	DataRange:   DAILY,
	Proxy:       proxy,
}

func GhTrendingQuery(language, format string) {
	infos := []*TrendingInfo{}
	switch language {
	case "java":
		infos, _ = JavaDefaultGHTrending.Query()
	case "python":
		infos, _ = PythonDefaultGHTrending.Query()
	case "go":
		infos, _ = GoDefaultGHTrending.Query()
	case "javascript":
		infos, _ = NodeJsDefaultGHTrending.Query()
	default:

	}

	switch format {
	case "json":
		data, _ := json.MarshalIndent(infos, "", " ")
		fmt.Fprintf(os.Stdout, "%s", string(data))
	case "table":
		for _, info := range infos {
			fmt.Printf("Repo:    %s  |  language: %s  |  Stars:  %s  |  forks:  %s  |  Stars Today:  %s \n", cyan.Render(info.Title), cyan.Render(language), cyan.Render(info.Stars), cyan.Render(info.Forks), cyan.Render(info.StarsDay))
			fmt.Printf("Desc:    %s\n", green.Render(info.Description))
			fmt.Printf("Link:    %s \n\n", gold.Render(info.UrlAddr))

		}
	default:

	}
}
