// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/libopenstorage/autopilot/api/autopilot"
	"github.com/libopenstorage/autopilot/api/autopilot/rest"
	"github.com/libopenstorage/autopilot/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
	app.Usage = "Autopilot Storage Optimization Engine"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,f",
			Usage:  "set the configuration file path",
			EnvVar: "CONFIG_FILE",
			Value:  "./etc/config.json",
		},
		cli.StringFlag{
			Name:   "data-dir",
			Usage:  "set the data directory for the process",
			EnvVar: "DATA_DIR",
			Value:  "./var/run/autopilot",
		},
	}

	app.Action = func(c *cli.Context) error {
		config := &config.Config{}

		log.SetLevel(logrus.DebugLevel)

		api := &autopilot.API{
			Log:     log,
			DataDir: c.GlobalString("data-dir"),
			Config:  config,
		}

		data, err := ioutil.ReadFile(c.GlobalString("config"))
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, config); err != nil {
			return err
		}

		log.Debugf("DataDir=%s", api.DataDir)

		if err := os.MkdirAll(api.DataDir, 0770); err != nil {
			return err
		}

		handler, err := rest.Handler(rest.Config{
			AutopilotAPI: api,
			Logger:       log,
			AuthBasicAuth: func(ctx context.Context, username string, pass string) (context.Context, interface{}, error) {
				return ctx, username, nil
			},
		})
		if err != nil {
			log.Fatalln(err)
		}
		s := &http.Server{
			Addr:           ":9000",
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Infof("starting server %s", s.Addr)
		log.Fatal(s.ListenAndServe())

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
