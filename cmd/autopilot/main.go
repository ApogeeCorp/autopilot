// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	log    = logrus.New()
	config = viper.New()
)

func main() {
	switch os.Getenv("LOG_FORMAT") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
		}
	case "json":
		fallthrough
	default:
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := logrus.ParseLevel(lvl); err == nil {
			log.SetLevel(level)
		}
	}

	app := cli.NewApp()

	app.Name = "autopilot"
	app.Version = "0.0.1"
	app.Usage = "Generate recommendations from Prometheus metrics"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,f",
			Usage:  "The path to the configuration file",
			EnvVar: "CONFIG_FILE",
		},
	}

	app.Before = loadConfig
	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func loadConfig(c *cli.Context) error {
	if c.Command.FullName() == "help" {
		return nil
	}
	if c.String("config") == "" {
		return fmt.Errorf("missing config parameter")
	}
	return nil
}

func run(c *cli.Context) error {

	return nil
}
