package http_handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func TestHttpHandler_HttpHandler(t *testing.T) {
	go StartServer()
	time.Sleep(1 * time.Second)
	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body := []byte(`{
						"key": "key",
						"value": "value"
						}`)

			resp, err := httpClient.Post("http://localhost:8080/put", "application/json", bytes.NewBuffer(body))
			defer resp.Body.Close()
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	//
	resp, err := httpClient.Get("http://localhost:8080/put-counter")
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Error("status code is not 200")
	}
	var putCounter int64
	_ = json.NewDecoder(resp.Body).Decode(&putCounter)
	if putCounter != 100 {
		t.Error("counter is not 100")
	}

	reqBody, err := json.Marshal(`
				"key": "key",
			`)
	if err != nil {
		t.Error(err)
	}
	responseBody := bytes.NewBuffer(reqBody)
	resp, err = httpClient.Post("http://localhost:8080/get", "application/json", responseBody)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
	}

	var objectFromCache RequestObject
	_ = json.NewDecoder(resp.Body).Decode(&objectFromCache)

}
