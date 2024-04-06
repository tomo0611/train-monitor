package models

type Train struct {
	No           string `json:"no"`
	Pos          string `json:"pos"`
	Direction    int    `json:"direction"`
	Nickname     string `json:"nickname"`
	Type         string `json:"type"`
	DisplayType  string `json:"displayType"`
	Dest         Dest   `json:"dest"`
	Via          string `json:"via"`
	DelaySeconds int    `json:"delaySeconds"`
	TypeChange   string `json:"typeChange"`
	NumberOfCars int    `json:"numberOfCars"`
}
