package nwmodel

// Stringers ----------------------------------------------------------------------------------
import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func (n node) String() string {
	return fmt.Sprintf("( <node> {ID: %v, Connections:%v, Machines:%v} )", n.ID, n.Connections, n.Machines)
}

// func (n node) modIDs() []modID {
// 	ids := make([]modID, 0)
// 	for _, slot := range n.Machines {
// 		if slot.TeamName != "" {
// 			ids = append(ids, slot.Module.id)
// 		}
// 	}
// 	return ids
// }

func (t team) String() string {
	var playerList []string
	for player := range t.players {
		playerList = append(playerList, string(player.GetName()))
	}
	return fmt.Sprintf("( <team> {Name: %v, Players:%v} )", t.Name, playerList)
}

func (p Player) String() string {
	return fmt.Sprintf("( <player> {Name: %v, team: %v} )", p.GetName(), p.TeamName)
}

func (r route) String() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}

	return fmt.Sprintf("( <route> {Endpoint: %v, Through: %v} )", r.Endpoint().ID, strings.Join(nodeList, ", "))
}

func (n node) StringFor(p *Player) string {

	// sort keys for consistent presentation
	addList := make([]string, 0)
	for add := range n.addressMap {
		addList = append(addList, add)
	}
	sort.Strings(addList)

	// compose list of all machines
	macList := ""
	for _, add := range addList {
		atIndicator := ""
		if p.macAddress == add {
			atIndicator = "*"
		}
		mac := n.addressMap[add]
		macList += "\n" + add + ":" + mac.StringFor(p) + atIndicator
	}

	connectList := strings.Trim(strings.Join(strings.Split(fmt.Sprint(n.Connections), " "), ","), "[]")

	return fmt.Sprintf("NodeID: %v\nConnects To: %s\nMachines: %v", n.ID, connectList, macList)
}

func (m machine) StringFor(p *Player) string {
	var feature string

	if m.isFeature() {
		feature = " (feature)"
	}

	switch {
	case m.TeamName != "":
		return "(" + m.details() + ")" + feature
	default:
		return "( -neutral- )" + feature
	}
}

func (m machine) details() string {
	return fmt.Sprintf("[%s] [%s] [%s] [%d/%d]", m.TeamName, m.builder, m.language, m.Health, m.MaxHealth)
}

func (r route) forMsg() string {
	nodeCount := len(r.Nodes)
	nodeList := make([]string, nodeCount)

	for i, node := range r.Nodes {
		// this loop is a little funny because we are reversing the order of the node list
		// it's reverse ordered in the data structure but to be human readable we'd like
		// the list to read from source to target
		nodeList[nodeCount-i-1] = strconv.Itoa(node.ID)
	}
	return fmt.Sprintf("(Endpoint: %v, Through: %v)", r.Endpoint().ID, strings.Join(nodeList, ", "))
}

// func (c CompileResult) String() string {
// 	ret := ""
// 	for k, v := range c.Graded {
// 		ret += fmt.Sprintf("(in: %v, out: %v)", k, v)
// 	}
// 	return ret
// }

// func (m machine) String() string {

// }
