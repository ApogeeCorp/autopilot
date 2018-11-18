// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package engine

import (
	"github.com/libopenstorage/autopilot/api/autopilot/types"
	"github.com/sirupsen/logrus"
)

// Engine is the autopilot recommendation engine
type Engine struct {
	Log *logrus.Logger
}

// Recommend returns a recommendation from the engine based on the rules and sample
func (e *Engine) Recommend(rules *types.RuleSet, samplePath string) ([]*types.Recommendation, error) {
	e.Log.Debugf("processing data here %s", samplePath)

	return nil, nil
}
