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

func TestTimeout(t *testing.T) {
	agent := "Test-Agent-1"

	call := make(chan bool, 1)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["User-Agent"][0] != agent {
			t.Error("Expected User-Agent: ", agent, " got: ", r.Header)
		}

		call <- true

		time.Sleep(2 * time.Second)
		http.Error(w, "Not found", http.StatusNotFound)
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
