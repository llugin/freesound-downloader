package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/llugin/freesound-downloader/client"
	"github.com/llugin/freesound-downloader/config"
)

func main() {
	authorize := flag.Bool("a", false, "authorize app")
	download := flag.String(
		"d",
		"",
		"download newest, provide access token",
	)
	list := flag.Bool("l", false, "list results")
	page := flag.Int("p", 1, "page")
	flag.Parse()

	c := client.Client{
		Config: config.Config,
	}

	if *authorize {
		err := c.Authorize()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *download != "" {
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

		ar, err := c.GetAccessToken(*download)
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
