package main

import (
	"github.com/codegangsta/cli"
	"github.com/elbuo8/carnegie/carnegie"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Carnegie"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file path",
		},
	}
	app.Action = func(c *cli.Context) {
		configPath := c.String("config")
		if configPath == "" {
			log.Println("config missing")
			return
		}
		config := viper.New()
		config.AddConfigPath(configPath)
		err := config.ReadInConfig()
		if err != nil {
			log.Println(err)
			return
		}
		lb, err := carnegie.New(config)
		if err != nil {
			log.Println(err)
			return
		}
		lb.Start()
	}

	app.Run(os.Args)
}
