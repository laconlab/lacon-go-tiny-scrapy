package crawler

import (
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

	agent := agents.next()
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

	agent = agents.next()
	if agent != "Agent-2" {
		t.Error("Expected: Agent-2 got: ", agent)
	}
}
