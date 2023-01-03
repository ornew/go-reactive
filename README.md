# Go Reactive

## Quick Start

```go
import (
	"github.com/ornew/go-reactive"
	"github.com/ornew/go-reactive/ref"
)

func main() {
	a := ref.New(100)
	b := ref.Computed(func() int {
		return a.Get() + 1
	})

	// a=100, b=101

	a.Set(200)

	// a=200, b=201
}
```

## Guide

### Ref

A `Ref` is a reference to a reactive object.

```go
var aRef ref.Ref[int] = ref.New(100)
```

A `Ref` is always a reference to a value typed `T`.

In a tracking context, `Get` returns a value and tracks the dependency.

```go
bRef := ref.Computed(func() int {
	// Here is tracking context, returns a value of aRef, and tracks aRef.
	// If a value of aRef changed, a value of bRef will be updated automatically.
	return aRef.Get() + 1
})

// Here is no tracking context, just returns a value of aRef.
aRef.Get()
```

A `Set` sets a value of `Ref` and triggers effects tracked it.

```go
aRef.Set(200)
// bRef tracked aRef, it will be updated.
bRef.Get() // 201
```

If copy `Ref`, it references to a same value.

```go
aRef.Set(100)
cRef = aRef
cRef.Get() // 100
aRef.Set(0)
aRef.Get() // 0
cRef.Get() // 0
cRef.Set(100)
aRef.Get() // 100
cRef.Get() // 100
```

If creates a `Ref` or updates a value, always copied a value, so the original value will be disconnected.

```go
var b bool
aRef.Set(b)
aRef.Set(true)
// b is false, aRef is true.
```

A `Ref` is a **shallow** reference. It no touches values referenced by `T`.

```go
type Ptr struct {
	Deep *int
}
pRef = ref.New(Ptr{Deep: nil})
t.Watch(func() { fmt.Println("changed") }, pRef)

var v int
pRef.Get().Deep = &v
// Nothing happens.
```

A `Computed` creates a new `Ref` by a computing function.
When `Ref` in the computing function is changed, this computed `Ref` value will be also updated.

```go
aStr := ref.Computed(func() string {
	return fmt.Sprintf("computed from aRef: %v", aRef.Get())
})
aRef.Set(false)
aStr.Get() // "computed from aRef: false"
```

This can be done with `Track` as following, but `Computed` has performance advantages.

```go
aStr := ref.New("")
t.Track(func() {
	value := fmt.Sprintf("track aRef: %v", aRef.Get())
	aStr.Set(value)
})
aRef.Set(false)
aStr.Get() // "track aRef: false"
```

### Tracker

A `Tracker` observes a `Ref` changes.

```go
var t tracker.Tracker = reactive.DefaultTracker
```

Some func such as ref.New can change tracker by args.

```go
// Use reactive.DefaultTracker.
ref.New(100)
// Use a custom tracker.
ref.New(100, ref.WithTracker(t))
```

A `Ref` belong only to one tracking scope.
If tries to track other tracker's `Ref`, defaults to error.

```go
aRef := ref.New(100, trackerA)

// Error.
trackerB.Track(func() { fmt.Println(aRef.Get()) })
// No error, but this doesn't track aRef so aRef never trigger this.
trackerB.Track(func() { fmt.Println(aRef.Unref()) })
```

A `Track` tracks automatically all `Ref` without explicitly specifying that the side-effect function depends on it.

```go
effect := func() {
	 fmt.Printf("aRef is %v", aRef.Get())
}

// The `effect` tracks automatically `aRef`.
// A `Track` will happen the side-effect in first time.
finish, _ := t.Track(effect) // "aRef is true"
```

<!--
A `Watch` also tracks `Ref`, but the `Ref` must be explicitly specified, and no tracks automatically.

```go
t.Track(effect)

// Should be explicitly specified dependencies.
reactive.Watch(t, effect, aRef)
reactive.Watch(t, effect, aRef, bRef)
```
-->

When a `Ref` value is changed, a `Tracker` triggers side-effects tracked it.

```go
aRef.Set(false) // "aRef is false"
```

A function returned by Track finishes the tracking a Ref of the side-effect.

```go
finish()

aRef.Set(true) // Nothing happens.
```
