package main

import "fmt"

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan string

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan string),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:

			fmt.Printf("Register\n")

			alreadyLoggedIn := false

			//Check to see if the user is already logged in
			for con := range h.connections {

				fmt.Printf("Check %s\n", con.email)

				if c.email == con.email {
					fmt.Printf("Closed")
					c.close("User already logged in")
					alreadyLoggedIn = true
					break
				}

			}

			if (!alreadyLoggedIn) {
				c.sendChat("MOTD AND STUFF")
				h.connections[c] = true
			}

		case c := <-h.unregister:

			fmt.Printf("Disconnected %s\n", c.String())
			delete(h.connections, c)
			close(c.send)
			c.close("Unregistered")

		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}
