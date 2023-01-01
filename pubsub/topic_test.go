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

import (
	"testing"
)

func TestTopic(t *testing.T) {
	var c Topic[int]
	src1 := make(chan int)
	c.Source(src1)
	sub1 := c.Subscribe()
	sub2 := c.Subscribe()
	src1 <- 100
	src1 <- 200
	src1 <- 300
	v, ok := <-sub1.Ch
	t.Log(v, ok)
	v, ok = <-sub1.Ch
	t.Log(v, ok)
	v, ok = <-sub1.Ch
	t.Log(v, ok)
	v, ok = <-sub2.Ch
	t.Log(v, ok)
	v, ok = <-sub2.Ch
	t.Log(v, ok)
	v, ok = <-sub2.Ch
	t.Log(v, ok)
	sub1.Unsubscribe()
	src1 <- 400
	v, ok = <-sub1.Ch
	t.Log(v, ok)
	v, ok = <-sub2.Ch
	t.Log(v, ok)
}
