# Go Reactive

```go
import (
	"github.com/ornew/go-reactive"
	"github.com/ornew/go-reactive/ref"
)

func main() {
	t := &reactive.ChannelTracker{}

	a := ref.New(t, 100)
	b := ref.Computed(t, func() int {
		return a.Get() + 1
	})

	// a=100, b=101

	a.Set(200)

	// a=200, b=201
}
```
