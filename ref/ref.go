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

import "github.com/ornew/go-reactive/tracker"

type Ref[T any] struct {
	t tracker.Tracker
	v *T
}

func (a *Ref[T]) Key() tracker.Key {
	return tracker.Key{
		Ptr: a.v,
	}
}

func (a *Ref[T]) Set(v T) {
	if a.v == nil {
		a.v = new(T)
	}
	*a.v = v
	a.t.Trigger(a.Key())
}

func (a *Ref[T]) Get() T {
	if a.v == nil {
		a.v = new(T)
	}
	a.t.Mark(a.Key())
	return *a.v
}

func (a *Ref[T]) Unref() (v T) {
	if a.v == nil {
		return v
	}
	return *a.v
}

type refOptions struct {
	t tracker.Tracker
}

type RefOption func(*refOptions)

func WithTracker(t tracker.Tracker) RefOption {
	return func(o *refOptions) {
		o.t = t
	}
}

func Zero[T any](opts ...RefOption) Ref[T] {
	opt := refOptions{}
	for _, o := range opts {
		o(&opt)
	}
	if opt.t == nil {
		opt.t = tracker.DefaultTracker
	}
	return Ref[T]{
		t: opt.t,
	}
}

func New[T any](v T, opts ...RefOption) Ref[T] {
	r := Zero[T](opts...)
	r.v = &v
	return r
}

func Computed[T comparable](fn func() T, opts ...RefOption) Ref[T] {
	r := Zero[T](opts...)
	r.t.Track(func() {
		v := fn()
		// Do not use Get() to avoid self-track, use Unref().
		if r.Unref() != v {
			r.Set(v)
		}
	})
	return r
}
