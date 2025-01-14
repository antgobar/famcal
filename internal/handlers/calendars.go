package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/antgobar/famcal/internal/googleprovider"
)

func (h Handler) getCalendars(w http.ResponseWriter, r *http.Request) {
	token, err := GetOauth2CookieValue(r)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	calendar, err := googleprovider.NewCalendar(*token, h.config.GoogleOauth2Config)
	if err != nil {
		log.Printf("Error creating calendar service: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calendars, err := calendar.GetCalendars()
	if err != nil {
		log.Printf("Error fetching calendar: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calendars)
}
