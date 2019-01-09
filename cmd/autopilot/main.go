/*
Copyright 2019 Openstorage.org

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/libopenstorage/autopilot/telemetry/providers"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "autopilot"
	app.Version = "0.0.1"
	app.Usage = "Autopilot Storage Optimization Engine"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,f",
			Usage:  "set the configuration file path",
			EnvVar: "CONFIG_FILE",
			Value:  "/etc/autopilot/config.yaml",
		},
		cli.StringFlag{
			Name:   "data-dir,d",
			Usage:  "set the data directory for the process",
			EnvVar: "DATA_DIR",
			Value:  "/var/run/autopilot",
		},
		cli.StringFlag{
			Name:   "log-level",
			Usage:  "set the log level",
			EnvVar: "LOG_LEVEL",
			Value:  "info",
		},
		cli.StringFlag{
			Name:   "log-format",
			Usage:  "set the log format",
			EnvVar: "LOG_FORMAT",
			Value:  "text",
		},
	}

	app.Before = setupLog

	app.Action = func(c *cli.Context) error {
		// create our CRD definitions
		if err := createCRD(c); err != nil {
			return err
		}

		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name: "policy",
			Subcommands: []cli.Command{
				cli.Command{
					Name:      "test",
					Action:    policyTestAction,
					Usage:     "Test a policy document using the configuration",
					UsageText: "test <file>",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func setupLog(c *cli.Context) error {
	// setup the log format
	switch c.String("log-format") {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	case "json":
		fallthrough
	default:
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	}

	// setup the log level
	if level, err := log.ParseLevel(c.String("log-level")); err == nil {
		log.SetLevel(level)
	}

	return nil
}
