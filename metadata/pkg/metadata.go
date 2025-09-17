package model

type Metadata struct {
	CubicleID string `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	Capacity  int    `json:"capacity"`
}
