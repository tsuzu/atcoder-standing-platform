package atcoder

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

type Client struct {
	ID, Pass string
	Client   *http.Client
	language string
	signedIn bool
}

func NewClient(id, pass string) (*Client, error) {
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}
	return &Client{
		ID:   id,
		Pass: pass,
		Client: &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) != 0 && via[0].URL.Path == "/login" {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}, nil
}

func (client *Client) Login() error {
	val := &url.Values{}
	val.Add("name", client.ID)
	val.Add("password", client.Pass)

	req, err := http.NewRequest("POST", "https://practice.contest.atcoder.jp/login", strings.NewReader(val.Encode()))

	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "AtCoder-Standing-Platform/0.01β")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "deflate, gzip")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 302 {
		return errors.New("Failed to log in to AtCoder")
	}

	client.signedIn = true
	return nil
}

func (client *Client) SetLanguage(lang string) {
	cookie := &http.Cookie{
		Name:   "language",
		Value:  lang,
		Domain: ".atcoder.jp",
		Path:   "/",
		MaxAge: int(1e9),
	}
	url, _ := url.Parse(".atcoder.jp")

	client.Client.Jar.SetCookies(url, []*http.Cookie{cookie})

	client.language = lang
}

func (client *Client) GetLanguage() string {
	return client.language
}

func (client *Client) getContests(pattern ...string) ([]Contest, error) {
	target := "https://atcoder.jp/contest"

	if client.language != "" {
		target = target + "?lang=" + url.QueryEscape(client.language)
	}

	query, err := goquery.NewDocument(target)

	if err != nil {
		return nil, err
	}

	subreg := regexp.MustCompile("(.*).contest.atcoder.jp")
	contests := make([]Contest, 0, 5)
	query.Find("h3").Each(func(_ int, s *goquery.Selection) {
		ok := false
		for i := range pattern {
			if pattern[i] == s.Text() {
				ok = true
				break
			}
		}
		if ok {
			s.NextAllFiltered("div").First().Find("tbody").Children().Has("a").Each(func(_ int, s *goquery.Selection) {
				var info Contest
				s.Find("td").EachWithBreak(func(i int, s *goquery.Selection) bool {
					if i == 0 {
						t := s.Find("a").Text()
						info.StartTime, _ = time.ParseInLocation("2006/01/02 15:04", t, JapanZone)
					} else if i == 1 {
						t := s.Find("a")
						var ok bool
						info.URL, ok = t.Attr("href")
						info.Title = t.Text()

						if ok {
							u, err := url.Parse(info.URL)

							if err == nil {
								if arr := subreg.FindStringSubmatch(u.Host); len(arr) > 1 {
									info.TitleInURL = arr[1]
								} else {
									logrus.Warning("The contest(" + info.URL + ") may not be at *.atcoder.jp")
								}
							}
						}
					} else if i == 2 {
						str := s.Text()

						t, err := time.Parse("15:04", str)

						if err == nil {
							after := t.Sub(invalidFinishTime)
							info.FinishTime = info.StartTime.Add(after)
						}
					} else {
						return false
					}
					return true
				})
				info.client = client.Client
				contests = append(contests, info)
			})
		}
	})

	return contests, nil
}

func (client *Client) GetUpcomingContests() ([]Contest, error) {
	return client.getContests("Upcoming Contests", "予定されたコンテスト")
}

func (client *Client) GetActiveContests() ([]Contest, error) {
	return client.getContests("Active Contests", "開催中のコンテスト")
}

func (client *Client) GetRecentContests() ([]Contest, error) {
	return client.getContests("Recent Contests", "終了後のコンテスト(最新10件)")
}
