// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file at the
// root of this project.

package autopilot

import (
	"io/ioutil"
	"os"
	"path"

	sparks "gitlab.com/ModelRocket/sparks/types"

	"github.com/go-openapi/runtime/middleware"
	uuid "github.com/satori/go.uuid"

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations"
	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
)

// RecommendationsGet Create a new telemetry sample from the provided definition
func (a *API) RecommendationsGet(ctx *Context, params operations.RecommendationsGetParams) middleware.Responder {
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

	recs, err := a.engine.Recommend(a.Config.Rules, samplePath)
	if err != nil {
		return sparks.NewError(err)
	}

	return operations.NewRecommendationsGetOK().WithPayload(recs)
}
