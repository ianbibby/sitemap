package sitemap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ianbibby/link"
)

type URL struct {
	Loc string `xml:"loc"`
}

type Sitemap struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []URL    `xml:"url"`
}

func (s Sitemap) String() string {
	b, err := xml.MarshalIndent(&s, "", "  ")
	if err != nil {
		panic(err)
	}

	header := `<?xml version="1.0" encoding="UTF-8"?>`
	return fmt.Sprintf("%s\n%s\n", header, string(b))
}

func Build(url string, maxDepth int) (Sitemap, error) {
	seen := make(map[string]struct{})
	q := make(map[string]struct{})
	nq := map[string]struct{}{
		url: struct{}{},
	}

	s := Sitemap{}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})

		for href, _ := range q {
			if _, ok := seen[href]; ok {
				continue
			}
			seen[href] = struct{}{}

			hrefs, err := get(href)
			if err != nil {
				return Sitemap{}, err
			}
			for _, href := range hrefs {
				nq[href] = struct{}{}
			}
		}
	}

	for href, _ := range seen {
		s.URLs = append(s.URLs, URL{Loc: href})
	}

	return s, nil
}

func get(rootURL string) ([]string, error) {
	resp, err := http.Get(rootURL)
	if err != nil {
		return nil, err
	}

	links, err := link.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var filtered []string
	baseURL := resp.Request.URL
	for _, l := range links {
		u, err := url.Parse(l.Href)
		if err != nil {
			return nil, err
		}

		if u.Scheme != "" && u.Scheme != baseURL.Scheme {
			continue
		}
		if u.Host != "" && u.Host != baseURL.Host {
			continue
		}

		add := (&url.URL{
			Scheme: baseURL.Scheme,
			Host:   baseURL.Host,
			Path:   u.Path,
		}).String()

		filtered = append(filtered, add)
	}

	return filtered, nil
}
