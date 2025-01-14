package googleprovider

import (
	"context"
	"time"

	"github.com/antgobar/famcal/internal/models"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendar struct {
	service *calendar.Service
}

func NewCalendar(tokenStr string, config *oauth2.Config) (*GoogleCalendar, error) {
	token := &oauth2.Token{
		AccessToken: tokenStr,
	}
	client := config.Client(context.Background(), token)
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &GoogleCalendar{srv}, nil
}

func (cal GoogleCalendar) GetCalendarEvents(calendarId string, n int64) ([]models.Event, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := cal.service.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(n).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}
	if len(events.Items) == 0 {
		return []models.Event{}, nil
	}

	var calEvents = []models.Event{}
	for _, item := range events.Items {
		startDate := item.Start.DateTime
		if startDate == "" {
			startDate = item.Start.Date
		}
		endDate := item.End.DateTime
		if endDate == "" {
			endDate = item.End.Date
		}
		event := models.Event{Summary: item.Summary, Start: startDate, End: endDate}
		calEvents = append(calEvents, event)
	}
	return calEvents, nil
}

type Calendars []models.CalendarDetails

var calendarDetailsStore = []*models.CalendarDetails{}

func (cal GoogleCalendar) GetCalendars() ([]*models.CalendarDetails, error) {
	calendarDetailsStore = []*models.CalendarDetails{}

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

func toCalendarDetails(entry *calendar.CalendarListEntry) *models.CalendarDetails {
	return &models.CalendarDetails{
		Id:          entry.Id,
		Summary:     entry.Summary,
		Description: entry.Description,
		TimeZone:    entry.TimeZone,
		Location:    entry.Location,
	}
}
