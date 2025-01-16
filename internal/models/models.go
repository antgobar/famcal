package models

type Member struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Colour int    `json:"colour"`
}

type CalendarDetails struct {
	Id          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
	Location    string `json:"location,omitempty"`
}

type Event struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}
