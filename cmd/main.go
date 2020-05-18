package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/Ashish-Bansal/redis-spectacles/cmd/interactive"
	"github.com/Ashish-Bansal/redis-spectacles/cmd/noninteractive"
)

const redisURLArgName string = "url"

func main() {
	redisURLFlag := &cli.StringFlag{
		Name:     redisURLArgName,
		Usage:    "Redis URL to scan",
		Required: true,
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "print",
				Usage: "Print all redis keyspace key prefixes",
				Action: func(c *cli.Context) error {
					noninteractive.ExecuteNonInteractive(c)
					return nil
				},
				Flags: []cli.Flag{
					redisURLFlag,
				},
			},
			{
				Name:  "interactive",
				Usage: "Starts interactive console to visualise prefixes",
				Action: func(c *cli.Context) error {
					interactive.ExecuteInteractive(c)
					return nil
				},
				Flags: []cli.Flag{
					redisURLFlag,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
