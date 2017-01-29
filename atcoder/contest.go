package atcoder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
)

type Contest struct {
	Title      string
	StartTime  time.Time
	FinishTime time.Time
	URL        string
	TitleInURL string
	client     *http.Client
}

func (contest *Contest) UpdateClient(client *http.Client) {
	contest.client = client
}

func (contest *Contest) GetMyStatus() (string, error) {
	resp, err := contest.client.Get(contest.URL)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "__privilege" {
			return c.Value, nil
		}
	}

	return "", nil
}

func (contest *Contest) IsJoined() (bool, error) {
	if s, err := contest.GetMyStatus(); err != nil {
		return false, err
	} else if s == "contestant" {
		return true, nil
	} else {
		return false, nil
	}
}

func (contest *Contest) Join() error {
	resp, err := contest.client.Get(contest.URL + "/participants/insert")
	defer resp.Body.Close()

	return err
}

func (contest *Contest) Standings() (*Standings, error) {
	resp, err := contest.client.Get(contest.URL + "/standings/json")

	if err != nil {
		return nil, err
	}

	logrus.Println(resp.Header)

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var standings StandingsResponse
	err = json.Unmarshal(b, &standings)

	if err != nil {
		return nil, err
	}

	if standings.Status != 200 {
		return nil, errors.New(standings.Message)
	}

	var res Standings

	if len(standings.Response) == 0 {
		return nil, errors.New("Problem not found")
	}
	res.Rank = standings.Response[0].Rank
	res.Tasks = make([]StandingsTask, 0, len(standings.Response[0].Tasks))

	for i := range standings.Response[0].Tasks {
		res.Tasks = append(res.Tasks, standings.Response[0].Tasks[i].StandingsTask)
	}

	res.Users = make([]StandingsUser, 0, len(standings.Response)-1)

	for i := 1; i < len(standings.Response); i++ {
		var user StandingsUser
		user.Failure = standings.Response[i].Failure
		user.Penalty = standings.Response[i].Penalty
		user.Score = standings.Response[i].Score
		user.UserName = standings.Response[i].UserName
		user.UserScreenName = standings.Response[i].UserScreenName

		user.Tasks = make([]StandingsUserTask, 0, len(standings.Response[i].Tasks))

		for j := range standings.Response[i].Tasks {
			user.Tasks = append(user.Tasks, standings.Response[i].Tasks[j].StandingsUserTask)
		}
		res.Users = append(res.Users, user)
	}

	return &res, nil
}
