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
package ref

import (
	"github.com/ornew/go-reactive/effect"
)

type Computed[T any] struct {
	r Ref[T]
}

func (c *Computed[T]) Get() T {
	return c.r.Get()
}

func Compute[T comparable](fn func() T, t Tracker) (c Computed[T]) {
	c.r.t = t
	effect.Track(func() {
		if c.r.v == nil {
			c.r.Set(fn())
			return
		}
		a := *c.r.v // Do not use Get() to avoid self-track.
		b := fn()
		if a != b {
			c.r.Set(b)
		}
	})
	return c
}
