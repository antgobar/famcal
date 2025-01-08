package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/antgobar/famcal/integrations"
)

func getCalendars(w http.ResponseWriter, r *http.Request) {
	service, err := integrations.CalendarService()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error creating calendar service: %v", err)
		return
	}

	calendars, err := integrations.GetCalendars(service)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error fetching calendar: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendars)
}
