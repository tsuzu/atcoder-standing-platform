package atcoder

import "time"

// JapanZone 東京のタイムゾーン
var JapanZone = time.FixedZone("JST", 9*60*60)

var invalidFinishTime = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)

// StandingsResponse was generated With JSON-to-Go
type StandingsResponse struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Response []struct {
		Rank  int `json:"rank,omitempty"`
		Tasks []struct {
			StandingsTask
			StandingsUserTask
		} `json:"tasks"`
		UserName       string `json:"user_name,omitempty"`
		UserScreenName string `json:"user_screen_name,omitempty"`
		Failure        string `json:"failure,omitempty"`
		Penalty        string `json:"penalty,omitempty"`
		Score          int    `json:"score,omitempty"`
	} `json:"response"`
	Count int `json:"count"`
}

type StandingsTask struct {
	TaskName       string `json:"task_name"`
	TaskScreenName string `json:"task_screen_name"`
}

type StandingsUser struct {
	UserName       string              `json:"user_name,omitempty"`
	UserScreenName string              `json:"user_screen_name,omitempty"`
	Failure        string              `json:"failure,omitempty"`
	Penalty        string              `json:"penalty,omitempty"`
	Score          int                 `json:"score,omitempty"`
	Tasks          []StandingsUserTask `json:"tasks"`
}

type StandingsUserTask struct {
	Extras      bool   `json:"extras"`
	Score       int    `json:"score"`
	Failure     int    `json:"failure"`
	ElapsedTime string `json:"elapsed_time"`
}

type Standings struct {
	Rank int `json:"rank,omitempty"` // What's this?

	Tasks []StandingsTask `json:"tasks"`

	Users []StandingsUser `json:"users"`
}
