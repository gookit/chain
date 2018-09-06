# middleware chain

Is a simple http middleware chain implement.

## Godoc

- [godoc for gopkg](https://godoc.org/gopkg.in/gookit/chain.v1)
- [godoc for github](https://godoc.org/github.com/gookit/chain)

## Usage

```go
package main

import (
	"github.com/gookit/chain"
	"net/http"
	"testing"
)

func main() {
	c := chain.New(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("a"))
			h.ServeHTTP(w, r)
			w.Write([]byte("A"))
		})
	}, func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("b"))
			h.ServeHTTP(w, r)
			w.Write([]byte("B"))
		})
	})

	c.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("c"))
			h.ServeHTTP(w, r)
			w.Write([]byte("C"))
		})
	})

	h := c.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("-(CORE)-"))
		w.WriteHeader(200)
	}))
 
	http.ListenAndServe(":8090", h)
}
```

## License

**MIT**
