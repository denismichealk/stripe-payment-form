package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func main() {
	// This is a public sample test API key.
	// Don’t submit any personally identifiable information in requests made with this key.
	// Sign in to see your own test API key embedded in code samples.
	stripe.Key = "sk_test_51P3DHGP80Ea2hQchCYudTu7mP5TNMEEqQU47aFavIzc4XwjjpBNB79PK5AI0whTVhAmzis1DhcGIl9LPLVHYRW8w00Pr30s8jg"

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/create-checkout-session", createCheckoutSession)
	http.HandleFunc("/session-status", retrieveCheckoutSession)
	addr := "localhost:4242"
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func createCheckoutSession(w http.ResponseWriter, r *http.Request) {
	domain := "http://localhost:4242"
	params := &stripe.CheckoutSessionParams{
		UIMode:    stripe.String("embedded"),
		ReturnURL: stripe.String(domain + "/return.html?session_id={CHECKOUT_SESSION_ID}"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				// Provide the exact Price ID (for example, pr_1234) of the product you want to sell
				Price:    stripe.String("price_1P3DORP80Ea2hQchpxpaEB9j"),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	s, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
	}

	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: s.ClientSecret,
	})
}

func retrieveCheckoutSession(w http.ResponseWriter, r *http.Request) {
	s, _ := session.Get(r.URL.Query().Get("session_id"), nil)

	writeJSON(w, struct {
		Status        string `json:"status"`
		CustomerEmail string `json:"customer_email"`
	}{
		Status:        string(s.Status),
		CustomerEmail: string(s.CustomerDetails.Email),
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}
