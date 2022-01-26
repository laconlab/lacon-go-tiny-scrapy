package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeAgents struct {
	agent string
}

func (a *fakeAgents) next() string {
	return a.agent
}

func Test200(t *testing.T) {
	agent := "Test-Agent-1"

	call := make(chan bool, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["User-Agent"][0] != agent {
			t.Error("Expected User-Agent: ", agent, " got: ", r.Header)
		}

		w.WriteHeader(http.StatusOK)

		call <- true
	}))
	defer ts.Close()

	resp := headerFilter(ts.URL, agent, time.Second)

	select {
	case <-call:
	default:
		t.Error("Server never called")
	}

	if !resp {
		t.Error("Expected true got false")
	}

}

func Test400(t *testing.T) {
	agent := "Test-Agent-1"

	call := make(chan bool, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["User-Agent"][0] != agent {
			t.Error("Expected User-Agent: ", agent, " got: ", r.Header)
		}

		call <- true
		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer ts.Close()

	resp := headerFilter(ts.URL, agent, time.Second)

	select {
	case <-call:
	default:
		t.Error("Server never called")
	}

	if resp {
		t.Error("Expected false got true")
	}
}

func Test500(t *testing.T) {
	agent := "Test-Agent-1"

	call := make(chan bool, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["User-Agent"][0] != agent {
			t.Error("Expected User-Agent: ", agent, " got: ", r.Header)
		}

		call <- true
		http.Error(w, "Ups", http.StatusInternalServerError)
	}))
	defer ts.Close()

	resp := headerFilter(ts.URL, agent, time.Second)

	select {
	case <-call:
	default:
		t.Error("Server never called")
	}

	if !resp {
		t.Error("Expected true got false")
	}
}

func TestTimeout(t *testing.T) {
	agent := "Test-Agent-1"

	call := make(chan bool, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["User-Agent"][0] != agent {
			t.Error("Expected User-Agent: ", agent, " got: ", r.Header)
		}

		call <- true

		time.Sleep(100 * time.Millisecond)
		http.Error(w, "Not found", http.StatusNotFound)
	}))
	defer ts.Close()

	resp := headerFilter(ts.URL, agent, 10*time.Millisecond)

	select {
	case <-call:
	default:
		t.Error("Server never called")
	}

	if !resp {
		t.Error("Expected true got false")
	}
}
