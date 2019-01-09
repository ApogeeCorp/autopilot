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

package telemetry

import (
	"fmt"
	"strings"
	"sync"

	sparks "gitlab.com/ModelRocket/sparks/types"
)

var (
	provMu    sync.RWMutex
	providers = make(map[string]NewFunc)
)

// Register makes a telemetry provider available by the provided name.
// If Register is called twice with the same name or if provider is nil,
// it panics.
func Register(name string, provider NewFunc) {
	name = strings.ToLower(name)

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

// NewInstance creates a new telemetry provider
func NewInstance(name, params string) (Provider, error) {
	provMu.RLock()
	newFn, ok := providers[name]
	provMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("telemetry: unknown provider %q (forgotten import?)", name)
	}
	return newFn(sparks.ParseStringParams(params))
}
