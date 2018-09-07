// Package chain is a simple http middleware chain implement.
package chain

import (
	"net/http"
)

// Middleware definition
type Middleware func(h http.Handler) http.Handler

// Chain middleware chains
type Chain []Middleware

// New a middleware chain.
// Usage:
// 	c := chain.New(middleware0, middleware1)
//	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("hello"))
// 		w.WriteHeader(200)
// 	})
// 	handler := c.Wrap(myHandler)
// 	http.ListenAndServe(":8080", handler)
func New(chain ...Middleware) Chain {
	return Chain(chain)
}

// Use more middleware
func (c *Chain) Use(mds ...Middleware) Chain {
	*c = append(*c, mds...)
	return *c
}

// Extend old chain, returns new chain
func (c Chain) Extend(mds ...Middleware) Chain {
	nc := make([]Middleware, len(c)+len(mds))
	// because make() add length, so cannot use append()
	// nc = append(nc, mds...)
	copy(nc, mds)
	copy(nc[len(mds):], c)

	return nc
}

// Wrap all middleware to the core http handler
func (c Chain) Wrap(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	max := len(c)
	lst := make([]int, max)

	// warp all handlers
	for i := range lst {
		current := max - i - 1
		h = c[current](h)
	}

	return h
}

// WrapFunc all middleware to the core http handler func
func (c Chain) WrapFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Wrap(nil)
	}

	return c.Wrap(fn)
}
