package nwmodel

import (
	"errors"
	"fmt"
	"log"
	"nwmessage"
	"strconv"
)

var mapCmdList = map[string]playerCommand{
	"nm": cmdNewBlankMap,

	"nrm": cmdNewRandMap,

	"an": cmdAddNodes,

	"rn": cmdRemoveNodes,

	"ln": cmdLinkNodes,

	//remoteness?????

	// ring

	// line

	// flower

	// "ap": cmdAddPoes,

	// "bake": cmdBakeSlots,

	// slot node_num, [criteria] what (takes challenge of this but not that, etc...)
}

func cmdAddNodes(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	var gnodeCount int
	if len(args) == 0 {
		gnodeCount = 1
	} else {
		nodeCount, err := strconv.Atoi(args[0])
		if err != nil {
			return nwmessage.PsError(err)
		}
		gnodeCount = nodeCount
	}

	gm.Map.addNodes(gnodeCount)

	gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("%d new nodes created", gnodeCount))
}

func cmdRemoveNodes(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	if len(args) == 0 {
		return nwmessage.PsError(errors.New("Need at least one node to remove"))
	}

	toRemove := make([]int, len(args))
	for i := range toRemove {
		targetNode, err := strconv.Atoi(args[i])
		if err != nil {
			return nwmessage.PsError(err)
		}

		if !gm.Map.nodeExists(targetNode) {
			return nwmessage.PsError(fmt.Errorf("Node %d not found", targetNode))
		}

		toRemove[i] = targetNode
	}

	gm.Map.removeNodes(toRemove)

	gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess(fmt.Sprintf("Removed nodes, ", toRemove))
}

func cmdLinkNodes(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	if len(args) != 2 {
		return nwmessage.PsError(errors.New("Need exactly two nodes to link"))
	}

	targ := make([]int, 2)
	for i := range targ {
		log.Printf("i in loop: %d", i)
		targetNode, err := strconv.Atoi(args[i])
		log.Printf("targetNode: %d", targetNode)
		if err != nil {
			return nwmessage.PsError(err)
		}
		if !gm.Map.nodeExists(targetNode) {
			return nwmessage.PsError(fmt.Errorf("Node %d does not exist", targetNode))
		}

		targ[i] = targetNode
	}

	log.Printf("linking node %d to node %d", targ[0], targ[1])
	gm.Map.Nodes[targ[0]].addConnection(gm.Map.Nodes[targ[1]])
	// toLink := make([]*node, len(args))
	// for i := range toLink {
	// 	targetNode, err := strconv.Atoi(args[i])
	// 	if err != nil {
	// 		return nwmessage.PsError(err)
	// 	}

	// 	if !gm.Map.nodeExists(targetNode) {
	// 		return nwmessage.PsError(fmt.Errorf("Node %d not found", targetNode))
	// 	}
	// 	toLink[i] = gm.Map.Nodes[targetNode]
	// }
	// log.Printf("toLink: %v", toLink)

	// for _, node := range toLink {
	// 	for otherNode := range toLink {
	// 		// log.Printf("linking node %d to node %d", node.ID, otherNode)
	// 		log.Printf("linking node %d to node %d", node.ID, gm.Map.Nodes[otherNode].ID)
	// 		node.addConnection(gm.Map.Nodes[otherNode])
	// 	}
	// }
	gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess("")
}

func cmdNewBlankMap(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {

	if len(args) != 0 {
		return nwmessage.PsError(errors.New("Command does not accept arguments"))
	}

	// nodeIdcount should be irrelevant since its now tied to maps
	// nodeIDCount = 0
	newBlankMap := newNodeMap()
	gm.Map = &newBlankMap
	gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess("Generating new blank map...")
}

func cmdNewRandMap(p *Player, gm *GameModel, args []string, c string) nwmessage.Message {
	// TODO fix d3 to update...
	for _, t := range gm.Teams {
		if t.poe != nil {
			return nwmessage.PsError(errors.New("Cannot alter map after a Point of Entry is set"))
		}
	}
	nodeCount, err := validateOneIntArg(args)
	if err != nil {
		return nwmessage.PsError(err)
	}

	// nodeIdcount should be irrelevant since its now tied to maps
	// nodeIDCount = 0
	gm.Map = newRandMap(nodeCount)
	gm.broadcastGraphReset()
	gm.psBroadcastExcept(p, nwmessage.PsAlert("Map was reset"))
	gm.broadcastState()
	return nwmessage.PsSuccess("Generating new random map...")
}
