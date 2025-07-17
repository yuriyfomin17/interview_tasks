package url_scrapper

import (
	"context"
	"testing"
	"time"
)

func TestUrlScrapper(t *testing.T) {
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan map[string]string, 1)

	expected := map[string]string{
		"http://ozon.ru":                 "http://ozon.ru url - not ok",
		"https://ozon.ru":                "https://ozon.ru url - not ok",
		"http://google.com":              "http://google.com url - ok",
		"http://somesite.com":            "http://somesite.com url - not ok",
		"http://non-existent.domain.tld": "http://non-existent.domain.tld url - not ok",
		"https://ya.ru":                  "https://ya.ru url - ok",
		"http://ya.ru":                   "http://ya.ru url - ok",
		"http://ееее":                    "http://ееее url - not ok",
	}

	go func() {
		result := UrlScrapper()
		done <- result
		close(done)
	}()

	select {
	case <-contextWithTimeout.Done():
		t.Error("Timeout was achieved")
		return
	case actual, ok := <-done:
		if !ok {
			t.Error("Channel was closed unexpectedly")
			return
		}
		if len(actual) != len(expected) {
			t.Errorf("Maps have different sizes. Expected: %d, got: %d", len(expected), len(actual))
			return
		}
		for k, v := range expected {
			if actualValue, exists := actual[k]; !exists || actualValue != v {
				t.Errorf("Mismatch for key %s. Expected: %s, got: %s", k, v, actualValue)
			}
		}
	}
}
