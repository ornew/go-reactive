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
package pubsub

type Topic[T any] struct {
	last *Subscriber[T]
}

func (c *Topic[T]) Source(ch <-chan T) {
	go func() {
		for {
			val, ok := <-ch
			if !ok {
				return
			}
			c.Publish(val)
		}
	}()
}

func (c *Topic[T]) Subscribe() *Subscriber[T] {
	sub := &Subscriber[T]{
		ch:   make(chan T),
		bc:   c,
		prev: c.last,
		next: nil,
	}
	sub.Ch = sub.ch
	if sub.prev != nil {
		sub.prev.next = sub
	}
	c.last = sub
	return c.last
}

func (c *Topic[T]) Publish(val T) {
	sub := c.last
	for sub != nil {
		c := sub.ch
		go func() {
			c <- val
		}()
		sub = sub.prev
	}
}
