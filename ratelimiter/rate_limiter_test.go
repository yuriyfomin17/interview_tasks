package ratelimiter

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

var httpClient = &http.Client{
	Timeout: 2 * time.Second,
}

func TestRateLimiter_ProcessTestRateLimiter(t *testing.T) {
	port := "8081"
	go func() {
		TestRateLimiter(10, port)
	}()
	time.Sleep(100 * time.Millisecond)
	testRateLimiterAndCheckExpectedStatusCode(t, port, 10, http.StatusTooManyRequests)
	time.Sleep(time.Second)
	testRateLimiterAndCheckExpectedStatusCode(t, port, 5, http.StatusOK)

}

func testRateLimiterAndCheckExpectedStatusCode(t *testing.T, port string, numRequests, expectedStatusCode int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := httpClient.Get(fmt.Sprintf("http://localhost:%s/request", port))
			if err != nil {
				t.Error("Could not get response")
			}
		}()
	}
	wg.Wait()
	resp, err := httpClient.Get(fmt.Sprintf("http://localhost:%s/request", port))
	if err != nil {
		t.Error("Could not get response")
	}
	if resp.StatusCode != expectedStatusCode {
		t.Error("Could not get status code", expectedStatusCode)
	}
}
