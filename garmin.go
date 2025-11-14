package garmin

import (
	"strconv"
	"time"
)

// API holds all the services.
type API struct {
	Activity       *ActivityService
	ActivityList   *ActivityListService
	Course         *CourseService
	Device         *DeviceService
	FitnessAge     *FitnessAgeService
	FitnessStats   *FitnessStatsService
	PersonalRecord *PersonalRecordService
	Sleep          *SleepService
	UserFocus      *UserFocusService
	UserProfile    *UserProfileService
	UserSummary    *UserSummaryService
	Weight         *WeightService
	Wellness       *WellnessService
}

// NewAPI creates a new API struct.
func NewAPI(client *Client) *API {
	s := service{c: client}
	return &API{
		Activity:       (*ActivityService)(&s),
		ActivityList:   (*ActivityListService)(&s),
		Course:         (*CourseService)(&s),
		Device:         (*DeviceService)(&s),
		FitnessAge:     (*FitnessAgeService)(&s),
		FitnessStats:   (*FitnessStatsService)(&s),
		PersonalRecord: (*PersonalRecordService)(&s),
		Sleep:          (*SleepService)(&s),
		UserFocus:      (*UserFocusService)(&s),
		UserProfile:    (*UserProfileService)(&s),
		UserSummary:    (*UserSummaryService)(&s),
		Weight:         (*WeightService)(&s),
		Wellness:       (*WellnessService)(&s),
	}
}

type service struct {
	c *Client
}

var DefaultClock clock

type Clock interface {
	Now() time.Time
}

type clock struct{}

func (*clock) Now() time.Time { return time.Now() }

type UnixTS time.Time

func (uts *UnixTS) UnmarshalJSON(b []byte) error { return uts.UnmarshalText(b) }
func (uts *UnixTS) UnmarshalText(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*uts = UnixTS(time.UnixMilli(i))
	return nil
}

func (uts UnixTS) Add(d time.Duration) UnixTS      { return UnixTS(time.Time(uts).Add(d)) }
func (uts UnixTS) Sub(t UnixTS) time.Duration      { return time.Time(uts).Sub(time.Time(t)) }
func (uts UnixTS) AddDate(y, m, d int) UnixTS      { return UnixTS(time.Time(uts).AddDate(y, m, d)) }
func (uts UnixTS) UTC() UnixTS                     { return UnixTS(time.Time(uts).UTC()) }
func (uts UnixTS) Local() UnixTS                   { return UnixTS(time.Time(uts).Local()) }
func (uts UnixTS) IsDST() bool                     { return time.Time(uts).IsDST() }
func (uts UnixTS) Zone() (string, int)             { return time.Time(uts).Zone() }
func (uts UnixTS) Truncate(d time.Duration) UnixTS { return UnixTS(time.Time(uts).Truncate(d)) }
func (uts UnixTS) Round(d time.Duration) UnixTS    { return UnixTS(time.Time(uts).Round(d)) }
func (uts UnixTS) In(loc *time.Location) UnixTS    { return UnixTS(time.Time(uts).In(loc)) }
func (uts UnixTS) Location() *time.Location        { return time.Time(uts).Location() }
func (uts UnixTS) Unix() int64                     { return time.Time(uts).Unix() }
func (uts UnixTS) UnixMilli() int64                { return time.Time(uts).UnixMilli() }
func (uts UnixTS) UnixMicro() int64                { return time.Time(uts).UnixMicro() }
func (uts UnixTS) UnixNano() int64                 { return time.Time(uts).UnixNano() }
func (uts UnixTS) After(ts UnixTS) bool            { return time.Time(uts).After(time.Time(ts)) }
func (uts UnixTS) Before(ts UnixTS) bool           { return time.Time(uts).Before(time.Time(ts)) }
func (uts UnixTS) Compare(ts UnixTS) int           { return time.Time(uts).Compare(time.Time(ts)) }
func (uts UnixTS) Equal(ts UnixTS) bool            { return time.Time(uts).Equal(time.Time(ts)) }
func (uts UnixTS) IsZero() bool                    { return time.Time(uts).IsZero() }
func (uts UnixTS) Date() (int, time.Month, int)    { return time.Time(uts).Date() }
func (uts UnixTS) Year() int                       { return time.Time(uts).Year() }
func (uts UnixTS) Month() time.Month               { return time.Time(uts).Month() }
func (uts UnixTS) Day() int                        { return time.Time(uts).Day() }
func (uts UnixTS) Weekday() time.Weekday           { return time.Time(uts).Weekday() }
func (uts UnixTS) ISOWeek() (int, int)             { return time.Time(uts).ISOWeek() }
func (uts UnixTS) Clock() (int, int, int)          { return time.Time(uts).Clock() }
func (uts UnixTS) Hour() int                       { return time.Time(uts).Hour() }
func (uts UnixTS) Minute() int                     { return time.Time(uts).Minute() }
func (uts UnixTS) Second() int                     { return time.Time(uts).Second() }
func (uts UnixTS) Nanosecond() int                 { return time.Time(uts).Nanosecond() }
func (uts UnixTS) YearDay() int                    { return time.Time(uts).YearDay() }
func (uts UnixTS) Format(format string) string     { return time.Time(uts).Format(format) }
