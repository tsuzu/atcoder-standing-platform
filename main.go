package main

import (
	"fmt"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/cs3238-tsuzu/atcoder-standing-platform/atcoder"
	"github.com/k0kubun/pp"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(os.Args[0], "[id]", "[password]")
	}

	client, err := atcoder.NewClient(os.Args[1], os.Args[2])

	if err != nil {
		panic(err)
	}

	fmt.Println(client.Login())

	client.SetLanguage("ja")

	if contests, err := client.GetRecentContests(); err != nil {
		panic(err)
	} else {
		for i := range contests {
			s, err := contests[i].Standings()

			if err != nil {
				logrus.Error(err)
			} else {
				pp.Println(s)
			}

			if ok, err := contests[i].IsJoined(); err != nil {
				logrus.Println(err)
			} else {
				logrus.Println(contests[i].Title, ":", ok)
			}
		}
	}

	//	fmt.Println(string(b))
}
