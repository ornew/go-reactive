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

type Key struct {
	Ptr  any
	Name string
}

type Tracker interface {
	Trigger(key Key)
	Mark(key Key)
	Track(effect func())
}

type tracker struct {
	ctx tracking
	em  effectMap
}

func (t *tracker) Mark(key Key) {
	effect := t.ctx.Active()
	if effect != nil {
		t.em.Add(key, effect)
	}
}

func (t *tracker) Track(effectFunc func()) {
	e := effect{
		ctx: &t.ctx,
		fn:  effectFunc,
	}
	e.Do()
}

var DefaultTracker Tracker = &Channel{}
