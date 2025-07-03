package http_handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"interview_tasks/custom_cache"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type HttpHandler struct {
	cache            CustomCacheInterface
	putObjectCounter int64
	getObjectCounter int64
}

// StartServer http handler
// 1. parse json
// 2. put in map
// 3. post - put data in map, get - get data from map based on key
// 4. count number of requests
func StartServer() {
	router := mux.NewRouter()
	customCache := custom_cache.NewCustomCache()
	httpServer := newHttpHandler(customCache)

	router = mux.NewRouter()

	router.HandleFunc("/put", httpServer.PutObjectInMap).Methods("POST")
	router.HandleFunc("/get", httpServer.GetObjectFromMap).Methods("POST")
	router.HandleFunc("/get-counter", httpServer.GetObjectCounter).Methods("GET")
	router.HandleFunc("/put-counter", httpServer.PutObjectCounter).Methods("GET")

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on localhost:8080")

	// start HTTP server
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")
}

func newHttpHandler(cache CustomCacheInterface) *HttpHandler {
	return &HttpHandler{
		cache: cache,
	}
}
func (h *HttpHandler) PutObjectInMap(w http.ResponseWriter, r *http.Request) {
	var request RequestObject
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w)
		return
	}

	atomic.AddInt64(&h.putObjectCounter, 1) // Atomically increment the counter
	h.cache.Put(request.Key, request.Value)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("success")
}

func (h *HttpHandler) GetObjectFromMap(w http.ResponseWriter, r *http.Request) {
	var request RequestKey
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w)
		return
	}

	atomic.AddInt64(&h.getObjectCounter, 1)
	objectFromCache, err := h.cache.Get(request.Key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(objectFromCache)
}

func (h *HttpHandler) GetObjectCounter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(atomic.LoadInt64(&h.getObjectCounter))
}

func (h *HttpHandler) PutObjectCounter(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(atomic.LoadInt64(&h.putObjectCounter))
}
