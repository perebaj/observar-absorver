package marinho

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Essay struct {
	Title   string
	Content string
	Date    string
	URL     string
}

type Essays []Essay

type RawEssay struct {
	URL  string
	HTML string
}

type RawEssays []RawEssay

const seed_url = "https://observareabsorver.blogspot.com/"

func ParseHTML2Essay(rawHTML string) (Essays, error) {
	//convert the raw html to a goquery document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error parsing the html body: %v", err)
	}

	var essays Essays
	doc.Find("div.date-outer").Each(func(i int, s *goquery.Selection) {
		dateString := s.Find("h2 > span").Text()
		title := s.Find("h3").Text()
		url, _ := s.Find("h3 > a").Attr("href")
		// any type of element that has text
		content := s.Text()
		essay := Essay{
			Title:   title,
			URL:     url,
			Date:    dateString,
			Content: content,
		}
		essays = append(essays, essay)
	})

	return essays, nil
}

func FetchPageHTML() (RawEssays, error) {
	var rawEssays RawEssays
	url := seed_url
	for {
		slog.Info("fetching", "url", url)
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching next link: %v", err)
		}

		defer res.Body.Close()
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, fmt.Errorf("error parsing the html body: %v", err)
		}

		html, err := doc.Html()
		if err != nil {
			return nil, fmt.Errorf("error getting the html content: %v", err)
		}

		rawEssay := RawEssay{
			URL:  url,
			HTML: html,
		}

		rawEssays = append(rawEssays, rawEssay)
		url = doc.Find("a.blog-pager-older-link").AttrOr("href", "")
		if url == "" {
			break
		}
	}
	return rawEssays, nil
}
