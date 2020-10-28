package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/llugin/freesound-downloader/client"
	"github.com/llugin/freesound-downloader/config"
)

func main() {
	listFlag := flag.Bool("l", false, "list results from page")
	pageFlag := flag.Int("p", 1, "page used for download and list")
	flag.Parse()

	c := &client.Client{
		Config: config.Config,
	}

	if *listFlag {
		err := list(c, *pageFlag)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := download(c, *pageFlag)
	if err != nil {
		log.Fatal(err)
	}
}

func download(c *client.Client, page int) error {
	authCode, err := c.Authorize()
	if err != nil {
		return err
	}

	query := client.Query{
		Query:    "",
		MaxLen:   60,
		PageSize: 16,
		Page:     page,
	}
	res, err := c.GetNewest(query)
	if err != nil {
		return err
	}

	ar, err := c.GetAccessToken(authCode)
	if err != nil {
		return err
	}

	err = c.Download(res, ar.AccessToken)
	if err != nil {
		return err
	}

	return nil
}

func list(c *client.Client, page int) error {
	query := client.Query{
		Query:    "",
		MaxLen:   60,
		PageSize: 16,
		Page:     page,
	}
	res, err := c.GetNewest(query)
	if err != nil {
		return err
	}
	for _, r := range res.Results {
		fmt.Printf("%+v\n", r)
	}
	return nil
}
