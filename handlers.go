package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func CreatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID int     `json:"order_id"`
		Amount  float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id := generateID()
	payment := &Payment{
		ID:      id,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Status:  "pending",
	}
	SavePayment(payment)

	go simulateWebhook(payment)

	resp := map[string]string{
		"payment_id":  id,
		"payment_url": "http://fake-payments.local/checkout/" + id,
	}
	log.Printf("[create] Payment created: %+v", payment)
	json.NewEncoder(w).Encode(resp)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, ok := GetPayment(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payment
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid webhook", http.StatusBadRequest)
		return
	}

	UpdatePaymentStatus(payload.ID, payload.Status)
	log.Printf("[webhook] Payment %s updated to %s", payload.ID, payload.Status)
	w.WriteHeader(http.StatusOK)
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func simulateWebhook(p *Payment) {
	time.Sleep(3 * time.Second)

	p.Status = "succeeded"
	body, _ := json.Marshal(p)

	resp, err := http.Post("http://localhost:8080/webhook", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("[webhook] Error simulating webhook: %v", err)
	} else {
		log.Printf("[webhook] Simulated webhook for payment %s (HTTP %d)", p.ID, resp.StatusCode)
	}
}
