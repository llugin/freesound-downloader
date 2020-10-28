package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/llugin/freesound-downloader/client"
	"github.com/llugin/freesound-downloader/config"
)

func main() {
	download := flag.Bool(
		"d",
		false,
		"download newest or results from page if provided",
	)
	list := flag.Bool("l", false, "list results")
	page := flag.Int("p", 1, "page")
	flag.Parse()

	c := client.Client{
		Config: config.Config,
	}

	if *download {
		authCode, err := c.Authorize()
		if err != nil {
			log.Fatal(err)
		}

		query := client.Query{
			Query:    "",
			MaxLen:   60,
			PageSize: 16,
			Page:     *page,
		}
		res, err := c.GetNewest(query)
		if err != nil {
			log.Fatal(err)
		}

		ar, err := c.GetAccessToken(authCode)
		if err != nil {
			log.Fatal(err)
		}

		err = c.Download(res, ar.AccessToken)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *list {
		query := client.Query{
			Query:    "",
			MaxLen:   60,
			PageSize: 16,
			Page:     *page,
		}
		res, err := c.GetNewest(query)
		if err != nil {
			log.Fatal(err)
		}
		for _, r := range res.Results {
			fmt.Printf("%+v\n", r)
		}
	}
}
