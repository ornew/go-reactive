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
package reactive

import (
	"context"
	"sync"

	"github.com/ornew/go-reactive/effect"
	"github.com/ornew/go-reactive/ref"
)

type ChannelTracker struct {
	s   sync.Once
	mu  sync.RWMutex
	eff map[ref.Key][]*effect.Effect
	ch  chan ref.Key
}

func (t *ChannelTracker) Track(key ref.Key) {
	eff := effect.GetActive()
	if eff != nil {
		t.Bind(key, eff)
	}
}

func (t *ChannelTracker) Trigger(key ref.Key) {
	t.ch <- key
}

func (t *ChannelTracker) Start(ctx context.Context) {
	var w sync.WaitGroup
	w.Add(1)
	t.s.Do(func() {
		t.ch = make(chan ref.Key)
		w.Done()
		go func() {
			defer func() {
				close(t.ch)
				t.ch = nil
			}()
			effects := make([]*effect.Effect, 0, 16)
			for {
				select {
				case <-ctx.Done():
					return
				case key, ok := <-t.ch:
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
					if r > cap(effects) {
						effects = make([]*effect.Effect, 0, r*2)
					}
					effects = effects[:r]
					copy(effects, t.eff[key])

					// WARN: There is posibility that an effect contains Bind(),
					// should Unlock() before Do() to avoid a deadlock.
					t.mu.RUnlock()
					for _, e := range effects {
						e.Do()
					}
				}
			}
		}()
	})
	w.Wait() // Wait to make t.ch
}

// Bind runs the effect when the key is triggered.
func (t *ChannelTracker) Bind(key ref.Key, eff *effect.Effect) {
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

func (t *ChannelTracker) Unbind(key ref.Key, eff *effect.Effect) {
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
