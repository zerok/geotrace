package server

type Location struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Timestamp string  `json:"timestamp"`
	Altitude  float64 `json:"altitude"`
	DeviceID  string  `json:"device_id"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Request struct {
	Locations []Location `json:"locations"`
}
