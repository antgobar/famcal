package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/antgobar/famcal/internal/googleprovider"
)

type eventsQuery struct {
	nEvents    int64
	calendarId string
}

func getEventsParams(r *http.Request) (*eventsQuery, error) {
	rawQuery := r.URL.RawQuery
	decodedQuery, err := url.QueryUnescape(rawQuery)
	if err != nil {
		return nil, errors.New("error decoding query parameters")
	}
	r.URL.RawQuery = decodedQuery

	calendarID := r.URL.Query().Get("calendarId")
	if calendarID == "" {
		return nil, errors.New("missing parameter: calendarID")
	}

	nEvents := r.URL.Query().Get("nEvents")
	if nEvents == "" {
		return nil, errors.New("missing parameter: nEvents")
	}

	nEventsInt, err := strconv.ParseInt(nEvents, 10, 64)
	if err != nil {
		return nil, err
	}

	return &eventsQuery{nEventsInt, calendarID}, nil
}

func (handler Handler) getEvents(w http.ResponseWriter, r *http.Request) {
	token, err := GetOauth2CookieValue(r)
	if err != nil {
		log.Printf("Error getting token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	calendar, err := googleprovider.NewCalendar(*token, handler.config.GoogleOauth2Config)
	if err != nil {
		log.Printf("Error creating calendar service: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	params, err := getEventsParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := calendar.GetCalendarEvents(params.calendarId, params.nEvents)
	if err != nil {
		log.Printf("Error fetching events: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
