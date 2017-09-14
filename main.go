package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getHostID() int {
	rand.Seed(time.Now().UTC().UnixNano())
	out := rand.Int()
	return out
}

// Output json struct
type Output struct {
	Hello   string `json:"hello"`
	NodeID  string `json:"node_id"`
	Message string `json:"message"`
}

func main() {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	hostID := getHostID()

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		msg := Output{
			Hello:   "World!",
			NodeID:  strconv.Itoa(hostID),
			Message: getenv("ENV_MESSAGE", "Hello"),
		}
		output, err := json.Marshal(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"error": true}`)
			return
		}
		log.Printf("Serving request to %s\n", r.URL)
		io.WriteString(w, string(output))
	}))

	addr := getenv("PORT", "8082")
	srv := &http.Server{Addr: fmt.Sprintf(":%s", addr), Handler: mux}
	log.Println("Listening on port", addr)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")

}
