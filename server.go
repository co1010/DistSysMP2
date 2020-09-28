package main

import (
	"MP2/utils"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

// Connections is global because almost every server function uses it
var connections map[string]net.Conn

// Starts a TCP server on the given port
func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	utils.CheckError(err)

	// Initialize connections as a map with string keys and net.Conn pointers as the values
	connections = make(map[string]net.Conn)

	// Goroutine to accept incoming connections
	go acceptConnections(ln)

	// Read server command line input
	readServerCommands(ln)
}

// Accepts incoming connections on the given listener
func acceptConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		// Each connection has its own goroutine
		go handleConnection(conn)
	}
}

// Handles the connection
func handleConnection(conn net.Conn) {
	for {
		// Decode connection and store in message
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		decoder.Decode(&message)

		// Use Message's exit and register fields to determine what to do
		if message.Exit {
			// Delete the user from connections, close the connection, and break the loop.
			delete(connections, message.From)
			conn.Close()
			fmt.Printf("User %s disconnected\n", message.From)
			break
		} else if message.Register {
			// Register the user in the connections map
			connections[message.From] = conn
			fmt.Printf("User %s connected\n", message.From)
		} else {
			// User connections map to find out who the message is going to
			to := connections[message.To]
			// Error handling if the user doesn't exist
			if to == nil {
				errorMessage := utils.Message{From: "Error", Content: "That user is not connected to this server"}
				utils.SendMessage((connections)[message.From], errorMessage)
			} else {
				utils.SendMessage(to, message)
			}
		}
	}
}

// Reads command line input
func readServerCommands(ln net.Listener) {
	commands := make(chan string)
	go utils.ReadCommands(commands)
	for {
		command := <-commands
		// If the user enters exit
		if command == "EXIT" {
			// Close all active connections
			terminateConnections()
			// Close the server
			err := ln.Close()
			utils.CheckError(err)
			os.Exit(0)
		}
	}
}

// Closes all active connections
func terminateConnections() {
	for _, connection := range connections {
		// Error handling
		if connection != nil {
			// Send message to clients with the exit flag true
			message := utils.Message{Exit: true}
			utils.SendMessage(connection, message)
		}
	}
}
