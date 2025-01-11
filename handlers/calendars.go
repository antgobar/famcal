package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/antgobar/famcal/googleprovider"
)

func getCalendars(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromCookie(r)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calendar, err := googleprovider.NewCalendar(token)
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
