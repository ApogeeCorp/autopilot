// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/libopenstorage/autopilot/api/autopilot"
	"github.com/libopenstorage/autopilot/api/autopilot/rest"
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

	app.Action = func(c *cli.Context) error {
		api := &autopilot.API{
			Log: log,
		}

		handler, err := rest.Handler(rest.Config{
			AutopilotAPI: api,
			Logger:       log,
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
		log.Fatal(s.ListenAndServe())

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
