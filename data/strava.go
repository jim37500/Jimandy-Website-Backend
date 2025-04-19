package data

type Token struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	Athlete      Athlete `json:"athlete"`
}

type Athlete struct {
	AthleteID uint64 `json:"id"`
}

type Activity struct {
	ActivityID uint64 `json:"id"`
}

type RouteMap struct {
	Polyline string `json:"summary_polyline"`
}

type Lap struct {
	LapID              uint64   `json:"id"`
	Activity           Activity `json:"activity"`
	Athlete            Athlete  `json:"athlete"`
	ElapsedTime        int      `json:"elapsed_time"`
	MovingTime         int      `json:"moving_time"`
	Distance           float32  `json:"distance"`
	AverageSpeed       float32  `json:"average_speed"`
	MaxSpeed           float32  `json:"max_speed"`
	AverageCadence     float32  `json:"average_cadence"`
	AverageHeartrate   float32  `json:"average_heartrate"`
	MaxHeartrate       float32  `json:"max_heartrate"`
	AverageWatts       float32  `json:"average_watts"`
	TotalElevationGain float32  `json:"total_elevation_gain"`
}

type Activities struct {
	ActivityID         uint64   `json:"id"`
	Name               string   `json:"name"`
	Athlete            Athlete  `json:"athlete"`
	SportType          string   `json:"sport_type"`
	Date               string   `json:"start_date"`
	ElapsedTime        int      `json:"elapsed_time"`
	MovingTime         int      `json:"moving_time"`
	Distance           float32  `json:"distance"`
	TotalElevationGain float32  `json:"total_elevation_gain"`
	AverageSpeed       float32  `json:"average_speed"`
	MaxSpeed           float32  `json:"max_speed"`
	AverageCadence     float32  `json:"average_cadence"`
	AverageHeartrate   float32  `json:"average_heartrate"`
	MaxHeartrate       float32  `json:"max_heartrate"`
	AverageWatts       float32  `json:"average_watts"`
	MaxWatts           float32  `json:"max_watts"`
	AverageTemperature int      `json:"average_temp"`
	RouteMap           RouteMap `json:"map"`
}
