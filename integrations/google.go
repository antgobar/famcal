package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) (*http.Client, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), token), nil
}

func getToken() (*oauth2.Token, error) {
	tokenFile := "token.json"
	token, err := getTokenFromFile(tokenFile)
	fmt.Println(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

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

func SaveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
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

func CalendarService() (*calendar.Service, error) {
	config, err := getConfigFromCredentials()
	if err != nil {
		return nil, err
	}
	client, err := getClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}

type Event struct {
	Summary   string
	EventDate string
}

func GetNextNCalendarEvents(srv *calendar.Service, calendarId string, n int64) ([]Event, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(n).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		return []Event{}, nil
	}

	var calEvents = []Event{}
	for _, item := range events.Items {
		date := item.Start.DateTime
		if date == "" {
			date = item.Start.Date
		}
		event := Event{item.Summary, date}
		fmt.Println(event)
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

func GetCalendars(srv *calendar.Service) ([]*CalendarDetails, error) {
	calendarDetails := []*CalendarDetails{}
	calendars, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	for _, item := range calendars.Items {
		calendar := toCalendarDetails(item)
		calendarDetails = append(calendarDetails, calendar)
	}
	return calendarDetails, nil
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
