package nwmodel

import (
	"feature"
	"fmt"
	"math/rand"
	"sync"

	"nwmessage"
)

type machine struct {
	sync.RWMutex
	// accepts   challengeCriteria // store what challenges can fill this machine
	challenge Challenge

	Powered  bool   `json:"powered"`
	builder  string // `json:"creator"`
	TeamName string `json:"owner"`

	attachedPlayers map[*Player]bool

	address string // mac address in node where machine resides

	// solution string // store solution used to pass. could be useful for later mechanics
	Type feature.Type `json:"type"` // NA for non-features, none or other feature.Type for features

	language  string // `json:"languageId"`
	Health    int    `json:"health"`
	MaxHealth int    `json:"maxHealth"`
}

type challengeCriteria struct {
	IDs        []int64  // list of acceptable challenge ids
	Tags       []string // acceptable categories of challenge
	Difficulty [][]int  // acceptable difficulties, [5] = level five, [3,5] = 3,4, or 5
}

// init methods

func newMachine() *machine {
	return &machine{
		Powered:         true,
		attachedPlayers: make(map[*Player]bool),
	}
}

func newFeature() *machine {
	m := newMachine()
	m.Type = feature.None
	return m
}

// machine methods -------------------------------------------------------------------------

// QUESTION: is technically reading state but an awkward place to employ a read lock
func (m *machine) detachAll(msg string) {
	for p := range m.attachedPlayers {
		p.macDetach()
		if msg != "" {
			p.Outgoing <- nwmessage.PsAlert(msg)
		}
	}
}

// Getters (need RLock)

func (m *machine) isNeutral() bool {
	m.RLock()
	if m.TeamName == "" {
		return true
	}
	m.RUnlock()
	return false
}

func (m *machine) isFeature() bool {
	m.RLock()
	if m.Type == nil {
		return false
	}
	m.RUnlock()
	return true
}

func (m *machine) belongsTo(teamName string) bool {
	m.RLock()
	if m.TeamName == teamName {
		return true
	}
	m.RUnlock()
	return false
}

// Setters (require Locks) -----------------------------------------------

func (m *machine) addPlayer(p *Player) {
	m.Lock()
	m.attachedPlayers[p] = true
	m.Unlock()
}

func (m *machine) remPlayer(p *Player) {
	m.Lock()
	delete(m.attachedPlayers, p)
	m.Unlock()
}

func (m *machine) reset() {
	m.Lock()
	m.builder = ""
	m.TeamName = ""
	m.language = ""
	m.Powered = true

	m.detachAll(fmt.Sprintf("mac:%s is resetting, you have been detached", m.address))

	// if m.Type != nil { // reset feature type?
	// 	m.Type = feature.None
	// }

	m.Health = 0
	m.resetChallenge()
	m.Unlock()
}

// resetChallenge should use m.accepts to get a challenge matching criteria TODO
func (m *machine) resetChallenge() {
	m.Lock()
	m.challenge = getRandomChallenge()
	m.MaxHealth = len(m.challenge.Cases)
	m.Unlock()
}

func (m *machine) claim(p *Player, r GradedResult) {
	m.Lock()
	m.builder = p.name
	m.TeamName = p.TeamName
	m.language = p.language
	// m.Powered = true

	m.Health = r.passed()
	m.Unlock()
}

// dummyClaim is used to claim a machine for a player without an execution result
func (m *machine) dummyClaim(teamName string, str string) {
	m.Lock()
	// m.builder = p.name
	m.TeamName = teamName
	m.language = "python" // TODO find ore elegent solution for this
	// m.Powered = true

	switch str {
	case "FULL":
		m.Health = m.MaxHealth
	case "RAND":
		m.Health = rand.Intn(m.MaxHealth) + 1
	case "MIN":
		m.Health = 1
	}
	m.Unlock()
}
