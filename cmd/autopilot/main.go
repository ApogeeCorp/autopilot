// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	sparks "gitlab.com/ModelRocket/sparks/util"
)

var (
	log = logrus.New()
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

	app.Commands = []cli.Command{
		{
			Name:  "collect",
			Usage: "Collect telemetry data from the provider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "provider,p",
					Usage:  "The telemetery provider to use to collect data",
					EnvVar: "PROVIDER",
					Value:  "prometheus",
				},
				cli.StringFlag{
					Name:   "url,u",
					Usage:  "The base URL for the telemetry provider",
					EnvVar: "PROVIDER_URL",
				},
				cli.StringFlag{
					Name:   "args,a",
					Usage:  "The provider args in the format 'arg1=val1;arg2=val2'",
					EnvVar: "PROVIDER_ARGS",
				},
				cli.StringFlag{
					Name:   "out,o",
					Usage:  "The output directory",
					EnvVar: "PROVIDER_OUT",
				},
			},
			Action: collect,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func collect(c *cli.Context) error {
	urlParam := c.String("url")
	if urlParam == "" {
		return errors.New("missing url parameter")
	}
	outParam := c.String("out")
	if outParam == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		outParam = wd + "/out"
	}
	if err := os.MkdirAll(outParam, 0700); err != nil {
		return err
	}

	// Parse the args into a map
	args := sparks.KVMap(c.String("args"), "&")

	switch c.String("provider") {
	case "prometheus":
		prom := prometheus.Prometheus{
			URL: urlParam,
			Log: log,
		}

		return prom.Collect(args, outParam)
	default:
		return errors.New("invalid provider")
	}
}
