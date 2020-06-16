package main

import (
	"fmt"
	"log"
	"mongoDbTest/handlers"
	"mongoDbTest/middleware"
	"mongoDbTest/services"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

//Routes configure router
func Routes(hydrationService services.Hydrations) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	controller := handlers.NewHydrationController(
		log.New(os.Stdout, "hydration-api ", log.LstdFlags),
		hydrationService)

	router.Use(middleware.LoggingMiddleware, middleware.HeadersMiddleware)

	router.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().UTC())
	})
	router.HandleFunc("/hydration", controller.GetHydrations()).Methods("GET").Queries()
	return router
}
