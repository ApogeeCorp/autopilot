// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"
	uuid "github.com/satori/go.uuid"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
)

// RecommendationsGet Create a new telemetry sample from the provided definition
func (a *API) RecommendationsGet(ctx *Context, params operations.RecommendationsGetParams) middleware.Responder {
	rules := a.Config.Rules

	// read in the rules
	if params.Rules != nil {
		rules = make([]*types.Rule, 0)
		data, err := ioutil.ReadAll(params.Rules)
		if err != nil {
			return sparks.NewError(err)
		}

		if err := json.Unmarshal(data, &rules); err != nil {
			return sparks.NewError(err)
		}
	}

	// read in the source
	source, err := ioutil.ReadAll(params.Sample)
	if err != nil {
		return sparks.NewError(err)
	}

	sampleID, err := uuid.NewV4()
	if err != nil {
		return sparks.NewError(err)
	}

	// check for the sample files
	samplePath := path.Join(a.DataDir, sampleID.String())
	if err := os.MkdirAll(samplePath, 0770); err != nil {
		return sparks.NewError(err)
	}
	a.Log.Debugf("creating new sample at %s", samplePath)

	prom := &prometheus.Prometheus{
		Log: a.Log,
	}

	if err := prom.Parse(source, nil, samplePath); err != nil {
		return sparks.NewError(err)
	}

	recs, err := a.engine.Recommend(rules, samplePath)
	if err != nil {
		return sparks.NewError(err)
	}

	return operations.NewRecommendationsGetOK().WithPayload(recs)
}
