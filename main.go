package main

import (
	"fmt"
	"log"
	"net"
)

type MessageType int

const (
	ClientConnected MessageType = iota + 1
	DeleteClient
	NewMessage
)

type Message struct {
	Type MessageType
	Conn net.Conn
	Text string
}

func client(conn net.Conn, messages chan Message) {
	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			conn.Close()
			messages <- Message{
				Type: DeleteClient,
				Conn: conn,
			}
			return
		}
		messages <- Message{
			Type: NewMessage,
			Text: string(buffer[0:n]),
			Conn: conn,
		}

	}
}

func server(message chan Message) {
	conns := map[string]net.Conn{}
	for {
		msg := <-message
		switch msg.Type {
		case ClientConnected:
			conns[msg.Conn.RemoteAddr().String()] = msg.Conn
		case DeleteClient:
			delete(conns, msg.Conn.RemoteAddr().String())
		case NewMessage:
			for _, conn := range conns {
				_, err := conn.Write([]byte(msg.Text))
				if err != nil {
					// TODO: remove the connection from the list
					fmt.Println("Could not send data to ...: %s", err)
				}

			}
		}
	}
}

// [null-ls] typechecking error: pattern ./...: directory prefix . does not contain main module or its selected dependencies
func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error ....\n")
	}
	log.Printf("Listening to TCP connection on port %s ...\n", "8080")

	messages := make(chan Message)
	go server(messages)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("ERROR\n")
			// handle error
		}
		log.Printf("accepted connection from %s", conn.RemoteAddr())
		messages <- Message{
			Type: ClientConnected,
			Conn: conn,
		}
		go client(conn, messages)
	}

}
