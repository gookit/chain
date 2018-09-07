# middleware chain

[![GoDoc](https://godoc.org/github.com/gookit/chain?status.svg)](https://godoc.org/github.com/gookit/chain)
[![Build Status](https://travis-ci.org/gookit/chain.svg?branch=master)](https://travis-ci.org/gookit/chain)
[![Coverage Status](https://coveralls.io/repos/github/gookit/chain/badge.svg?branch=master)](https://coveralls.io/github/gookit/chain?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/chain)](https://goreportcard.com/report/github.com/gookit/chain)

A simple http middleware chain implement.

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
	middleware0 := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("a"))
			h.ServeHTTP(w, r)
			w.Write([]byte("A"))
		})
	}
	middleware1 := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("b"))
			h.ServeHTTP(w, r)
			w.Write([]byte("B"))
		})
	}

	c := chain.New(middleware0, middleware1)

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
	// Output: abc-(CORE)-CBA
}
```

## License

**MIT**
