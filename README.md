# Go Reactive

```go
t := &tracker.SingleChannel{}
t.Start(context.TODO())

a := ref.New(100, t)
b := ref.Compute(func() int {
	return a.Get() + 1
}, t)

// a=100, b=101

a.Set(200)

// a=200, b=201
```
