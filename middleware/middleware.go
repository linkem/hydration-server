package middleware

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var middlewareLog *log.Logger = log.New(os.Stdout, "MiddleWare ", log.LstdFlags)

// LoggingMiddleware log request body
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middlewareLog.Printf("Route:[%s] %s", r.Method, r.URL.Path)
		middlewareLog.Printf("RawQuery: %s", r.URL.RawQuery)
		middlewareLog.Printf("Query: %v", r.URL.Query())
		middlewareLog.Printf("# of params: %d", len(r.URL.Query()))

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			middlewareLog.Printf("Error readying body; \nError: %s", err.Error())
		}
		middlewareLog.Printf("Body: %v", body)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// HeadersMiddleware set headers
func HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
