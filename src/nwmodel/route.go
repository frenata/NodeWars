package nwmodel

import (
	"strconv"
)

type route struct {
	Nodes  []*node `json:"nodes"`
	player *Player
}

// route methods --------------------------------------------

func (r route) Endpoint() *node {
	if len(r.Nodes) < 1 {
		return nil
	}

	return r.Nodes[0]
}

func (r route) runsThrough(n *node) bool { // runsThrough does not check the endpoint
	for i := 0; i < len(r.Nodes)-1; i++ { // don't check last node
		node := r.Nodes[i]
		if n == node {
			return true
		}
	}
	return false
}

// func (r route) containsNode(n *node) (int, bool) {
// 	for i := 0; i < len(r.Nodes) - 1; i++ { // don't check last node
// 		node := r.Nodes[i]
// 		if n == node {
// 			return i, true
// 		}
// 	}
// 	return 0, false
// }

func (r route) length() int {
	return len(r.Nodes)
}

// asIds reverses the order of the nodes and stores ids only
func (r route) asIds() []nodeID {
	nodeCount := len(r.Nodes)
	list := make([]nodeID, nodeCount)

	for i := 0; i < nodeCount; i++ {
		list[i] = r.Nodes[nodeCount-1-i].ID
	}

	// fmt.Printf("Route: %v\n,Nodecount: %d\nList: %v\n", r, nodeCount, list)
	return list
}

// traficMap and packet used to create more easily rendered statemessages
type trafficMap struct {
	Traffic map[string][]packet `json:"traffic"`
}

type packet struct {
	Owner     string `json:"owner"`
	Direction string `json:"dir"`
}

func newTrafficMap() *trafficMap {
	return &trafficMap{
		Traffic: make(map[string][]packet),
	}
}

// TODO route is stored reversed?
func (t *trafficMap) addRoute(r *route, color string) {
	// fmt.Printf("asIds test: %v", r.asIds())
	nodeIDs := r.asIds()
	for i := 0; i < len(nodeIDs)-1; i++ {
		n1 := nodeIDs[i]
		n2 := nodeIDs[i+1]

		var dir, edge string

		// if we're connecting to poe, ignore
		if n1 == n2 {
			continue
		}

		if n1 > n2 {
			dir = "down"
			edge = strconv.Itoa(n1) + "e" + strconv.Itoa(n2)
		} else {
			dir = "up"
			edge = strconv.Itoa(n2) + "e" + strconv.Itoa(n1)
		}
		t.appendPacket(packet{color, dir}, edge)
	}
}

func (t *trafficMap) appendPacket(p packet, edge string) {
	_, ok := t.Traffic[edge]
	if !ok {
		t.Traffic[edge] = []packet{p}
	} else {
		t.Traffic[edge] = append(t.Traffic[edge], p)
	}
}
