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

	"github.com/libopenstorage/autopilot/api/autopilot/rest/operations/sample"
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/libopenstorage/autopilot/telemetry/providers/prometheus"
)

// RecommendationsGet Returns the recommendations for a particular sample
func (a *API) RecommendationsGet(ctx *Context, params sample.RecommendationsGetParams) middleware.Responder {
	rules := &types.RuleSet{}

	// lookup the rules
	if err := a.DB.First(&rules, "id=?", params.Rules).Error; err != nil {
		return sparks.NewError(err)
	}

	// lookup the sample
	s := &types.Sample{}
	if err := a.DB.First(s, "id=?", params.SampleID).Error; err != nil {
		return sparks.NewError(err)
	}

	// check for the sample files
	samplePath := path.Join(os.Getenv("SAMPLE_VOL"), s.ID.String())
	if _, err := os.Stat(samplePath); err != nil {
		return sparks.NewError(err)
	}

	recs, err := a.engine.Recommend(rules, samplePath)
	if err != nil {
		return sparks.NewError(err)
	}

	return sample.NewRecommendationsGetOK().WithPayload(recs)
}

// SampleCreate Create a new telemetry sample from the provided definition
func (a *API) SampleCreate(ctx *Context, params sample.SampleCreateParams) middleware.Responder {
	if *params.Type != string(types.SourceTypePrometheus) {
		return sparks.ErrInvalidParameter.Reason("invalid type")
	}

	source, err := ioutil.ReadAll(params.Sample)
	if err != nil {
		return sparks.NewError(err)
	}

	tx := a.DB.Begin()
	s := &types.Sample{
		Type: types.SourceTypePrometheus,
	}

	if err := tx.Create(s).Error; err != nil {
		tx.Rollback()
		return sparks.NewError(err)
	}

	// check for the sample files
	samplePath := path.Join(os.Getenv("SAMPLE_VOL"), s.ID.String())
	if err := os.MkdirAll(samplePath, 0600); err != nil {
		tx.Rollback()
		return sparks.NewError(err)
	}
	a.Log.Debugf("creating new sample at %s", samplePath)

	prom := &prometheus.Prometheus{
		Log: a.Log,
	}

	if err := prom.Parse(source, nil, samplePath); err != nil {
		tx.Rollback()
		return sparks.NewError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return sparks.NewError(err)
	}

	return sample.NewSampleCreateCreated().WithPayload(s)
}

// SampleDelete Returns the request collected object
func (a *API) SampleDelete(ctx *Context, params sample.SampleDeleteParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleDelete")
}

// SampleGet Returns the request collected object
func (a *API) SampleGet(ctx *Context, params sample.SampleGetParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleGet")
}

// SampleList Returns an array of telemetry samples defined in the system
func (a *API) SampleList(ctx *Context, params sample.SampleListParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleList")
}

// SampleUpdate Update the properties of the specified sample
func (a *API) SampleUpdate(ctx *Context, params sample.SampleUpdateParams) middleware.Responder {
	return sparks.ErrNotImplemented("sampleUpdate")
}
