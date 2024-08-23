package trending

import (
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type GithubTrending struct {
}

func Parse() {
	requstPath := "https://github.com/trending/java?since=daily"

	proxy, _ := url.Parse("http://localhost:7897")
	cusTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	client := &http.Client{
		Transport: cusTransport,
		Timeout:   time.Second * 10,
	}

	response, err := client.Get(requstPath)

	if err != nil {
		log.Errorf("query body error. %v", err)
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)

	if err != nil {
		log.Errorf("create document error %v", err)
	}

	doc.Find(".Box-row").Each(func(i int, s *goquery.Selection) {
		atext := s.Find("h2 a").Text()
		regex, _ := regexp.Compile("[\n ]")
		res := regex.ReplaceAllString(atext, "")
		//text := s.Find(".text-normal").Text()

		log.Infof("query atext:  %s", res)
	})
}
