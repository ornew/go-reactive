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

import (
	"context"
	"sync"

	"github.com/ornew/go-reactive/effect"
	"github.com/ornew/go-reactive/pool"
	"github.com/ornew/go-reactive/pubsub"
	"github.com/ornew/go-reactive/ref"
)

type SingleChannel struct {
	mu  sync.RWMutex
	eff map[ref.Key][]*effect.Effect
	top pubsub.Topic[ref.Key]
}

func (t *SingleChannel) Start(ctx context.Context) {
	s := t.top.Subscribe()
	go func() {
		defer s.Unsubscribe()
		p := pool.NewSlice[*effect.Effect](0)
		for {
			select {
			case <-ctx.Done():
				return
			case key, ok := <-s.Ch:
				if !ok {
					return
				}
				t.mu.RLock()
				r := len(t.eff[key])
				if r == 0 {
					t.mu.RUnlock()
					continue
				}

				// Copy effects.
				eff, put := p.Get(r)
				eff = eff[:r]
				copy(eff, t.eff[key])

				// WARN: There is posibility that an effect contains Link(),
				// should Unlock() before Do() to avoid a deadlock.
				t.mu.RUnlock()
				for _, e := range eff {
					e.Do()
				}
				put(eff) // Return to pool.
			}
		}
	}()
}

// Link runs the effect when the key is triggered.
func (t *SingleChannel) Link(key ref.Key, eff *effect.Effect) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.eff == nil {
		t.eff = map[ref.Key][]*effect.Effect{}
	}
	effs, ok := t.eff[key]
	if !ok {
		t.eff[key] = []*effect.Effect{}
	}
	for _, e := range effs {
		if e == eff {
			return
		}
	}
	t.eff[key] = append(effs, eff)
}

func (t *SingleChannel) Unlink(key ref.Key, eff *effect.Effect) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.eff == nil {
		return
	}
	effs, ok := t.eff[key]
	if !ok {
		return
	}
	for i, e := range effs {
		if e == eff {
			neweffs := effs[:i]
			if len(effs) > i+1 {
				neweffs = append(neweffs, effs[i+1:]...)
			}
			t.eff[key] = neweffs
			return
		}
	}
}

func (t *SingleChannel) Track(key ref.Key) {
	eff := effect.GetActive()
	if eff != nil {
		t.Link(key, eff)
	}
}

func (t *SingleChannel) Trigger(key ref.Key) {
	t.top.Publish(key)
}
