package crawler

type HttpAgents struct {
	id     int
	Agents []string `yaml:"User-Agents"`
}

func (a *HttpAgents) Next() string {
	ret := a.Agents[a.id]
	a.id = (a.id + 1) % len(a.Agents)
	return ret
}
