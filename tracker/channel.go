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
	"sync"
)

type Channel struct {
	tracker
	start sync.Once
	ch    chan Key
}

var _ Tracker = (*Channel)(nil)

func (t *Channel) Trigger(key Key) {
	if t.ch == nil {
		t.Start()
	}
	t.ch <- key
}

func (t *Channel) Start() {
	t.start.Do(func() {
		t.ch = make(chan Key)
		go func() {
			var closed bool
			defer func() {
				if !closed {
					// Close when t.ch was not closed (e.g., panic in effect).
					close(t.ch)
				}
			}()
			buf := make([]*effect, 16)
			for {
				key, ok := <-t.ch
				if !ok {
					closed = true
					return
				}
				t.em.Load(&buf, key)
				for _, effect := range buf {
					// TODO: Handle panic.
					effect.Do()
				}
			}
		}()
	})
}
