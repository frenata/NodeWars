package protocol

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"nwmessage"
	"nwmodel"
	"regrequest"

	"github.com/gorilla/websocket"
)

// Handshake related constants
const versionNumber = "1.0.0"

// VersionTag is used for handshake and generally to identify the server type and version
const VersionTag = "NodeWars:" + versionNumber

// var upgrader = websocket.Upgrader{}
// Allows cross-origin web socket upgrade. Remove for production
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Ask about reduntant error messaging...
func doHandshake(ws *websocket.Conn) error {

	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Printf("Could not read from socket: %v", err)
		return err
	}

	if string(p) == VersionTag {
		message := []byte("Welcome to NodeWars")
		if err := ws.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error while affirming handshake")
			return err
		}
	} else {
		errorMessage := []byte("Version mismation, Server is running " + VersionTag + ", closing connection.")
		if err := ws.WriteMessage(websocket.TextMessage, errorMessage); err != nil {
			log.Println("Error while aborting handshake")
			return err
		}
		return errors.New(string(errorMessage))
	}
	return nil
}

// HandleConnections is the point of entry for all websocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request, d *Dispatcher) {
	fmt.Println("New player connected...")
	// Upgrade GET to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Attempt handshake
	err = doHandshake(ws)
	if err != nil {
		log.Printf("Handshake error: %v", err)
		ws.Close()
		return
	}

	// Assuming we're all good, register client

	// create a single use channel to receive a registered player on
	tempChan := make(chan bool)

	p := nwmodel.NewPlayer(ws)
	d.registrationQueue <- regrequest.Reg(p, tempChan) // add player registration to dispatcher jobs,

	defer func() {

		d.registrationQueue <- regrequest.Dereg(p)

	}() // clean up player when we're done

	// block till player is registered...
	_ = <-tempChan

	// Spin up gorouting to monitor outgoing and send those messages to player.Socket
	// log.Println("Spinning up outgoing handler for player...")

	p.Outgoing(nwmessage.PsPrompt(p.GetName() + "@lobby>"))

	// Handle socket stream
	for {
		msg, err := nwmessage.MsgFromClient(p)

		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		// log.Println("received player message")
		d.clientMessages <- msg
	}
}
