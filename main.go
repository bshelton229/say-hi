package main

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func generateID() int {
	rand.Seed(time.Now().UTC().UnixNano())
	out := rand.Int()
	return out
}

// Output json struct
type Output struct {
	Hello             string            `json:"hello"`
	GeneratedID       string            `json:"generated_id"`
	AdditionalMessage string            `json:"additional_message"`
	RequestHeaders    map[string]string `json:"request_headers"`
	SayHiEnv          map[string]string `json:"say_hi_env"`
}

func main() {
	generatedID := generateID()

	mux := http.NewServeMux()

	mux.Handle("/down", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"msg":"DOWN!"}`)
	}))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		headerOut := map[string]string{}
		envOut := map[string]string{}
		cookieReg := regexp.MustCompile("(?i)(cookie)")

		for k, v := range r.Header {
			if !cookieReg.MatchString(k) {
				headerOut[k] = v[0]
			}
		}

		// Build up map of env vars starting with SAY_HI_ENV_
		for _, v := range os.Environ() {
			if strings.HasPrefix(v, "SAY_HI_ENV_") {
				parsed := strings.SplitN(v, "=", 2)
				envOut[parsed[0]] = parsed[1]
			}
		}

		msg := Output{
			Hello:             "World!",
			GeneratedID:       strconv.Itoa(generateID),
			AdditionalMessage: "Added message",
			RequestHeaders:    headerOut,
			SayHiEnv:          envOut,
		}
		output, err := json.Marshal(msg)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"error":true}`)
			return
		}

		log.Printf("Serving request to %s\n", r.URL)
		io.WriteString(w, string(output))
	}))

	var listen string
	flag.StringVar(&listen, "listen", "0.0.0.0:8082", "Listen string")
	flag.Parse()

	srv := &http.Server{Addr: listen, Handler: mux}
	log.Println("Listening on ", listen)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
