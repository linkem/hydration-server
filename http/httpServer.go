package httpServer

import (
	"fmt"
	"log"
	"mongoDbTest/handlers"
	"mongoDbTest/middleware"
	"mongoDbTest/models"
	"mongoDbTest/services"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var logger *log.Logger = log.New(os.Stdout, "Http ", log.LstdFlags)

func routes(hydrationService services.Hydrations) *mux.Router {
	router := mux.NewRouter()
	controller := handlers.NewHydrationController(
		log.New(os.Stdout, "hydration-api ", log.LstdFlags),
		hydrationService)

	router.Use(middleware.LoggingMiddleware, middleware.HeadersMiddleware)

	router.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, time.Now().UTC())
	})

	router.HandleFunc("/hydration", controller.GetHydrations()).Methods(http.MethodGet)
	return router
}

func InitHTTPServer(c models.Config, hydrationService func() services.Hydrations) {
	if c.Server.Enabled {
		logger.Printf("Start HTTP at port %s\n", c.Server.Port)
		router := routes(hydrationService())
		go func() {
			log.Fatal(http.ListenAndServe(":"+c.Server.Port, router))
		}()
	}
}
