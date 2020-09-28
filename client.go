package main

import (
	"MP2/utils"
	"encoding/gob"
	"net"
	"os"
	"strings"
)

// Read incoming messages from given connection
func incomingMessages(conn net.Conn) {
	for {
		// Save decoded message in message variable
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		decoder.Decode(&message)

		// If the exit field is true, the server has shut down.
		if message.Exit {
			println("Server has shut down, exiting...")
			// Close connection and exit
			err := conn.Close()
			utils.CheckError(err)
			os.Exit(0)
		} else {
			utils.PrintMessage(message)
		}
	}
}

// Connect to the server given by the User struct, return the established connection
func connectToServer(user utils.User) net.Conn {
	// Set address as ip:port
	address := user.IP + ":" + user.PORT

	// Attempt to connect to IP:PORT of destination
	conn, err := net.Dial("tcp", address)
	utils.CheckError(err)

	return conn
}

// Register a username with a connection
func registerWithServer(conn net.Conn, username string) {
	// Use the register field of the message so the server knows to register the username
	message := utils.Message{From: username, Register: true}
	utils.SendMessage(conn, message)
}

// Read client input from the command line
func readClientCommands(conn net.Conn, username string) {
	ch := make(chan string)
	go utils.ReadCommands(ch)

	for {
		text := <-ch
		// Commands are separated by spaces, split the text into a slice on spaces
		input := strings.Split(text, " ")

		if input[0] == "send" {
			to := input[1]
			from := username
			content := strings.Join(input[2:], " ")

			// Send message to the given user with the given content. Register and Exit flags are false
			message := utils.Message{To: to, From: from, Content: content, Register: false, Exit: false}
			go utils.SendMessage(conn, message)

		} else if input[0] == "EXIT" {
			// Send message with exit flag to the server
			println("Exit command received. Exiting...")
			message := utils.Message{From: username, Exit: true}
			utils.SendMessage(conn, message)
			// Close the connection and exit
			err := conn.Close()
			utils.CheckError(err)
			os.Exit(0)
		}
	}
}
