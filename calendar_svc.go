package garmin

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type CalendarService service

// TODO:
// GET https://connect.garmin.com/calendar-service/events?startDate=2024-08-18&endDate=2024-08-22&eventType=running
// GET https://connect.garmin.com/calendar-service/events?startDate=2024-08-24&pageIndex=1&limit=20&sortOrder=eventDate_asc
// GET https://connect.garmin.com/calendar-service/events?startDate=1995-01-01&endDate=2024-08-23&pageIndex=1&limit=20&sortOrder=eventDate_desc

type EventRequest struct {
	EventType *string
	SortOrder *string
	Start     *time.Time
	End       *time.Time
	Limit     *int
	PageIndex *int
}

type CalendarPreferences struct {
	View               int   `json:"view"`
	ShowEvent          bool  `json:"showEvent"`
	ShowGoal           bool  `json:"showGoal"`
	ShowOptions        bool  `json:"showOptions"`
	ShowWorkout        bool  `json:"showWorkout"`
	ShowWeeklyTotal    bool  `json:"showWeeklyTotal"`
	EventColor         int   `json:"eventColor"`
	GoalColor          int   `json:"goalColor"`
	WorkoutColor       int   `json:"workoutColor"`
	ExpandActivity     bool  `json:"expandActivity"`
	ExpandGoal         bool  `json:"expandGoal"`
	ExpandCalendar     bool  `json:"expandCalendar"`
	ExpandTrainingPlan bool  `json:"expandTrainingPlan"`
	ExpandGroups       bool  `json:"expandGroups"`
	Groups             any   `json:"groups"`
	TrainingPlans      []any `json:"trainingPlans"`
	ActivityTypes      []struct {
		Key   string `json:"key"`
		Color int    `json:"color"`
		Show  bool   `json:"show"`
	} `json:"activityTypes"`
	HealthTypes []struct {
		Key  string `json:"key"`
		Show bool   `json:"show"`
	} `json:"healthTypes"`
	HideJetLag bool `json:"hideJetLag"`
}

func (c *CalendarService) Preferences() (*CalendarPreferences, error) {
	var res CalendarPreferences
	return &res, c.c.apiGet(&res, "/calendar-service/preferences", nil)
}

type Calendar struct {
	StartDayOfMonth      int            `json:"startDayOfMonth"`
	NumOfDaysInMonth     int            `json:"numOfDaysInMonth"`
	NumOfDaysInPrevMonth int            `json:"numOfDaysInPrevMonth"`
	StartDate            string         `json:"startDate"`
	EndDate              string         `json:"endDate"`
	Month                int            `json:"month"`
	Year                 int            `json:"year"`
	CalendarItems        []CalendarItem `json:"calendarItems"`
}

type CalendarItem struct {
	ID                    any    `json:"id"`
	GroupID               any    `json:"groupId"`
	TrainingPlanID        any    `json:"trainingPlanId"`
	ItemType              string `json:"itemType"`
	ActivityTypeID        int    `json:"activityTypeId"`
	WellnessActivityUUID  any    `json:"wellnessActivityUuid"`
	Title                 string `json:"title"`
	Date                  string `json:"date"`
	Duration              any    `json:"duration"`
	Distance              any    `json:"distance"`
	Calories              any    `json:"calories"`
	FloorsClimbed         any    `json:"floorsClimbed"`
	AvgRespirationRate    any    `json:"avgRespirationRate"`
	UnitOfPoolLength      any    `json:"unitOfPoolLength"`
	Weight                any    `json:"weight"`
	Difference            any    `json:"difference"`
	CourseID              any    `json:"courseId"`
	CourseName            any    `json:"courseName"`
	SportTypeKey          any    `json:"sportTypeKey"`
	URL                   string `json:"url"`
	IsStart               any    `json:"isStart"`
	IsRace                bool   `json:"isRace"`
	RecurrenceID          any    `json:"recurrenceId"`
	IsParent              any    `json:"isParent"`
	ParentID              any    `json:"parentId"`
	UserBadgeID           any    `json:"userBadgeId"`
	BadgeCategoryTypeID   any    `json:"badgeCategoryTypeId"`
	BadgeCategoryTypeDesc any    `json:"badgeCategoryTypeDesc"`
	BadgeAwardedDate      any    `json:"badgeAwardedDate"`
	BadgeViewed           any    `json:"badgeViewed"`
	HideBadge             any    `json:"hideBadge"`
	StartTimestampLocal   any    `json:"startTimestampLocal"`
	EventTimeLocal        struct {
		StartTimeHhMm string `json:"startTimeHhMm"`
		TimeZoneID    string `json:"timeZoneId"`
	} `json:"eventTimeLocal"`
	DiveNumber                 any    `json:"diveNumber"`
	MaxDepth                   any    `json:"maxDepth"`
	AvgDepth                   any    `json:"avgDepth"`
	SurfaceInterval            any    `json:"surfaceInterval"`
	ElapsedDuration            any    `json:"elapsedDuration"`
	LapCount                   any    `json:"lapCount"`
	BottomTime                 any    `json:"bottomTime"`
	AtpPlanID                  any    `json:"atpPlanId"`
	WorkoutID                  any    `json:"workoutId"`
	ProtectedWorkoutSchedule   bool   `json:"protectedWorkoutSchedule"`
	ActiveSets                 any    `json:"activeSets"`
	Strokes                    any    `json:"strokes"`
	NoOfSplits                 any    `json:"noOfSplits"`
	MaxGradeValue              any    `json:"maxGradeValue"`
	TotalAscent                any    `json:"totalAscent"`
	DifferenceStress           any    `json:"differenceStress"`
	ClimbDuration              any    `json:"climbDuration"`
	MaxSpeed                   any    `json:"maxSpeed"`
	AverageHR                  any    `json:"averageHR"`
	ActiveSplitSummaryDuration any    `json:"activeSplitSummaryDuration"`
	MaxSplitDistance           any    `json:"maxSplitDistance"`
	MaxSplitSpeed              any    `json:"maxSplitSpeed"`
	Location                   string `json:"location"`
	ShareableEventUUID         string `json:"shareableEventUuid"`
	SplitSummaryMode           any    `json:"splitSummaryMode"`
	CompletionTarget           struct {
		Value    float64 `json:"value"`
		Unit     string  `json:"unit"`
		UnitType string  `json:"unitType"`
	} `json:"completionTarget"`
	WorkoutUUID        any  `json:"workoutUuid"`
	NapStartTimeLocal  any  `json:"napStartTimeLocal"`
	PhasedTrainingPlan any  `json:"phasedTrainingPlan"`
	ShareableEvent     bool `json:"shareableEvent"`
	PrimaryEvent       bool `json:"primaryEvent"`
	Subscribed         bool `json:"subscribed"`
	AutoCalcCalories   any  `json:"autoCalcCalories"`
	DecoDive           any  `json:"decoDive"`
}

func (c *CalendarService) GetMonth(year int, month time.Month) (*Calendar, error) {
	var cal Calendar
	p := fmt.Sprintf("/calendar-service/year/%d/month/%d", year, int(month))
	return &cal, c.c.apiGet(&cal, p, nil)
}

func (c *CalendarService) GetMonthByDate(date time.Time) (*Calendar, error) {
	return c.GetMonth(date.Year(), date.Month())
}

func (c *CalendarService) GetWeek(year int, month time.Month, startDay int) (*Calendar, error) {
	var cal Calendar
	p := fmt.Sprintf("/calendar-service/year/%d/month/%d/day/%d/start/%d", year, int(month), startDay+7, startDay)
	return &cal, c.c.apiGet(&cal, p, nil)
}

func (c *CalendarService) GetWeekByDate(date time.Time) (*Calendar, error) {
	return c.GetWeek(date.Year(), date.Month(), date.Day()-1)
}

type YearCalendar struct {
	StartDayOfJanuary int  `json:"startDayOfJanuary"`
	LeapYear          bool `json:"leapYear"`
	YearItems         []struct {
		Date    string `json:"date"`
		Display int    `json:"display"`
	} `json:"yearItems"`
	YearSummaries []struct {
		ActivityTypeID     int `json:"activityTypeId"`
		NumberOfActivities int `json:"numberOfActivities"`
		TotalDistance      int `json:"totalDistance"`
		TotalDuration      int `json:"totalDuration"`
		TotalCalories      int `json:"totalCalories"`
	} `json:"yearSummaries"`
}

func (c *CalendarService) GetYear(year int) (*YearCalendar, error) {
	var y YearCalendar
	return &y, c.c.apiGet(&y, fmt.Sprintf("/calendar-service/year/%d", year), nil)
}

type UpcomingEvent struct {
	ID               int    `json:"id"`
	GroupID          any    `json:"groupId"`
	EventName        string `json:"eventName"`
	Date             string `json:"date"`
	URL              string `json:"url"`
	RegistrationURL  string `json:"registrationUrl"`
	CourseID         any    `json:"courseId"`
	CompletionTarget struct {
		Value    float64 `json:"value"`
		Unit     string  `json:"unit"`
		UnitType string  `json:"unitType"`
	} `json:"completionTarget"`
	EventTimeLocal struct {
		StartTimeHhMm string `json:"startTimeHhMm"`
		TimeZoneID    string `json:"timeZoneId"`
	} `json:"eventTimeLocal"`
	Note               any    `json:"note"`
	WorkoutID          any    `json:"workoutId"`
	EventImageUUID     any    `json:"eventImageUUID"`
	Location           string `json:"location"`
	LocationStartPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"locationStartPoint"`
	EventType    string `json:"eventType"`
	EventPrivacy struct {
		Label          string `json:"label"`
		IsShareable    bool   `json:"isShareable"`
		IsDiscoverable bool   `json:"isDiscoverable"`
	} `json:"eventPrivacy"`
	ShareableEventUUID string `json:"shareableEventUuid"`
	EventCustomization struct {
		CustomGoal struct {
			Value    float64 `json:"value"`
			Unit     string  `json:"unit"`
			UnitType string  `json:"unitType"`
		} `json:"customGoal"`
		IsPrimaryEvent           bool `json:"isPrimaryEvent"`
		AssociatedWithActivityID any  `json:"associatedWithActivityId"`
		IsTrainingEvent          bool `json:"isTrainingEvent"`
		IsGoalMet                any  `json:"isGoalMet"`
		TrainingPlanID           any  `json:"trainingPlanId"`
		TrainingPlanType         any  `json:"trainingPlanType"`
	} `json:"eventCustomization"`
	Provider       string `json:"provider"`
	EventRef       string `json:"eventRef"`
	Statuses       any    `json:"statuses"`
	Race           bool   `json:"race"`
	Subscribed     bool   `json:"subscribed"`
	EventOrganizer bool   `json:"eventOrganizer"`
}

func (c *CalendarService) Upcoming(days, limit int) (res []UpcomingEvent, e error) {
	return res, c.c.apiGet(&res, "/calendar-service/events/upcoming", url.Values{
		"numDaysForward": []string{strconv.FormatInt(int64(days), 10)},
		"limit":          []string{strconv.FormatInt(int64(limit), 10)},
	})
}

// GET https://connect.garmin.com/calendar-service/event/primary

// GET https://connect.garmin.com/calendar-service/event/e108b689-6e93-47d3-b4c6-5686fa68b6fb/shareable
type SharableEvent struct {
	ID               int    `json:"id"`
	GroupID          any    `json:"groupId"`
	EventName        string `json:"eventName"`
	Date             string `json:"date"`
	URL              string `json:"url"`
	RegistrationURL  string `json:"registrationUrl"`
	CourseID         any    `json:"courseId"`
	CompletionTarget struct {
		Value    float64 `json:"value"`
		Unit     string  `json:"unit"`
		UnitType string  `json:"unitType"`
	} `json:"completionTarget"`
	EventTimeLocal struct {
		StartTimeHhMm string `json:"startTimeHhMm"`
		TimeZoneID    string `json:"timeZoneId"`
	} `json:"eventTimeLocal"`
	Note               any    `json:"note"`
	WorkoutID          any    `json:"workoutId"`
	EventImageUUID     any    `json:"eventImageUUID"`
	Location           string `json:"location"`
	LocationStartPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"locationStartPoint"`
	EventType    string `json:"eventType"`
	EventPrivacy struct {
		Label          string `json:"label"`
		IsShareable    bool   `json:"isShareable"`
		IsDiscoverable bool   `json:"isDiscoverable"`
	} `json:"eventPrivacy"`
	ShareableEventUUID string `json:"shareableEventUuid"`
	EventCustomization struct {
		CustomGoal struct {
			Value    float64 `json:"value"`
			Unit     string  `json:"unit"`
			UnitType string  `json:"unitType"`
		} `json:"customGoal"`
		IsPrimaryEvent           bool `json:"isPrimaryEvent"`
		AssociatedWithActivityID any  `json:"associatedWithActivityId"`
		IsTrainingEvent          bool `json:"isTrainingEvent"`
		IsGoalMet                any  `json:"isGoalMet"`
		TrainingPlanID           any  `json:"trainingPlanId"`
		TrainingPlanType         any  `json:"trainingPlanType"`
	} `json:"eventCustomization"`
	Provider       string `json:"provider"`
	EventRef       string `json:"eventRef"`
	Statuses       any    `json:"statuses"`
	CourseName     any    `json:"courseName"`
	WorkoutName    any    `json:"workoutName"`
	EventImageURLs any    `json:"eventImageURLs"`
	Race           bool   `json:"race"`
	Subscribed     bool   `json:"subscribed"`
	EventOrganizer bool   `json:"eventOrganizer"`
}

// PUT https://connect.garmin.com/calendar-service/event/e108b689-6e93-47d3-b4c6-5686fa68b6fb/customization
//
// {"customGoal":{"value":"8999.69999999","unit":"second","unitType":"time"},"isPrimaryEvent":true,"associatedWithActivityId":null,"isTrainingEvent":true,"isGoalMet":null,"trainingPlanId":null,"trainingPlanType":null}

type RaceEventProvider struct {
	EventsProviderKey  string `json:"eventsProviderKey"`
	EventsProviderName string `json:"eventsProviderName"`
	VerifiedStatus     string `json:"verifiedStatus"`
	LogoURL            any    `json:"logoURL"`
	CanMakeOfficial    bool   `json:"canMakeOfficial"`
}

func (c *CalendarService) RaceEventProviders() (res []RaceEventProvider, e error) {
	return res, c.c.apiGet(&res, "/calendar-service/race-events/providers", nil)
}

// https://connect.garmin.com/race-search/events?searchPhrase=&poiLat=37.76893&poiLon=-122.26193&withinMeters=80467&fromDate=2024-08-23&toDate=2025-08-23&includeInPerson=true&includeVirtual=false&verifiedStatuses=OFFICIAL%2CVERIFIED&limit=200
// https://connect.garmin.com/race-search/events?searchPhrase=&poiLat=37.76893&poiLon=-122.26193&withinMeters=80467&fromDate=2024-08-23&toDate=2025-08-23&includeInPerson=true&includeVirtual=false&verifiedStatuses=VERIFIED&limit=200
// https://connect.garmin.com/race-search/events?eventType=trail_running&searchPhrase=&poiLat=37.76893&poiLon=-122.26193&withinMeters=160935&fromDate=2024-08-23&toDate=2025-08-23&includeInPerson=true&includeVirtual=false&verifiedStatuses=NONE&limit=200
// https://connect.garmin.com/race-search/events?eventType=trail_running&searchPhrase=Run&poiLat=37.76707&poiLon=-122.24584&withinMeters=160935&fromDate=2024-08-23&toDate=2025-08-23&includeInPerson=true&includeVirtual=false&verifiedStatuses=NONE&limit=200

type RaceSearchRequest struct {
	Latitude         *float64
	Longitude        *float64
	WithinMeters     *int
	FromDate         *time.Time
	ToDate           *time.Time
	IncludeInPerson  *bool
	IncludeVirtual   *bool
	VerifiedStatuses *string
	Limit            *int
	SearchPhrase     *string
}

type RaceSearchResult struct {
	Provider          string `json:"provider"`
	EventRef          string `json:"eventRef"`
	EventName         string `json:"eventName"`
	EventDate         string `json:"eventDate"`
	EventStartTime    any    `json:"eventStartTime"`
	EventURL          string `json:"eventUrl"`
	RegistrationURL   any    `json:"registrationUrl"`
	CompletionTargets []struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"completionTargets"`
	LocationStartPoint struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"locationStartPoint"`
	EventType          string `json:"eventType"`
	DistanceToEvent    int    `json:"distanceToEvent"`
	AdministrativeArea struct {
		CountryCode  string `json:"countryCode"`
		CityEn       string `json:"cityEn"`
		StateEn      string `json:"stateEn"`
		CityNative   string `json:"cityNative"`
		StateNative  string `json:"stateNative"`
		NativeLocale string `json:"nativeLocale"`
	} `json:"administrativeArea"`
	HasCourse        bool   `json:"hasCourse"`
	VerifiedStatus   string `json:"verifiedStatus"`
	GarminEventUUID  string `json:"garminEventUuid"`
	Sig              string `json:"sig"`
	DetailsEndpoints []struct {
		View string `json:"view"`
		URL  string `json:"url"`
	} `json:"detailsEndpoints"`
	IsOfficial bool `json:"isOfficial"`
}
