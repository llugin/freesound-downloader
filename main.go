package main

import (
	"flag"
	"fmt"
	"log"

	"./client"
	"./config"
)

func main() {

	authorize := flag.Bool("a", false, "authorize app")
	download := flag.String("d", "", "download newest, provide authorization")

	flag.Parse()

	c := client.Client{
		Config: config.Config,
	}
	if *authorize {
		fmt.Println("authorize")
		err := c.Authorize()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *download != "" {
		fmt.Println("download")
		query := client.Query{
			Query:  "",
			MaxLen: 60,
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
}
