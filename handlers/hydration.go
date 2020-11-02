package handlers

import (
	"encoding/json"
	"log"
	"mongoDbTest/services"
	"net/http"
	"strconv"
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
		pageSizeString := r.URL.Query().Get("pageSize")
		pageString := r.URL.Query().Get("page")
		pageSize := 5
		page := 1
		var err error
		if pageString != "" {
			page, err = strconv.Atoi(pageString)
			if err != nil {
				http.Error(w, "handlers.Hydrations.GetHydrations; Error parsing page", http.StatusBadRequest)
				return
			}
		}
		if pageSizeString != "" {
			pageSize, err = strconv.Atoi(pageSizeString)
			if err != nil {
				http.Error(w, "handlers.Hydrations.GetHydrations; Error parsing pageSize", http.StatusBadRequest)
				return
			}
		}
		hydration, err := h.service.GetHydrations(r.Context(), time.Time{}, time.Now(), page, pageSize)
		if err != nil {
			http.Error(w, "handlers.Hydrations.GetHydrations error", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		h.l.Printf("Results: %d", len(*hydration))
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
