package middleware

import (
	"io/ioutil"
	"log"
	"net/http"
)

// LoggingMiddleware log request body
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Printf("Error readying body; \nError: %s", err.Error())
		}
		log.Println(body)
		log.Println()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// HeadersMiddleware set headers
func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
