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

type Subscriber[T any] struct {
	ch   chan T
	Ch   <-chan T
	bc   *Topic[T]
	prev *Subscriber[T]
	next *Subscriber[T]
}

func (s *Subscriber[T]) Pull() T {
	return <-s.Ch
}

func (s *Subscriber[T]) Unsubscribe() {
	if s.prev != nil {
		s.prev.next = s.next
	}
	if s.next != nil {
		s.next.prev = s.prev
	} else {
		s.bc.last = s.prev
	}
	s.bc = nil
	s.prev = nil
	s.next = nil
	close(s.ch)
}
