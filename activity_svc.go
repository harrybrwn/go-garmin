package garmin

import (
	"fmt"
	"net/url"
)

type ActivityService service

type ActivityUUID struct {
	UUID string `json:"uuid"`
}

func (au *ActivityUUID) String() string { return au.UUID }

type Activity struct {
	ID                 int64               `json:"activityId"`
	UUID               ActivityUUID        `json:"activityUUID"`
	Name               string              `json:"activityName"`
	UserProfileID      int                 `json:"userProfileId"`
	IsMultiSportParent bool                `json:"isMultiSportParent"`
	Type               ActivityType        `json:"activityTypeDTO"`
	EventType          EventType           `json:"eventTypeDTO"`
	AccessControlRule  AccessControlRule   `json:"accessControlRuleDTO"`
	TimeZoneUnit       TimeZoneUnit        `json:"timeZoneUnitDTO"`
	MetadataDTO        ActivityMetadataDTO `json:"metadataDTO"`
	SummaryDTO         ActivitySummaryDTO  `json:"summaryDTO"`
	LocationName       string              `json:"locationName"`
	SplitSummaries     []SplitSummary      `json:"splitSummaries"`
}

type AccessControlRule struct {
	TypeID  int    `json:"typeId"`
	TypeKey string `json:"typeKey"`
}

type TimeZoneUnit struct {
	UnitID   int     `json:"unitId"`
	UnitKey  string  `json:"unitKey"`
	Factor   float64 `json:"factor"`
	TimeZone string  `json:"timeZone"`
}

type ActivityMetadataDTO struct {
	IsOriginal                      bool `json:"isOriginal"`
	DeviceApplicationInstallationID int  `json:"deviceApplicationInstallationId"`
	AgentApplicationInstallationID  any  `json:"agentApplicationInstallationId"`
	AgentString                     any  `json:"agentString"`
	FileFormat                      struct {
		FormatID  int    `json:"formatId"`
		FormatKey string `json:"formatKey"`
	} `json:"fileFormat"`
	AssociatedCourseID  any              `json:"associatedCourseId"`
	LastUpdateDate      string           `json:"lastUpdateDate"`
	UploadedDate        string           `json:"uploadedDate"`
	VideoURL            any              `json:"videoUrl"`
	HasPolyline         bool             `json:"hasPolyline"`
	HasChartData        bool             `json:"hasChartData"`
	HasHrTimeInZones    bool             `json:"hasHrTimeInZones"`
	HasPowerTimeInZones bool             `json:"hasPowerTimeInZones"`
	UserInfoDTO         ActivityUserInfo `json:"userInfoDto"`
	ChildIds            []any            `json:"childIds"`
	ChildActivityTypes  []any            `json:"childActivityTypes"`
	Sensors             any              `json:"sensors"`
	ActivityImages      []any            `json:"activityImages"`
	Manufacturer        string           `json:"manufacturer"`
	DiveNumber          any              `json:"diveNumber"`
	LapCount            int              `json:"lapCount"`
	AssociatedWorkoutID any              `json:"associatedWorkoutId"`
	IsAtpActivity       any              `json:"isAtpActivity"`
	DeviceMetaDataDTO   struct {
		DeviceID        string `json:"deviceId"`
		DeviceTypePk    int    `json:"deviceTypePk"`
		DeviceVersionPk int    `json:"deviceVersionPk"`
	} `json:"deviceMetaDataDTO"`
	HasIntensityIntervals      bool `json:"hasIntensityIntervals"`
	HasSplits                  bool `json:"hasSplits"`
	EBikeMaxAssistModes        any  `json:"eBikeMaxAssistModes"`
	EBikeBatteryUsage          any  `json:"eBikeBatteryUsage"`
	EBikeBatteryRemaining      any  `json:"eBikeBatteryRemaining"`
	EBikeAssistModeInfoDTOList any  `json:"eBikeAssistModeInfoDTOList"`
	HasRunPowerWindData        bool `json:"hasRunPowerWindData"`
	CalendarEventInfo          any  `json:"calendarEventInfo"`
	GroupRideUUID              any  `json:"groupRideUUID"`
	AutoCalcCalories           bool `json:"autoCalcCalories"`
	Favorite                   bool `json:"favorite"`
	ManualActivity             bool `json:"manualActivity"`
	RunPowerWindDataEnabled    bool `json:"runPowerWindDataEnabled"`
	Trimmed                    bool `json:"trimmed"`
	GCJ02                      bool `json:"gcj02"`
	PersonalRecord             bool `json:"personalRecord"`
	ElevationCorrected         bool `json:"elevationCorrected"`
}

type ActivityUserInfo struct {
	UserProfileID         int64  `json:"userProfilePk"`
	Displayname           string `json:"displayname"`
	Fullname              string `json:"fullname"`
	ProfileImageURLLarge  string `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string `json:"profileImageUrlSmall"`
	UserPro               bool   `json:"userPro"`
}

type ActivitySummaryDTO struct {
	StartTimeLocal                 string  `json:"startTimeLocal"`
	StartTimeGMT                   string  `json:"startTimeGMT"`
	StartLatitude                  float64 `json:"startLatitude"`
	StartLongitude                 float64 `json:"startLongitude"`
	Distance                       float64 `json:"distance"`
	Duration                       float64 `json:"duration"`
	MovingDuration                 float64 `json:"movingDuration"`
	ElapsedDuration                float64 `json:"elapsedDuration"`
	ElevationGain                  float64 `json:"elevationGain"`
	ElevationLoss                  float64 `json:"elevationLoss"`
	MaxElevation                   float64 `json:"maxElevation"`
	MinElevation                   float64 `json:"minElevation"`
	AverageSpeed                   float64 `json:"averageSpeed"`
	AverageMovingSpeed             float64 `json:"averageMovingSpeed"`
	MaxSpeed                       float64 `json:"maxSpeed"`
	Calories                       float64 `json:"calories"`
	BmrCalories                    float64 `json:"bmrCalories"`
	AverageHR                      float64 `json:"averageHR"`
	MaxHR                          float64 `json:"maxHR"`
	AverageRunCadence              float64 `json:"averageRunCadence"`
	MaxRunCadence                  float64 `json:"maxRunCadence"`
	AveragePower                   float64 `json:"averagePower"`
	MaxPower                       float64 `json:"maxPower"`
	MinPower                       float64 `json:"minPower"`
	NormalizedPower                float64 `json:"normalizedPower"`
	TotalWork                      float64 `json:"totalWork"`
	GroundContactTime              float64 `json:"groundContactTime"`
	StrideLength                   float64 `json:"strideLength"`
	VerticalOscillation            float64 `json:"verticalOscillation"`
	TrainingEffect                 float64 `json:"trainingEffect"`
	AnaerobicTrainingEffect        float64 `json:"anaerobicTrainingEffect"`
	AerobicTrainingEffectMessage   string  `json:"aerobicTrainingEffectMessage"`
	AnaerobicTrainingEffectMessage string  `json:"anaerobicTrainingEffectMessage"`
	VerticalRatio                  float64 `json:"verticalRatio"`
	EndLatitude                    float64 `json:"endLatitude"`
	EndLongitude                   float64 `json:"endLongitude"`
	MaxVerticalSpeed               float64 `json:"maxVerticalSpeed"`
	WaterEstimated                 float64 `json:"waterEstimated"`
	TrainingEffectLabel            string  `json:"trainingEffectLabel"`
	ActivityTrainingLoad           float64 `json:"activityTrainingLoad"`
	MinActivityLapDuration         float64 `json:"minActivityLapDuration"`
	DirectWorkoutFeel              int     `json:"directWorkoutFeel"`
	DirectWorkoutRpe               int     `json:"directWorkoutRpe"`
	ModerateIntensityMinutes       int     `json:"moderateIntensityMinutes"`
	VigorousIntensityMinutes       int     `json:"vigorousIntensityMinutes"`
	Steps                          int     `json:"steps"`
	RecoveryHeartRate              int     `json:"recoveryHeartRate"`
	AvgGradeAdjustedSpeed          float64 `json:"avgGradeAdjustedSpeed"`
	DifferenceBodyBattery          int     `json:"differenceBodyBattery"`
}

func (as *ActivityService) Get(id int64) (*Activity, error) {
	var a Activity
	p := fmt.Sprintf("/activity-service/activity/%d", id)
	return &a, as.c.apiGet(&a, p, nil)
}

func (as *ActivityService) Activities(req *ActivitySearch) ([]ListedActivity, error) {
	return (*ActivityListService)(as).Activities(req)
}

type ActivityDetails struct {
	ActivityID        int64 `json:"activityId"`
	MeasurementCount  int   `json:"measurementCount"`
	MetricsCount      int   `json:"metricsCount"`
	MetricDescriptors []struct {
		MetricsIndex int    `json:"metricsIndex"`
		Key          string `json:"key"`
		Unit         struct {
			ID     int     `json:"id"`
			Key    string  `json:"key"`
			Factor float64 `json:"factor"`
		} `json:"unit"`
	} `json:"metricDescriptors"`
	ActivityDetailMetrics []struct {
		Metrics []any `json:"metrics"`
	} `json:"activityDetailMetrics"`
	GeoPolylineDTO struct {
		StartPoint struct {
			Lat                       float64 `json:"lat"`
			Lon                       float64 `json:"lon"`
			Altitude                  any     `json:"altitude"`
			Time                      int64   `json:"time"`
			TimerStart                bool    `json:"timerStart"`
			TimerStop                 bool    `json:"timerStop"`
			DistanceFromPreviousPoint any     `json:"distanceFromPreviousPoint"`
			DistanceInMeters          any     `json:"distanceInMeters"`
			Speed                     float64 `json:"speed"`
			CumulativeAscent          any     `json:"cumulativeAscent"`
			CumulativeDescent         any     `json:"cumulativeDescent"`
			ExtendedCoordinate        bool    `json:"extendedCoordinate"`
			Valid                     bool    `json:"valid"`
		} `json:"startPoint"`
		EndPoint struct {
			Lat                       float64 `json:"lat"`
			Lon                       float64 `json:"lon"`
			Altitude                  any     `json:"altitude"`
			Time                      int64   `json:"time"`
			TimerStart                bool    `json:"timerStart"`
			TimerStop                 bool    `json:"timerStop"`
			DistanceFromPreviousPoint any     `json:"distanceFromPreviousPoint"`
			DistanceInMeters          any     `json:"distanceInMeters"`
			Speed                     float64 `json:"speed"`
			CumulativeAscent          any     `json:"cumulativeAscent"`
			CumulativeDescent         any     `json:"cumulativeDescent"`
			ExtendedCoordinate        bool    `json:"extendedCoordinate"`
			Valid                     bool    `json:"valid"`
		} `json:"endPoint"`
		MinLat   float64 `json:"minLat"`
		MaxLat   float64 `json:"maxLat"`
		MinLon   float64 `json:"minLon"`
		MaxLon   float64 `json:"maxLon"`
		Polyline []struct {
			Lat                       float64 `json:"lat"`
			Lon                       float64 `json:"lon"`
			Altitude                  any     `json:"altitude"`
			Time                      int64   `json:"time"`
			TimerStart                bool    `json:"timerStart"`
			TimerStop                 bool    `json:"timerStop"`
			DistanceFromPreviousPoint any     `json:"distanceFromPreviousPoint"`
			DistanceInMeters          any     `json:"distanceInMeters"`
			Speed                     float64 `json:"speed"`
			CumulativeAscent          any     `json:"cumulativeAscent"`
			CumulativeDescent         any     `json:"cumulativeDescent"`
			ExtendedCoordinate        bool    `json:"extendedCoordinate"`
			Valid                     bool    `json:"valid"`
		} `json:"polyline"`
	} `json:"geoPolylineDTO"`
	HeartRateDTOs    any  `json:"heartRateDTOs"`
	PendingData      any  `json:"pendingData"`
	DetailsAvailable bool `json:"detailsAvailable"`
}

// Details gets an activities details given the activity ID.
func (as *ActivityService) Details(id int64) (*ActivityDetails, error) {
	var ad ActivityDetails
	p := fmt.Sprintf("/activity-service/activity/%d/details", id)
	return &ad, as.c.apiGet(&ad, p, url.Values{
		"maxChartSize":    []string{"250"},
		"maxPolylineSize": []string{"2000"},
		"maxHeatMapSize":  []string{"2000"},
	})
}

type ActivityTypedSplits struct {
	ActivityID   int64                `json:"activityId"`
	ActivityUUID ActivityUUID         `json:"activityUUID"`
	Splits       []ActivityTypedSplit `json:"splits"`
}

type ActivityTypedSplit struct {
	StartTimeLocal                  string  `json:"startTimeLocal"`
	StartTimeGMT                    string  `json:"startTimeGMT"`
	StartLatitude                   float64 `json:"startLatitude"`
	StartLongitude                  float64 `json:"startLongitude"`
	Distance                        float64 `json:"distance"`
	Duration                        float64 `json:"duration"`
	MovingDuration                  float64 `json:"movingDuration"`
	ElapsedDuration                 float64 `json:"elapsedDuration"`
	ElevationGain                   float64 `json:"elevationGain"`
	ElevationLoss                   float64 `json:"elevationLoss"`
	AverageSpeed                    float64 `json:"averageSpeed"`
	AverageMovingSpeed              float64 `json:"averageMovingSpeed,omitempty"`
	MaxSpeed                        float64 `json:"maxSpeed"`
	Calories                        float64 `json:"calories"`
	BmrCalories                     float64 `json:"bmrCalories"`
	AverageHR                       float64 `json:"averageHR"`
	MaxHR                           float64 `json:"maxHR"`
	AverageRunCadence               float64 `json:"averageRunCadence"`
	MaxRunCadence                   float64 `json:"maxRunCadence"`
	AveragePower                    float64 `json:"averagePower"`
	MaxPower                        float64 `json:"maxPower"`
	NormalizedPower                 float64 `json:"normalizedPower,omitempty"`
	GroundContactTime               float64 `json:"groundContactTime,omitempty"`
	StrideLength                    float64 `json:"strideLength,omitempty"`
	VerticalOscillation             float64 `json:"verticalOscillation,omitempty"`
	VerticalRatio                   float64 `json:"verticalRatio,omitempty"`
	TotalExerciseReps               int     `json:"totalExerciseReps"`
	EndLatitude                     float64 `json:"endLatitude"`
	EndLongitude                    float64 `json:"endLongitude"`
	AvgVerticalSpeed                float64 `json:"avgVerticalSpeed"`
	AvgGradeAdjustedSpeed           float64 `json:"avgGradeAdjustedSpeed"`
	AvgElapsedDurationVerticalSpeed float64 `json:"avgElapsedDurationVerticalSpeed"`
	Type                            string  `json:"type"`
	MessageIndex                    int     `json:"messageIndex"`
	EndTimeGMT                      string  `json:"endTimeGMT"`
	StartElevation                  float64 `json:"startElevation"`
}

func (as *ActivityService) TypedSplits(id int64) (*ActivityTypedSplits, error) {
	// GET /activity-service/activity/<id>/typedsplits
	var ats ActivityTypedSplits
	p := fmt.Sprintf("/activity-service/activity/%d/typedsplits", id)
	return &ats, as.c.apiGet(&ats, p, nil)
}

type Splits struct {
	ActivityID int64           `json:"activityId"`
	LapDTOs    []LapDTO        `json:"lapDTOs"`
	EventDTOs  []SplitEventDTO `json:"eventDTOs"`
}

type LapDTO struct {
	StartTimeGMT          string  `json:"startTimeGMT"`
	StartLatitude         float64 `json:"startLatitude"`
	StartLongitude        float64 `json:"startLongitude"`
	Distance              float64 `json:"distance"`
	Duration              float64 `json:"duration"`
	MovingDuration        float64 `json:"movingDuration"`
	ElapsedDuration       float64 `json:"elapsedDuration"`
	ElevationGain         float64 `json:"elevationGain"`
	ElevationLoss         float64 `json:"elevationLoss"`
	MaxElevation          float64 `json:"maxElevation"`
	MinElevation          float64 `json:"minElevation"`
	AverageSpeed          float64 `json:"averageSpeed"`
	AverageMovingSpeed    float64 `json:"averageMovingSpeed"`
	MaxSpeed              float64 `json:"maxSpeed"`
	Calories              float64 `json:"calories"`
	BmrCalories           float64 `json:"bmrCalories"`
	AverageHR             float64 `json:"averageHR"`
	MaxHR                 float64 `json:"maxHR"`
	AverageRunCadence     float64 `json:"averageRunCadence"`
	MaxRunCadence         float64 `json:"maxRunCadence"`
	AveragePower          float64 `json:"averagePower"`
	MaxPower              float64 `json:"maxPower"`
	MinPower              float64 `json:"minPower"`
	NormalizedPower       float64 `json:"normalizedPower"`
	TotalWork             float64 `json:"totalWork"`
	GroundContactTime     float64 `json:"groundContactTime"`
	StrideLength          float64 `json:"strideLength"`
	VerticalOscillation   float64 `json:"verticalOscillation"`
	VerticalRatio         float64 `json:"verticalRatio"`
	EndLatitude           float64 `json:"endLatitude"`
	EndLongitude          float64 `json:"endLongitude"`
	MaxVerticalSpeed      float64 `json:"maxVerticalSpeed"`
	AvgGradeAdjustedSpeed float64 `json:"avgGradeAdjustedSpeed"`
	LapIndex              int     `json:"lapIndex"`
	LengthDTOs            []any   `json:"lengthDTOs"`
	ConnectIQMeasurement  []any   `json:"connectIQMeasurement"`
	IntensityType         string  `json:"intensityType"`
	MessageIndex          int     `json:"messageIndex"`
}

type SplitEventDTO struct {
	StartTimeGMT            string `json:"startTimeGMT"`
	StartTimeGMTDoubleValue int64  `json:"startTimeGMTDoubleValue"`
	SectionTypeDTO          struct {
		ID             int    `json:"id"`
		Key            string `json:"key"`
		SectionTypeKey string `json:"sectionTypeKey"`
	} `json:"sectionTypeDTO"`
}

func (as *ActivityService) Splits(id int64) (*Splits, error) {
	var s Splits
	return &s, as.c.apiGet(&s, fmt.Sprintf("/activity-service/activity/%d/splits", id), nil)
}

type SplitSummaries struct {
	ActivityID     int64          `json:"activityId"`
	ActivityUUID   ActivityUUID   `json:"activityUUID"`
	SplitSummaries []SplitSummary `json:"splitSummaries"`
}

type SplitSummary struct {
	Distance              float64 `json:"distance"`
	Duration              float64 `json:"duration"`
	MovingDuration        float64 `json:"movingDuration"`
	ElevationGain         float64 `json:"elevationGain"`
	ElevationLoss         float64 `json:"elevationLoss"`
	AverageSpeed          float64 `json:"averageSpeed"`
	AverageMovingSpeed    float64 `json:"averageMovingSpeed,omitempty"`
	MaxSpeed              float64 `json:"maxSpeed"`
	Calories              float64 `json:"calories"`
	BMRCalories           float64 `json:"bmrCalories"`
	AverageHR             float64 `json:"averageHR"`
	MaxHR                 float64 `json:"maxHR"`
	AverageRunCadence     float64 `json:"averageRunCadence"`
	MaxRunCadence         float64 `json:"maxRunCadence"`
	AveragePower          float64 `json:"averagePower"`
	MaxPower              float64 `json:"maxPower"`
	NormalizedPower       float64 `json:"normalizedPower"`
	GroundContactTime     float64 `json:"groundContactTime,omitempty"`
	StrideLength          float64 `json:"strideLength"`
	VerticalOscillation   float64 `json:"verticalOscillation,omitempty"`
	VerticalRatio         float64 `json:"verticalRatio,omitempty"`
	TotalExerciseReps     int     `json:"totalExerciseReps"`
	AvgVerticalSpeed      float64 `json:"avgVerticalSpeed"`
	AvgGradeAdjustedSpeed float64 `json:"avgGradeAdjustedSpeed"`
	SplitType             string  `json:"splitType"`
	NoOfSplits            int     `json:"noOfSplits"`
	MaxElevationGain      float64 `json:"maxElevationGain"`
	AverageElevationGain  float64 `json:"averageElevationGain"`
	MaxDistance           int     `json:"maxDistance"`
}

func (as *ActivityService) SplitSummaries(id int64) (*SplitSummaries, error) {
	var s SplitSummaries
	p := fmt.Sprintf("/activity-service/activity/%d/split_summaries", id)
	return &s, as.c.apiGet(&s, p, nil)
}

type TimeInZone struct {
	ZoneNumber      int     `json:"zoneNumber"`
	SecsInZone      float64 `json:"secsInZone"`
	ZoneLowBoundary int     `json:"zoneLowBoundary"`
}

func (as *ActivityService) HeartRateTimeInZones(id int64) (res []TimeInZone, err error) {
	p := fmt.Sprintf("/activity-service/activity/%d/hrTimeInZones", id)
	return res, as.c.apiGet(&res, p, nil)
}

func (as *ActivityService) PowerTimeInZones(id int64) (res []TimeInZone, e error) {
	p := fmt.Sprintf("/activity-service/activity/%d/powerTimeInZones", id)
	return res, as.c.apiGet(&res, p, nil)
}

type ActivityWeather struct {
	IssueDate                 string  `json:"issueDate"`
	Temp                      int     `json:"temp"`
	ApparentTemp              int     `json:"apparentTemp"`
	DewPoint                  int     `json:"dewPoint"`
	RelativeHumidity          int     `json:"relativeHumidity"`
	WindDirection             int     `json:"windDirection"`
	WindDirectionCompassPoint string  `json:"windDirectionCompassPoint"`
	WindSpeed                 int     `json:"windSpeed"`
	WindGust                  any     `json:"windGust"`
	Latitude                  float64 `json:"latitude"`
	Longitude                 float64 `json:"longitude"`
	WeatherStationDTO         struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Timezone any    `json:"timezone"`
	} `json:"weatherStationDTO"`
	WeatherTypeDTO struct {
		WeatherTypePk any    `json:"weatherTypePk"`
		Desc          string `json:"desc"`
		Image         any    `json:"image"`
	} `json:"weatherTypeDTO"`
}

func (as *ActivityService) Weather(id int64) (*ActivityWeather, error) {
	var aw ActivityWeather
	p := fmt.Sprintf("/activity-service/activity/%d/weather", id)
	return &aw, as.c.apiGet(&aw, p, nil)
}

type ActivityType struct {
	TypeID       int    `json:"typeId"`
	TypeKey      string `json:"typeKey"`
	ParentTypeID int    `json:"parentTypeId"`
	IsHidden     bool   `json:"isHidden"`
	Restricted   bool   `json:"restricted"`
	Trimmable    bool   `json:"trimmable"`
}

func (as *ActivityService) Types() (at []ActivityType, err error) {
	return at, as.c.apiGet(&at, "/activity-service/activity/activityTypes", nil)
}

type EventType struct {
	TypeID    int    `json:"typeId"`
	TypeKey   string `json:"typeKey"`
	SortOrder int    `json:"sortOrder"`
}

func (as *ActivityService) EventTypes() (et []EventType, err error) {
	return et, as.c.apiGet(&et, "/activity-service/activity/eventTypes", nil)
}

func (as *ActivityService) Workouts(id int64) {
	// GET https://connect.garmin.com/activity-service/activity/<id>/workouts
	panic("not implemented")
}
