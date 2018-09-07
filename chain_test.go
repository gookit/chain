package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
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

	c := New(middleware0, middleware1)

	middleware2 := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("c"))
			h.ServeHTTP(w, r)
			w.Write([]byte("C"))
		})
	}
	c.Use(middleware2)

	coreHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("-(CORE)-"))
		w.WriteHeader(200)
	}

	// create fake request
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// test Wrap
	h := c.Wrap(http.HandlerFunc(coreHandler))
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "abc-(CORE)-CBA" {
		t.Errorf("error res body: %s", w.Body.String())
	}

	// test WrapFunc
	w = httptest.NewRecorder()
	h = c.WrapFunc(coreHandler)
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "abc-(CORE)-CBA" {
		t.Errorf("error res body: %s", w.Body.String())
	}

	w = httptest.NewRecorder()
	h = c.WrapFunc(nil)
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "abc404 page not found\nCBA" {
		t.Errorf("error res body: %q", w.Body.String())
	}

	// test Extend
	nc := c.Extend(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("0"))
			h.ServeHTTP(w, r)
			w.Write([]byte("1"))
		})
	})

	w = httptest.NewRecorder()
	h = nc.WrapFunc(coreHandler)
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "0abc-(CORE)-CBA1" {
		t.Errorf("error res body: %s", w.Body.String())
	}
}
