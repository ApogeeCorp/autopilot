// Copyright 2018 Portworx Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package telemetry

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/libopenstorage/autopilot/api/autopilot/types"
)

var (
	provMu    sync.RWMutex
	providers = make(map[string]NewFunc)
)

// Register makes a telemetry provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(t types.ProviderType, provider NewFunc) {
	name := strings.ToLower(string(t))

	provMu.Lock()
	defer provMu.Unlock()
	if provider == nil {
		panic("telemetry: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("telemetry: Register called twice for provider " + name)
	}

	providers[name] = provider
}

// NewInstance creates a new telemetry provider instance with the specified parameters
func NewInstance(prov types.Provider) (Provider, error) {
	if prov == nil {
		return nil, errors.New("invalid provider")
	}
	name := strings.ToLower(string(prov.Type()))
	provMu.RLock()
	createFn, ok := providers[name]
	provMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("telemetry: unknown provider %q (forgotten import?)", name)
	}
	return createFn(prov)
}
