package url_scrapper

import (
	"net/http"
	"sync"
	"time"
)

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

func UrlScrapper() map[string]string {
	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://ееее",
	}
	urlMapResponses := make(map[string]string)
	wg := &sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := httpClient.Get(url)
			if err != nil {
				urlMapResponses[url] = url + " url - not ok"
				return
			}
			defer func() {
				currentBodyErr := resp.Body.Close()
				if currentBodyErr != nil {
					urlMapResponses[url] = url + " url - not ok"
					return
				}
			}()
			if resp.StatusCode != 200 {
				urlMapResponses[url] = url + " url - not ok"
				return
			}
			urlMapResponses[url] = url + " url - ok"
		}()
	}
	wg.Wait()
	return urlMapResponses
}
