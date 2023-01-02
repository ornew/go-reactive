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

import "github.com/ornew/go-reactive/effect"

type Key struct {
	Ptr  any
	Name string
}

type Tracker interface {
	Track(key Key)
	Trigger(key Key)
}

type Ref[T any] struct {
	t Tracker
	v *T
}

func (a *Ref[T]) Set(v T) {
	if a.v == nil {
		a.v = new(T)
	}
	*a.v = v
	a.t.Trigger(Key{
		Ptr: a.v,
	})
}

func (a *Ref[T]) Get() T {
	if a.v == nil {
		a.v = new(T)
	}
	a.t.Track(Key{
		Ptr: a.v,
	})
	return *a.v
}

func New[T any](v T, t Tracker) Ref[T] {
	return Ref[T]{
		t: t,
		v: &v,
	}
}

func Compute[T comparable](fn func() T, t Tracker) (r Ref[T]) {
	r.t = t
	effect.Track(func() {
		if r.v == nil {
			r.Set(fn())
			return
		}
		a := *r.v // Do not use Get() to avoid self-track.
		b := fn()
		if a != b {
			r.Set(b)
		}
	})
	return r
}
