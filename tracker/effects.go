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
package tracker

import "sync"

type effect struct {
	ctx *tracking
	fn  func()
}

func (e *effect) Do() {
	e.ctx.Activate(e)
	defer e.ctx.Inactivate()
	e.fn()
}

// There is posibility that an effect contains Add(),
// should m.Unlock() before Do() to avoid a deadlock.
type effectMap struct {
	m    sync.RWMutex
	link map[Key][]*effect
}

func (t *effectMap) Load(effects *[]*effect, key Key) {
	t.m.RLock()
	r := len(t.link[key])
	if r == 0 {
		t.m.RUnlock()
		(*effects) = (*effects)[0:0]
		return
	}

	// Copy effects.
	if r > cap(*effects) {
		*effects = make([]*effect, 0, r*2)
	}
	*effects = (*effects)[:r]
	copy(*effects, t.link[key])

	t.m.RUnlock()
}

// bind runs the effect when the key is triggered.
func (t *effectMap) Add(key Key, eff *effect) {
	t.m.Lock()
	defer t.m.Unlock()
	if t.link == nil {
		t.link = map[Key][]*effect{}
	}
	effs, ok := t.link[key]
	if !ok {
		t.link[key] = []*effect{}
	}
	for _, e := range effs {
		if e == eff {
			return
		}
	}
	t.link[key] = append(effs, eff)
}

func (t *effectMap) Remove(key Key, eff *effect) {
	t.m.Lock()
	defer t.m.Unlock()
	if t.link == nil {
		return
	}
	effs, ok := t.link[key]
	if !ok {
		return
	}
	for i, e := range effs {
		if e == eff {
			neweffs := effs[:i]
			if len(effs) > i+1 {
				neweffs = append(neweffs, effs[i+1:]...)
			}
			t.link[key] = neweffs
			return
		}
	}
}
