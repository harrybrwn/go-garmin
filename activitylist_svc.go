package garmin

import (
	"net/url"
	"strconv"
)

type ActivityListService service

type ListedActivity struct {
	ID                                    int64        `json:"activityId"`
	Name                                  string       `json:"activityName"`
	StartTimeLocal                        string       `json:"startTimeLocal"`
	StartTimeGMT                          string       `json:"startTimeGMT"`
	ActivityType                          ActivityType `json:"activityType"`
	EventType                             EventType    `json:"eventType"`
	Distance                              float64      `json:"distance"`
	Duration                              float64      `json:"duration"`
	ElapsedDuration                       float64      `json:"elapsedDuration"`
	MovingDuration                        float64      `json:"movingDuration"`
	ElevationGain                         float64      `json:"elevationGain"`
	ElevationLoss                         float64      `json:"elevationLoss"`
	AverageSpeed                          float64      `json:"averageSpeed"`
	MaxSpeed                              float64      `json:"maxSpeed"`
	StartLatitude                         float64      `json:"startLatitude"`
	StartLongitude                        float64      `json:"startLongitude"`
	HasPolyline                           bool         `json:"hasPolyline"`
	HasImages                             bool         `json:"hasImages"`
	OwnerID                               int          `json:"ownerId"`
	OwnerDisplayName                      string       `json:"ownerDisplayName"`
	OwnerFullName                         string       `json:"ownerFullName"`
	OwnerProfileImageURLSmall             string       `json:"ownerProfileImageUrlSmall"`
	OwnerProfileImageURLMedium            string       `json:"ownerProfileImageUrlMedium"`
	OwnerProfileImageURLLarge             string       `json:"ownerProfileImageUrlLarge"`
	Calories                              float64      `json:"calories"`
	BmrCalories                           float64      `json:"bmrCalories"`
	AverageHR                             float64      `json:"averageHR"`
	MaxHR                                 float64      `json:"maxHR"`
	AverageRunningCadenceInStepsPerMinute float64      `json:"averageRunningCadenceInStepsPerMinute"`
	MaxRunningCadenceInStepsPerMinute     float64      `json:"maxRunningCadenceInStepsPerMinute"`
	Steps                                 int          `json:"steps"`
	UserRoles                             []string     `json:"userRoles"`
	Privacy                               struct {
		TypeID  int    `json:"typeId"`
		TypeKey string `json:"typeKey"`
	} `json:"privacy"`
	UserPro                 bool    `json:"userPro"`
	CourseID                int     `json:"courseId,omitempty"`
	HasVideo                bool    `json:"hasVideo"`
	TimeZoneID              int     `json:"timeZoneId"`
	BeginTimestamp          int64   `json:"beginTimestamp"`
	SportTypeID             int     `json:"sportTypeId"`
	AvgPower                float64 `json:"avgPower"`
	MaxPower                float64 `json:"maxPower"`
	AerobicTrainingEffect   float64 `json:"aerobicTrainingEffect"`
	AnaerobicTrainingEffect float64 `json:"anaerobicTrainingEffect"`
	NormPower               float64 `json:"normPower"`
	AvgVerticalOscillation  float64 `json:"avgVerticalOscillation"`
	AvgGroundContactTime    float64 `json:"avgGroundContactTime"`
	AvgStrideLength         float64 `json:"avgStrideLength"`
	AvgLeftBalance          float64 `json:"avgLeftBalance"`
	VO2MaxValue             float64 `json:"vO2MaxValue"`
	AvgVerticalRatio        float64 `json:"avgVerticalRatio"`
	DeviceID                int64   `json:"deviceId"`
	MinElevation            float64 `json:"minElevation"`
	MaxElevation            float64 `json:"maxElevation"`
	MaxDoubleCadence        float64 `json:"maxDoubleCadence"`
	SummarizedDiveInfo      struct {
		SummarizedDiveGases []any `json:"summarizedDiveGases"`
	} `json:"summarizedDiveInfo"`
	MaxVerticalSpeed               float64                      `json:"maxVerticalSpeed"`
	Manufacturer                   string                       `json:"manufacturer"`
	LocationName                   string                       `json:"locationName"`
	LapCount                       int                          `json:"lapCount"`
	EndLatitude                    float64                      `json:"endLatitude"`
	EndLongitude                   float64                      `json:"endLongitude"`
	WaterEstimated                 float64                      `json:"waterEstimated"`
	TrainingEffectLabel            string                       `json:"trainingEffectLabel"`
	ActivityTrainingLoad           float64                      `json:"activityTrainingLoad"`
	MinActivityLapDuration         float64                      `json:"minActivityLapDuration"`
	AerobicTrainingEffectMessage   string                       `json:"aerobicTrainingEffectMessage"`
	AnaerobicTrainingEffectMessage string                       `json:"anaerobicTrainingEffectMessage"`
	SplitSummaries                 []ListedActivitySplitSummary `json:"splitSummaries"`
	HasSplits                      bool                         `json:"hasSplits"`
	ModerateIntensityMinutes       int                          `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes       int                          `json:"vigorousIntensityMinutes"`
	AvgGradeAdjustedSpeed          float64                      `json:"avgGradeAdjustedSpeed"`
	DifferenceBodyBattery          int                          `json:"differenceBodyBattery"`
	Pr                             bool                         `json:"pr"`
	AutoCalcCalories               bool                         `json:"autoCalcCalories"`
	ElevationCorrected             bool                         `json:"elevationCorrected"`
	AtpActivity                    bool                         `json:"atpActivity"`
	Favorite                       bool                         `json:"favorite"`
	DecoDive                       bool                         `json:"decoDive"`
	Parent                         bool                         `json:"parent"`
	Purposeful                     bool                         `json:"purposeful"`
	ManualActivity                 bool                         `json:"manualActivity"`
}

type ListedActivitySplitSummary struct {
	NoOfSplits           int     `json:"noOfSplits"`
	TotalAscent          float64 `json:"totalAscent"`
	Duration             float64 `json:"duration"`
	SplitType            string  `json:"splitType"`
	NumClimbSends        int     `json:"numClimbSends"`
	MaxElevationGain     float64 `json:"maxElevationGain"`
	AverageElevationGain float64 `json:"averageElevationGain"`
	MaxDistance          int     `json:"maxDistance"`
	Distance             float64 `json:"distance"`
	AverageSpeed         float64 `json:"averageSpeed"`
	MaxSpeed             float64 `json:"maxSpeed"`
	NumFalls             int     `json:"numFalls"`
	ElevationLoss        float64 `json:"elevationLoss"`
}

type ActivitySearch struct {
	Limit           *int
	Start           *int
	Search          *string
	ExcludeChildren *bool
	ActivityType    *string
	Favorite        *int
}

func (as *ActivitySearch) WithLimit(lim int) *ActivitySearch {
	as.Limit = &lim
	return as
}

func (as *ActivitySearch) WithStart(start int) *ActivitySearch {
	as.Start = &start
	return as
}

func (as *ActivitySearch) WithSearch(search string) *ActivitySearch {
	as.Search = &search
	return as
}

func (as *ActivitySearch) WithActivityType(typ string) *ActivitySearch {
	as.ActivityType = &typ
	return as
}

func (as *ActivitySearch) WithExcludeChildren(v bool) *ActivitySearch {
	as.ExcludeChildren = &v
	return as
}

func (as *ActivitySearch) WithFavorites() *ActivitySearch {
	v := 1
	as.Favorite = &v
	return as
}

func (as *ActivitySearch) params() url.Values {
	v := make(url.Values, 3)
	if as == nil {
		return v
	}
	if as.Limit != nil {
		v.Set("limit", strconv.FormatInt(int64(*as.Limit), 10))
	}
	if as.Start != nil {
		v.Set("start", strconv.FormatInt(int64(*as.Start), 10))
	}
	if as.ActivityType != nil {
		v.Set("activityType", *as.ActivityType)
	}
	if as.Favorite != nil {
		v.Set("favorite", strconv.FormatInt(int64(*as.Favorite), 10))
	}
	if as.ExcludeChildren != nil {
		v.Set("excludeChildren", strconv.FormatBool(*as.ExcludeChildren))
	}
	if as.Search != nil {
		v.Set("search", *as.Search)
	}
	return v
}

func (al *ActivityListService) Activities(req *ActivitySearch) (list []ListedActivity, e error) {
	// GET /activitylist-service/activities/search/activities?limit=20&start=0
	// GET /activitylist-service/activities/search/activities?activityType=running&limit=20&excludeChildren=false&start=0
	// GET /activitylist-service/activities/search/activities?favorite=1&limit=20&start=0
	// GET /activitylist-service/activities/search/activities?search=Trail&limit=20&start=0
	return list, al.c.apiGet(&list, "/activitylist-service/activities/search/activities", req.params())
}

// FirstLast returns the first and last activity IDs.
func (al *ActivityListService) FirstLast() (int64, int64, error) {
	var res struct {
		First int64 `json:"firstActivityId"`
		Last  int64 `json:"lastActivityId"`
	}
	return res.First, res.Last, al.c.apiGet(&res, "/activitylist-service/activities/first-last", nil)
}
