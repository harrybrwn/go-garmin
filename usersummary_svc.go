package garmin

import (
	"fmt"
	"net/url"
	"time"
)

type Stat[Stat any] struct {
	CalendarDate string `json:"calendarDate"`
	Values       Stat   `json:"values"`
}

type UserSummaryService service

type StressStat struct {
	HighStressDuration   any `json:"highStressDuration"`
	LowStressDuration    int `json:"lowStressDuration"`
	OverallStressLevel   int `json:"overallStressLevel"`
	RestStressDuration   int `json:"restStressDuration"`
	MediumStressDuration any `json:"mediumStressDuration"`
}

func datepath(base string, a, b time.Time) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		base,
		a.Format(time.DateOnly),
		b.Format(time.DateOnly),
	)
}

func (uss *UserSummaryService) DailyStress(start, end time.Time) (s []Stat[StressStat], err error) {
	p := datepath("/usersummary-service/stats/stress/daily", start, end)
	return s, uss.c.apiGet(&s, p, nil)
}

type WeeklyStressStat struct {
	CalendarDate string  `json:"calendarDate"`
	Value        float64 `json:"value"`
}

func (uss *UserSummaryService) WeeklyStress(weeks int, end time.Time) (s []WeeklyStressStat, err error) {
	p := fmt.Sprintf("/usersummary-service/stats/stress/weekly/%s/%d", end.Format(time.DateOnly), weeks)
	return s, uss.c.apiGet(&s, p, nil)
}

type HeartRateStat struct {
	WellnessMaxAvgHR int `json:"wellnessMaxAvgHR"`
	WellnessMinAvgHR int `json:"wellnessMinAvgHR"`
	// only available for daily data
	RestingHR int `json:"restingHR"`
	// only available for weekly data
	AvgRestingHR int `json:"avgRestingHR"`
}

func (uss *UserSummaryService) DailyHeartRate(start, end time.Time) (s []Stat[HeartRateStat], err error) {
	p := datepath("/usersummary-service/stats/heartRate/daily", start, end)
	return s, uss.c.apiGet(&s, p, nil)
}

func (uss *UserSummaryService) WeeklyHeartRate(weeks int, end time.Time) (s []Stat[HeartRateStat], err error) {
	p := fmt.Sprintf("/usersummary-service/stats/heartRate/weekly/%s/%d", end.Format(time.DateOnly), weeks)
	return s, uss.c.apiGet(&s, p, nil)
}

type BodyBatteryStat struct {
	LowBodyBattery  int `json:"lowBodyBattery"`
	HighBodyBattery int `json:"highBodyBattery"`
}

func (uss *UserSummaryService) DailyBodyBattery(start, end time.Time) (s []Stat[BodyBatteryStat], err error) {
	p := datepath("/usersummary-service/stats/bodybattery/daily", start, end)
	return s, uss.c.apiGet(&s, p, nil)
}

type DailyStepsStat struct {
	StepGoal      int `json:"stepGoal"`
	TotalSteps    int `json:"totalSteps"`
	TotalDistance int `json:"totalDistance"`
}

type DailySteps struct {
	Values       []Stat[DailyStepsStat] `json:"values"`
	Aggregations StepsAggregations      `json:"aggregations"`
}

type StepsAggregations struct {
	TotalStepsAverage       float64 `json:"totalStepsAverage"`
	TotalStepsWeeklyAverage float64 `json:"totalStepsWeeklyAverage"`
}

func (uss *UserSummaryService) DailySteps(start, end time.Time) (s *DailySteps, err error) {
	// GET https://connect.garmin.com/usersummary-service/stats/daily/2024-08-10/2024-08-16?statsType=STEPS&currentDate=2024-08-16
	now := uss.c.Clock.Now()
	s = new(DailySteps)
	return s, uss.c.apiGet(
		s,
		datepath("/usersummary-service/stats/daily", start, end),
		url.Values{
			"statsType":   []string{"STEPS"},
			"currentDate": []string{now.Format(time.DateOnly)},
		},
	)
}

type MonthlyStepsStat struct {
	WellnessDataDaysCount int `json:"wellnessDataDaysCount"`
	TotalSteps            int `json:"totalSteps"`
	TotalDistance         int `json:"totalDistance"`
	TotalStepsGoal        int `json:"TotalStepsGoal"`
}

func (uss *UserSummaryService) MonthlySteps(months int, end time.Time) (s []Stat[MonthlyStepsStat], err error) {
	p := fmt.Sprintf("/usersummary-service/stats/steps/monthly/%s/%d", end.Format(time.DateOnly), months)
	return s, uss.c.apiGet(&s, p, nil)
}

type WeeklyStepsStat struct {
	TotalSteps            float64 `json:"totalSteps"`
	AverageSteps          float64 `json:"averageSteps"`
	WellnessDataDaysCount int     `json:"wellnessDataDaysCount"`
	AverageDistance       float64 `json:"averageDistance"`
	TotalDistance         float64 `json:"totalDistance"`
}

func (uss *UserSummaryService) WeeklySteps(weeks int, end time.Time) (s []Stat[WeeklyStepsStat], err error) {
	p := fmt.Sprintf("/usersummary-service/stats/steps/weekly/%s/%d", end.Format(time.DateOnly), weeks)
	return s, uss.c.apiGet(&s, p, nil)
}

func (uss *UserSummaryService) MonthlyPushes(months int, end time.Time) ([]any, error) {
	// GET https://connect.garmin.com/usersummary-service/stats/pushes/monthly/2024-08-16/12
	// TODO Find the struct fields returned
	panic("not implemented")
}

func (uss *UserSummaryService) WeeklyPushes(weeks int, end time.Time) ([]any, error) {
	// GET https://connect.garmin.com/usersummary-service/stats/pushes/weekly/2024-08-16/52
	// TODO Find the struct fields returned
	panic("not implemented")
}

type IntensityMinutesStat struct {
	CalendarDate  string `json:"calendarDate"`
	WeeklyGoal    int    `json:"weeklyGoal"`
	ModerateValue int    `json:"moderateValue"`
	VigorousValue int    `json:"vigorousValue"`
}

func (uss *UserSummaryService) DailyIntensityMinutes(start, end time.Time) (s []IntensityMinutesStat, err error) {
	return s, uss.c.apiGet(&s, datepath("/usersummary-service/stats/im/daily", start, end), nil)
}

func (uss *UserSummaryService) WeeklyIntensityMinutes(start, end time.Time) (s []IntensityMinutesStat, err error) {
	return s, uss.c.apiGet(&s, datepath("/usersummary-service/stats/im/weekly", start, end), nil)
}
