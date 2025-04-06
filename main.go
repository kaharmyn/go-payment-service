package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	// Routes
	r.Post("/create-payment", CreatePaymentHandler)
	r.Post("/webhook", WebhookHandler)
	r.Get("/status/{id}", StatusHandler)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting payment-service on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
