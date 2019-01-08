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
	"github.com/libopenstorage/autopilot/engine"
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
	}

	app.Before = func(c *cli.Context) error {
		switch os.Getenv("LOG_FORMAT") {
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

		if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
			if level, err := log.ParseLevel(lvl); err == nil {
				log.SetLevel(level)
			}
		}

		return nil
	}

	app.Action = func(c *cli.Context) error {
		log.SetLevel(log.DebugLevel)

		eng, err := engine.NewEngine()
		if err != nil {
			return err
		}

		if err := eng.Start(); err != nil {
			return err
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
