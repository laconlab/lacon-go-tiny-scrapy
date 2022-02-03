package crawler

import (
	"sync"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestRoundRobinAgents(t *testing.T) {
	cfg := `
    userAgents:
        - "Agent-1"
        - "Agent-2"
        - "Agent-3"
    `

	agents := &HttpAgents{}
	if err := yaml.Unmarshal([]byte(cfg), agents); err != nil {
		t.Error(err)
	}

	iter := agents.getIter()

	agent := iter()
	if agent != "Agent-2" {
		t.Error("Expected: Agent-2 got: ", agent)
	}

	agent = iter()
	if agent != "Agent-3" {
		t.Error("Expected: Agent-3 got: ", agent)
	}

	agent = iter()
	if agent != "Agent-1" {
		t.Error("Expected: Agent-1 got: ", agent)
	}

	agent = iter()
	if agent != "Agent-2" {
		t.Error("Expected: Agent-2 got: ", agent)
	}
}

func TestAgentsRace(t *testing.T) {
	cfg := `
    userAgents:
        - "Agent-1"
        - "Agent-2"
        - "Agent-3"
    `

	agents := &HttpAgents{}
	if err := yaml.Unmarshal([]byte(cfg), agents); err != nil {
		t.Error(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(10)
	defer wg.Wait()

	for i := 0; i < 10; i++ {
		go func(w *sync.WaitGroup, iter func() string) {
			defer w.Done()
			iter()
		}(&wg, agents.getIter())
	}

}
