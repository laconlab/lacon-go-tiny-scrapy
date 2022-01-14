// round robin agent list
package crawler

import (
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

type HttpAgents struct {
	id     int
	Agents []string `yaml:"User-Agents"`
	mu     sync.Mutex
}

func newHttpAgents(filePath string) *HttpAgents {
	cfg, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	agents := &HttpAgents{}
	if err := yaml.Unmarshal(cfg, &agents); err != nil {
		log.Fatal(err)
	}

	return agents
}

func (a *HttpAgents) next() string {
	a.mu.Lock()
	defer a.mu.Unlock() // overkill

	ret := a.Agents[a.id]
	a.id = (a.id + 1) % len(a.Agents)
	return ret
}
