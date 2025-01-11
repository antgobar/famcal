package googleprovider

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendar struct {
	service *calendar.Service
}

func HandleRequestToken(authCode string) (*oauth2.Token, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token, nil
}

func loadConfig() (*oauth2.Config, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{calendar.CalendarScope},
		Endpoint:     google.Endpoint,
	}, nil
}

func AuthUrl() (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", err
	}
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
}

func NewCalendar(tokenStr string) (*GoogleCalendar, error) {
	token := &oauth2.Token{
		AccessToken: tokenStr,
	}

	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background(), token)

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &GoogleCalendar{srv}, nil
}

type Event struct {
	Summary string `json:"summary"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

func (cal GoogleCalendar) GetCalendarEvents(calendarId string, n int64) ([]Event, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := cal.service.Events.List(calendarId).ShowDeleted(false).
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
	Id          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	TimeZone    string `json:"timeZone,omitempty"`
	Location    string `json:"location,omitempty"`
}

type Calendars []CalendarDetails

var calendarDetailsStore = []*CalendarDetails{}

func (cal GoogleCalendar) GetCalendars() ([]*CalendarDetails, error) {
	calendarDetailsStore = []*CalendarDetails{}

	calendars, err := cal.service.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	for _, item := range calendars.Items {
		calendar := toCalendarDetails(item)
		calendarDetailsStore = append(calendarDetailsStore, calendar)
	}
	return calendarDetailsStore, nil
}

func toCalendarDetails(entry *calendar.CalendarListEntry) *CalendarDetails {
	return &CalendarDetails{
		Id:          entry.Id,
		Summary:     entry.Summary,
		Description: entry.Description,
		TimeZone:    entry.TimeZone,
		Location:    entry.Location,
	}
}
