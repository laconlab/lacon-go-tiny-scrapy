package selector

type websitesStates struct {
	idOfLastUsedWebsite int
	states              []*websiteState
}

func websiteStateFromSettings() *websitesStates {
	return &websitesStates {
		idOfLastUsedWebsite: 0,
		states:              getWebsitesStates(),
	}
}

func getWebsitesStates() []*websiteState {
	websites := newWebsiteConfig()
	states := make([]*websiteState, 0, len(websites))
	for name, setting := range websites {
		state := newWebsiteStateFromSettings(name, setting)
		states = append(states, state)
	}
	return states
}

func (ws *websitesStates) isEndState() bool {
	for _, state := range ws.states {
		if !state.isEndState() {
			return false
		}
	}
	return true
}

func (ws *websitesStates) createNextHTTPRequest() HTTPRequest {
	website := ws.getNextWebsiteState()
	return website.createNextHTTPRequest()
}

func (ws *websitesStates) getNextWebsiteState() *websiteState {
	if ws.isEndState() {
		panic("There is no more sites to download")
	}

	websiteId := ws.getNextWebsiteId()
	return ws.states[websiteId]
}

func (ws *websitesStates) getNextWebsiteId() int {
	for {
		currentId := ws.roundRobinIdUpdate()
		if !ws.isEndStateById(currentId) {
			return currentId
		}
	}
}

func (ws *websitesStates) roundRobinIdUpdate() int {
	ws.idOfLastUsedWebsite++

	if ws.idOfLastUsedWebsite >= len(ws.states) {
		ws.idOfLastUsedWebsite = 0
	}

	return ws.idOfLastUsedWebsite
}

func (ws *websitesStates) isEndStateById(id int) bool {
	return ws.states[id].isEndState()
}
