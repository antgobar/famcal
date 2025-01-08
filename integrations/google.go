package integrations

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func HandleRequestToken(authCode string) (*oauth2.Token, error) {
	config, err := getConfigFromCredentials()
	if err != nil {
		return nil, err
	}
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token, nil
}

func getConfigFromCredentials() (*oauth2.Config, error) {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GoogleAuthUrl() (string, error) {
	config, err := getConfigFromCredentials()
	if err != nil {
		return "", err
	}
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
}

func CalendarService(tokenStr string) (*calendar.Service, error) {
	token := &oauth2.Token{
		AccessToken: tokenStr,
	}

	config, err := getConfigFromCredentials()
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background(), token)

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

type Event struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func GetCalendarEvents(srv *calendar.Service, calendarId string, n int64) ([]Event, error) {
	calId := calendarOriginIdFromId(calendarDetailsStore, calendarId)
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(calId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(n).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	if len(events.Items) == 0 {
		return []Event{}, nil
	}

	var calEvents = []Event{}
	for _, item := range events.Items {
		startDate := item.Start.DateTime
		if startDate == "" {
			startDate = item.Start.Date
		}
		endDate := item.End.DateTime
		if endDate == "" {
			endDate = item.End.Date
		}
		event := Event{item.Summary, startDate, endDate}
		calEvents = append(calEvents, event)
	}
	return calEvents, nil
}

type CalendarDetails struct {
	Id          int `json:"id"`
	OriginId    string
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
	Location    string `json:"location,omitempty"`
}

type Calendars []CalendarDetails

func calendarOriginIdFromId(calendars []*CalendarDetails, id string) string {
	for _, calendar := range calendars {
		calId, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		if calendar.Id == calId {
			return calendar.OriginId
		}
	}
	return ""
}

var calendarDetailsStore = []*CalendarDetails{}

func GetCalendars(srv *calendar.Service) ([]*CalendarDetails, error) {
	calendarDetailsStore = []*CalendarDetails{}

	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	for i, item := range calendars.Items {
		calendar := toCalendarDetails(item, i)
		calendarDetailsStore = append(calendarDetailsStore, calendar)
	}
	return calendarDetailsStore, nil
}

func toCalendarDetails(entry *calendar.CalendarListEntry, id int) *CalendarDetails {
	return &CalendarDetails{
		Id:          id,
		OriginId:    entry.Id,
		Summary:     entry.Summary,
		Description: entry.Description,
		TimeZone:    entry.TimeZone,
		Location:    entry.Location,
	}
}
