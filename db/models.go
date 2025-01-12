package db

type Period struct {
    Id  int `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}
type Day struct {
    Day_Total string `json:"total"`
	Periods   []Period
}

type Timers struct {
	Days  map[string]Day `json:"days"`
}

type Data struct {
	Sessions map[string]Timers `json:"sessions"`
}