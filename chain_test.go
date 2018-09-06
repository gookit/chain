package chain

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	c := New(func(h http.Handler) http.Handler {
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

	// create fake request
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "abc-(CORE)-CBA" {
		t.Errorf("error res body: %s", w.Body.String())
	}

	h = c.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("-(CORE1)-"))
	})
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "abc-(CORE1)-CBA" {
		t.Errorf("error res body: %s", w.Body.String())
	}

	h = c.Wrap(nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("error res status code: %d", w.Code)
	}

	if w.Body.String() != "abc404 page not found\nCBA" {
		t.Errorf("error res body: %q", w.Body.String())
	}
}