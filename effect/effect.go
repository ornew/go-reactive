/*
Copyright 2022 Arata Furukawa.

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
package effect

import (
	"sync"
)

var (
	mu     sync.RWMutex
	active *Effect
)

func GetActive() *Effect {
	mu.RLock()
	e := active
	mu.RUnlock()
	return e
}

func setActive(e *Effect) {
	mu.Lock()
	if active != nil {
		mu.Unlock()
		panic("nesting effect is not allowed.")
	}
	active = e
	mu.Unlock()
}

func clearActive() {
	mu.Lock()
	// if active == nil {
	// 	panic("acrive effect is not found.")
	// }
	active = nil
	mu.Unlock()
}

type Effect struct {
	effect func()
}

func (e *Effect) Do() {
	setActive(e)
	e.effect()
	clearActive()
}

func Track(effect func()) {
	e := &Effect{effect: effect}
	e.Do()
}
