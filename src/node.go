package main

import "strconv"

var NodeID int = 0

type Node struct {
	ID          string  `json:"id"`
	Connections []*Node `json:"connections"`
	Owner       *Team   `json:"owner"`
	Contested   bool    `json:"contested"`
	Ice         []Ice   `json:"ice"`
}

type Ice struct {
	TestID string `json:"test_id"`
	Owner  Team   `json:"owner"`
}

func NewNode(name string) *Node {
	id := name
	if name == "" {
		NodeID++
		id = strconv.Itoa(NodeID)
	}

	connections := make([]*Node, 0)
	ice := make([]Ice, 0)

	return &Node{id, connections, nil, false, ice}
}

func (n *Node) addConnection(m *Node) {
	n.Connections = append(n.Connections, m)
}

// check for contested status and set flag (and owner?)
func (n *Node) isContested() {

}

// returns true if nodes have a valid path
func (n Node) canConnect(m Node) bool {

	return false
}
