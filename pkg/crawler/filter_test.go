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

func Test200 (t *testing.T){
    ua := "Test-Agent-1"

    call := make(chan bool, 1)
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header["User-Agent"][0] != ua {
            t.Error("Expected User-Agent: ", ua, " got: ", r.Header)
        }

        call<- true
    }))
    defer ts.Close()

    agents := &fakeAgents{
        agent: ua,
    }

    headerFilter := newHeaderFilter(agents, time.Second)

    resp := headerFilter.filter(ts.URL)

    select{
    case <-call:
    default:
        t.Error("Server never called")
    }

    if !resp {
        t.Error("Expected true got false")
    }

}

func Test400(t *testing.T) {
    ua := "Test-Agent-1"

    call := make(chan bool, 1)
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header["User-Agent"][0] != ua {
            t.Error("Expected User-Agent: ", ua, " got: ", r.Header)
        }

        call<- true
        http.Error(w, "Not found", http.StatusNotFound)
    }))
    defer ts.Close()

    agents := &fakeAgents{
        agent: ua,
    }

    headerFilter := newHeaderFilter(agents, time.Second)

    resp := headerFilter.filter(ts.URL)

    select{
    case <-call:
    default:
        t.Error("Server never called")
    }

    if resp {
        t.Error("Expected false got true")
    }
}

func TestTimeout(t *testing.T) {
    ua := "Test-Agent-1"

    call := make(chan bool, 1)
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header["User-Agent"][0] != ua {
            t.Error("Expected User-Agent: ", ua, " got: ", r.Header)
        }

        call<- true

        time.Sleep(2 * time.Second)
        http.Error(w, "Not found", http.StatusNotFound)
    }))
    defer ts.Close()

    agents := &fakeAgents{
        agent: ua,
    }

    headerFilter := newHeaderFilter(agents, time.Second)

    resp := headerFilter.filter(ts.URL)

    select{
    case <-call:
    default:
        t.Error("Server never called")
    }

    if !resp {
        t.Error("Expected true got false")
    }
}
