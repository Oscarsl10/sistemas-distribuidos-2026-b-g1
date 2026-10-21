// Command demo runs the model of csp-auth-api on :8081 with a stub of the next hop on :8082,
// so the traceparent, the logs and /metrics can be seen with curl.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	instrumentation "csp-auth-api-instrumentation"
)

func main() {
	logger := instrumentation.NewLogger(os.Stdout, "csp-auth-api", "3.0.0", "dev")

	go func() {
		nextHop := http.NewServeMux()
		nextHop.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
			logger.Info("next hop received the message",
				"event", "nexthop.received", "traceparent", r.Header.Get("traceparent"))
			w.WriteHeader(http.StatusAccepted)
		})
		log.Fatal(http.ListenAndServe("127.0.0.1:8082", nextHop))
	}()

	client := &http.Client{Timeout: 3 * time.Second}
	server := instrumentation.NewServer(logger, instrumentation.NewMetrics(), client, "http://127.0.0.1:8082/")
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", server))
}
