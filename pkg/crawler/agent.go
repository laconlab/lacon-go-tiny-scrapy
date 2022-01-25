// round robin agent list
package crawler

import (
	"log"
)

type HttpAgents struct {
	Agents []string `yaml:"User-Agents"`
}

func (a *HttpAgents) getIter() func() string {
	if len(a.Agents) == 0 {
		log.Fatal("No agents")
	}

	id := 0
	return func() string {
		id = (id + 1) % len(a.Agents)
		return a.Agents[id]
	}
}
