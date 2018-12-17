// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package engine

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alecthomas/template"
	"github.com/go-openapi/strfmt"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/libopenstorage/autopilot/config"
	"github.com/libopenstorage/autopilot/telemetry"
	log "github.com/sirupsen/logrus"
)

// Engine is the autopilot recommendation engine
type Engine struct {
	config    *config.Config
	providers map[string]telemetry.Provider
	stop      chan bool
	wg        sync.WaitGroup
}

// NewEngine returns a new engine
func NewEngine(c *config.Config) (*Engine, error) {
	// initialize the provider instances
	provs := make(map[string]telemetry.Provider)
	for _, p := range c.Providers {
		prov, err := telemetry.Get(p)
		if err != nil {
			return nil, err
		}
		provs[p.Name()] = prov
	}

	return &Engine{
		providers: provs,
		config:    c,
		stop:      make(chan bool),
	}, nil
}

// Start starts the engine monitors and the collector scheduler
func (e *Engine) Start() error {
	log.Debug("staring engine...")

	// start the monitors
	for _, m := range e.config.Monitors {
		_, ok := e.providers[m.Provider]
		if !ok {
			e.Stop()
			return fmt.Errorf("invalid provider %s", m.Provider)
		}
		rules := make([]*types.Rule, 0)
		for _, r := range m.Rules {
			if rule, ok := e.config.GetRule(r); ok {
				rules = append(rules, rule)
			}
		}

		dur, err := time.ParseDuration(*m.Interval)
		if err != nil {
			e.Stop()
			return err
		}
		e.wg.Add(1)

		go func(prov string, interval time.Duration, rules []*types.Rule) {
			defer e.wg.Done()

			ticker := time.NewTicker(interval)

			for {
				select {
				case <-ticker.C:
					_, err := e.Recommend(prov, rules)
					if err != nil {
						log.Errorln(err)
					} else {
						// TODO: emit recommendations
					}
				case <-e.stop:
					break
				}
			}
		}(m.Provider, dur, rules)
	}

	return nil
}

// Stop stops the engine
func (e *Engine) Stop() {
	close(e.stop)
	e.wg.Wait()
}

// Recommend returns recommendations from the provider by name with the rules request
func (e *Engine) Recommend(name string, rules []*types.Rule) ([]*types.Recommendation, error) {
	rval := make([]*types.Recommendation, 0)

	prov, ok := e.providers[name]
	if !ok {
		return nil, errors.New("unknown provider")
	}

	for _, rule := range rules {
		vectors, err := prov.Query(telemetry.Params{"query": rule.Expr})
		if err != nil {
			return nil, err
		}

		recommendation := &types.Recommendation{
			Timestamp: strfmt.DateTime(time.Now()),
		}

		for _, v := range vectors {
			var proposalValue bytes.Buffer

			fmap := template.FuncMap{
				"formatAsDate": formatAsDate,
				"formatFloat":  formatFloat,
			}
			proposal := &types.Proposal{
				Rule:    rule.Name,
				Cluster: v.Metric.Cluster,
				Node:    v.Metric.Node,
				Volume:  *v.Metric.Volume,
			}

			t := template.Must(template.New("Issue").
				Funcs(fmap).
				Parse("{{index .Value 0 | formatAsDate}}: " + rule.Issue + ` ` + rule.Proposal))

			err := t.Execute(&proposalValue, v)
			if err != nil {
				return nil, err
			}

			proposal.Action = proposalValue.String()

			recommendation.Proposals = append(recommendation.Proposals, proposal)
		}
		rval = append(rval, recommendation)
	}

	return rval, nil
}

func formatAsDate(timestamp float64) string {
	unixTimeUTC := time.Unix(int64(timestamp), 0) //gives unix time stamp in utc
	return unixTimeUTC.Format(time.RFC3339)       // converts utc time to RFC3339 format
}

func formatFloat(value string) string {
	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%.1f", floatVal)
}
