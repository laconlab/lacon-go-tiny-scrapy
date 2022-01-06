package crawler

import "net/http"

type HttpAgents struct {
	id     int
	Agents []string `yaml:"User-Agents"`
}

func (a *HttpAgents) Next() string {
	ret := a.Agents[a.id]
	a.id = (a.id + 1) % len(a.Agents)
	return ret
}

func AgentHttpRequest(url, agent string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", agent)

	client := &http.Client{}
	return client.Do(req)
}
