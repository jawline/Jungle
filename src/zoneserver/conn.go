package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"strconv"
)

const DisconnectCommand string = "Disconnect"
const ChatCommand string = "Chat"

var lastConnectionId uint64 = 0

func nextConnectionId() string {
	next := strconv.FormatUint(lastConnectionId, 10)
	lastConnectionId++
	return next
}

type loginRequest struct {
	Email    string
	Password string
}

type connection struct {

	//The connection id
	id string

	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string

	// The connections username
	email string
}

type CommandMessage struct {
	Command string
	Data    string
}

/**
 * Read the next message from a connection or return an error
 */
func (c *connection) next() (string, error) {
	var message string

	err := websocket.Message.Receive(c.ws, &message)

	if err != nil {
		return "", err
	}

	return message, nil
}

/**
 * Send the specified data through socket
 */
func (c *connection) sendData(data string) error {
	err := websocket.Message.Send(c.ws, data)
	return err
}

/**
 * Send a chat message
 */
func (c *connection) sendChat(message string) error {

	//Create the new message structure
	chatMessage := CommandMessage{Command: ChatCommand, Data: message}

	//Convert it to json
	messageEncoded, err := json.Marshal(&chatMessage)

	//Check for errors
	if err != nil {
		return err
	}

	return c.sendData(string(messageEncoded))
}

/**
 * Send a close message and then close the connection
 */
func (c *connection) close(reason string) {

	message, _ := json.Marshal(&CommandMessage{DisconnectCommand, reason})

	websocket.Message.Send(c.ws, string(message))
	c.ws.Close()
}

func (c *connection) reader() {

	for {

		message, err := c.next()

		if err != nil {
			break
		}

		structured := CommandMessage{}
		json.Unmarshal([]byte(message), &structured)

		if err != nil {
			break
		}

		if structured.Command == "Chat" {
			h.broadcast <- message
		}
	}

	c.close("Connection issues")
}

func (c *connection) writer() {

	for message := range c.send {

		err := websocket.Message.Send(c.ws, message)

		if err != nil {
			break
		}

	}

	c.close("Connection issues")
}

func (c *connection) String() string {
	return "ID: " + c.id + " email " + c.email
}

func wsHandler(ws *websocket.Conn) {
	c := &connection{id: nextConnectionId(), send: make(chan string, 256), ws: ws}

	login, err := c.next()

	if err != nil {
		c.close("Connection error")
		return
	}

	loginDetails := loginRequest{}

	err = json.Unmarshal([]byte(login), &loginDetails)

	if err != nil {
		c.close("Badly formed login request")
		return
	}

	if loginDetails.Email != "blake_l@imclueless.co.uk" {
		c.close("Invalid Username/Password")
	}

	c.email = loginDetails.Email

	h.register <- c

	fmt.Printf("Registered %s\n", c.String())

	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}
