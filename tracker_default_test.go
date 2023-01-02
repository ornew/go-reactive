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
package reactive_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ornew/go-reactive"
	"github.com/ornew/go-reactive/effect"
	"github.com/ornew/go-reactive/ref"
)

func TestSimpleChannel_Ref(t *testing.T) {
	ctx := context.TODO()
	tr := &reactive.SingleChannel{}
	tr.Start(ctx)

	t.Log("a=100")
	a := ref.New(tr, 100)

	// no effect tracking context - no track
	assert.Equal(t, 100, a.Get())
	// assert.Len(t, tr.subs, 0)

	var effectCount int
	var b int
	effect.Track(func() {
		// effect tracking context - track a
		b = a.Get() + 1
		t.Logf("b=%v computed", b)
		effectCount++
	})
	assert.Equal(t, 1, effectCount)
	// assert.Len(t, tr.subs, 1)
	assert.Equal(t, 101, b)

	t.Log("a=200")
	a.Set(200)

	if assert.Eventually(t, func() bool { return b == 201 }, time.Second, 10*time.Millisecond) {
		assert.Equal(t, 2, effectCount)
		assert.Equal(t, 200, a.Get())
		assert.Equal(t, 201, b)
	}

	t.Log("c=a")
	c := a
	t.Log("c=300")
	c.Set(300)
	// for k, s := range tr.subs {
	// 	t.Logf("key %p, sub: %p", k, s)
	// }
	// assert.Len(t, tr.subs, 1)

	if assert.Eventually(t, func() bool { return b == 301 }, time.Second, 10*time.Millisecond) {
		assert.Equal(t, 3, effectCount)
		assert.Equal(t, 300, a.Get())
		assert.Equal(t, 301, b)
		assert.Equal(t, 300, c.Get())
	}

	t.Log("d=0")
	d := ref.New(tr, 0)
	var e int
	effect.Track(func() {
		// effetc tracking context - track a
		e = d.Get() + 1
		t.Logf("e=%v computed", e)
	})
	t.Log("d=1")
	d.Set(1)
	if assert.Eventually(t, func() bool { return e == 2 }, time.Second, 10*time.Millisecond) {
		assert.Equal(t, 1, d.Get())
		assert.Equal(t, 2, e)
	}
}

func TestSimpleChannel_Compute(t *testing.T) {
	ctx := context.TODO()
	tr := &reactive.SingleChannel{}
	tr.Start(ctx)

	t.Log("a=100")
	a := ref.New(tr, 100)
	b := ref.Computed(tr, func() string {
		return strconv.Itoa(a.Get())
	})
	t.Logf("b=%q computed", b.Get())
	assert.Equal(t, 100, a.Get())
	assert.Equal(t, "100", b.Get())

	t.Log("a=12345")
	a.Set(12345)
	if assert.Eventually(t, func() bool { return b.Get() == "12345" }, time.Second, 10*time.Millisecond) {
		assert.Equal(t, 12345, a.Get())
		t.Logf("b=%q computed", b.Get())
		assert.Equal(t, "12345", b.Get())
	}
}
