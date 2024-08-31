package trending

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	ForkNumber  string
	StarNumber  string
}

const (
	DAILY   = "daily"
	WEEKLY  = "weekly"
	MONTHLY = "monthly"
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
	regex, _ := regexp.Compile("[\n ]")
	doc.Find(".Box-row").Each(func(i int, s *goquery.Selection) {
		aLink := s.Find("h2 a")
		// repo的title
		title := aLink.Text()
		res := regex.ReplaceAllString(title, "")

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
		forkNumber = regex.ReplaceAllString(forkNumber, "")
		// repe start 数量
		startNumber := s.Find("div:nth-child(4)").Find("a:first-of-type").Text()
		startNumber = regex.ReplaceAllString(startNumber, "")

		info := TrendingInfo{
			Title:       res,
			Description: description,
			UrlAddr:     gh.Host + aHref,
			ForkNumber:  forkNumber,
			StarNumber:  startNumber,
		}
		trends = append(trends, &info)
		log.Debugf("query title:  %s, href: %s, descr: %s, fork %s, starNumber: %s", res, aHref, description, forkNumber, startNumber)
	})

	return trends, nil
}

var JavaDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "java",
	DataRange:   WEEKLY,
	Proxy:       proxy,
}

var GoDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "go",
	DataRange:   WEEKLY,
	Proxy:       proxy,
}

var PythonDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "python",
	DataRange:   WEEKLY,
	Proxy:       proxy,
}

var NodeJsDefaultGHTrending = &GithubTrending{
	Host:        "https://github.com",
	RequestPath: "trending",
	Language:    "javascript",
	DataRange:   WEEKLY,
	Proxy:       proxy,
}
