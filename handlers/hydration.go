package handlers

import (
	"encoding/json"
	"log"
	"mongoDbTest/services"
	"net/http"
	"time"
)

//Hydrations is hydrationsController
type Hydrations struct {
	service services.Hydrations
	l       *log.Logger
}

func NewHydrationController(l *log.Logger, service services.Hydrations) *Hydrations {
	h := &Hydrations{l: l, service: service}
	return h
}

// GetHydrations get all data
func (h *Hydrations) GetHydrations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hydration, err := h.service.GetHydrations(time.Time{}, time.Now(), 1, 200)
		if err != nil {
			http.Error(w, "handlers.Hydrations.GetHydrations error", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(hydration)
	}
}

//CreateHydration create hydration
// func CreateHydration(s *server.Server) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// decoder := json.NewDecoder(r.Body)
// 		// var body Hydration
// 		// err := decoder.Decode(&body)
// 		// if err != nil {
// 		// 	panic(err)
// 		// }
// 		// validationError := validate.Struct(body)
// 		// fmt.Println(validationError)
// 		// fmt.Println(body)
// 		// fmt.Fprintf(w, "delete")

// 		fmt.Println("Create")
// 		fmt.Fprintf(w, "Create")
// 	}
// }
