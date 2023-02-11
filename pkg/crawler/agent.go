// round robin agent list
package crawler

type HttpAgents struct {
	Agents []string `yaml:"userAgents"`
	id     int
}

func (a *HttpAgents) Next() string {
	if len(a.Agents) == 0 {
		return ""
	}
	a.id = (a.id + 1) % len(a.Agents)
	return a.Agents[a.id]
}
