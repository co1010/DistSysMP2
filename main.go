package main

import (
	"MP2/utils"
	"github.com/akamensky/argparse"
	"os"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("User", "Saves server and user information")
	// Create string list flag
	var serverInfo = parser.StringList("s", "string", &argparse.Options{Required: true, Help: "Input server ip, server port, and your username"})
	// Parse input
	err := parser.Parse(os.Args)
	utils.CheckError(err)

	// Look for server keyword- start server if it's there
	if (*serverInfo)[0] == "server" {
		startServer((*serverInfo)[1])
	}

	// Wrap command line information into neat User struct
	user := utils.User{IP: (*serverInfo)[0], PORT: (*serverInfo)[1], Username: (*serverInfo)[2]}

	// Connect to the sever
	conn := connectToServer(user)

	// Register with the server
	registerWithServer(conn, user.Username)

	// Goroutine to read incoming messages from server
	go incomingMessages(conn)

	// Read client command line input
	readClientCommands(conn, user.Username)
}
