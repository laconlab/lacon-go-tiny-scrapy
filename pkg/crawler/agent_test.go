package crawler

import (
	"io/ioutil"
	"sync"
	"syscall"
	"testing"
)

func TestRoundRobinAgents(t *testing.T) {
	fileName := "tmp.yml"
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	cfg := `
    User-Agents:
        - "Agent-1"
        - "Agent-2"
        - "Agent-3"
    `
	ioutil.WriteFile(f.Name(), []byte(cfg), 0644)

	agents := newHttpAgents(f.Name())

	agent := agents.next()
	if agent != "Agent-1" {
		t.Error("Expected: Agent-1 got: ", agent)
	}

	agent = agents.next()
	if agent != "Agent-2" {
		t.Error("Expected: Agent-2 got: ", agent)
	}

	agent = agents.next()
	if agent != "Agent-3" {
		t.Error("Expected: Agent-3 got: ", agent)
	}

	agent = agents.next()
	if agent != "Agent-1" {
		t.Error("Expected: Agent-1 got: ", agent)
	}
}

func TestAgentsRace(t *testing.T) {
	fileName := "tmp.yml"
	f, err := ioutil.TempFile("", fileName)
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())

	cfg := `
    User-Agents:
        - "Agent-1"
        - "Agent-2"
        - "Agent-3"
    `
	ioutil.WriteFile(f.Name(), []byte(cfg), 0644)

	agents := newHttpAgents(f.Name())

	wg := sync.WaitGroup{}
	wg.Add(10)
	defer wg.Wait()

	for i := 0; i < 10; i++ {
		go func(w *sync.WaitGroup, agents *HttpAgents) {
			defer w.Done()
			agents.next()
		}(&wg, agents)
	}

}
