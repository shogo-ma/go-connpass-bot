package models

import "github.com/franela/goreq"

type Connpass struct {
	Events           []Event `json:"events"`
	ResultsAvailable int64   `json:"results_available"`
	ResultsReturned  int64   `json:"results_returned"`
	ResultsStart     int64   `json:"results_start"`
}

type Event struct {
	Accepted    int64  `json:"accepted"`
	Address     string `json:"address"`
	Catch       string `json:"catch"`
	Description string `json:"description"`
	EndedAt     string `json:"ended_at"`
	EventID     int64  `json:"event_id"`
	EventType   string `json:"event_type"`
	EventURL    string `json:"event_url"`
	HashTag     string `json:"hash_tag"`
	//Lat              float64 `json:"lat"`
	//Lon              float64 `json:"lon"`
	Limit            int64  `json:"limit"`
	OwnerDisplayName string `json:"owner_display_name"`
	OwnerID          int64  `json:"owner_id"`
	OwnerNickname    string `json:"owner_nickname"`
	Place            string `json:"place"`
	Series           struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"series"`
	StartedAt string `json:"started_at"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
	Waiting   int64  `json:"waiting"`
}

type Params struct {
	EventId int
	Keyword string
	Ym      int64
	Ymd     int64
	Count   int
	Order   int
}

const API_ENDPOINT = "https://connpass.com/api/v1/event/"

func Request(p *Params) (Connpass, error) {
	var cps Connpass
	res, err := goreq.Request{
		Uri:         API_ENDPOINT,
		QueryString: p,
	}.Do()

	if err != nil {
		return cps, err
	}

	res.Body.FromJsonTo(&cps)
	return cps, err
}
