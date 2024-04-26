package marinho

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ScraperEssay is a struct that contains the content of the essay blog post
type ScraperEssay struct {
	Title   string
	Content string
	Date    string
	URL     string
}

// ScraperEssays is a collection of Essay
type ScraperEssays []ScraperEssay

// RawEssay is a struct that contains the content of the html page
type RawEssay struct {
	URL  string
	HTML string
}

// RawEssays is a array of RawEssay
type RawEssays []RawEssay

const seedURL = "https://observareabsorver.blogspot.com/"

// ParseHTML2Essay parses the raw html and returns the essays
func ParseHTML2Essay(rawHTML string) (ScraperEssays, error) {
	//convert the raw html to a goquery document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error parsing the html body: %v", err)
	}

	var essays ScraperEssays
	doc.Find("div.date-outer").Each(func(_ int, s *goquery.Selection) {
		dateString := s.Find("h2 > span").Text()
		title := s.Find("h3").Text()
		url, _ := s.Find("h3 > a").Attr("href")
		// any type of element that has text
		content := s.Text()
		essay := ScraperEssay{
			Title:   title,
			URL:     url,
			Date:    dateString,
			Content: content,
		}
		// if the title or the url is empty, we don't want to append it the final result
		if essay.Title != "" || essay.URL != "" {
			essays = append(essays, essay)
		}
	})

	return essays, nil
}

// FetchPageHTML fetches the html content of the page
func FetchPageHTML() (RawEssays, error) {
	var rawEssays RawEssays
	url := seedURL
	for {
		slog.Info("fetching", "url", url)
		res, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error fetching next link: %v", err)
		}
		func() {
			_ = res.Body.Close()
		}()
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
