package url_scrapper

import (
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 1 * time.Second,
}

func UrlScrapper() map[string]string {
	numWorkers := 3
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
	jobs, outChannel := make(chan string, len(urls)), make(chan string, 2*len(urls))
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()
	urlMapResponses := make(map[string]string)
	for i := 0; i < numWorkers; i++ {
		go func() {
			for url := range jobs {
				workerFunc(url, outChannel)
			}
		}()
	}
	for i := 0; i < len(urls); i++ {
		urlMapResponses[<-outChannel] = <-outChannel
	}
	close(outChannel)
	return urlMapResponses
}

func workerFunc(url string, outChan chan<- string) {
	resp, err := httpClient.Get(url)
	if err != nil {
		outChan <- url
		outChan <- url + " url - not ok"
		return
	}
	defer func() {
		currentBodyErr := resp.Body.Close()
		if currentBodyErr != nil {
			outChan <- url
			outChan <- url + " url - not ok"
			return
		}
	}()
	if resp.StatusCode != 200 {
		outChan <- url
		outChan <- url + " url - not ok"
		return
	}
	outChan <- url
	outChan <- url + " url - ok"
}
