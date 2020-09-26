package main

import (
	"MP2/utils"
	"bufio"
	"github.com/akamensky/argparse"
	"os"
	"strings"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("User", "Saves server and user information")
	// Create string list flag
	var serverInfo = parser.StringList("s", "string", &argparse.Options{Required: true, Help: "Input server ip, server port, and your username"})
	// Parse input
	err := parser.Parse(os.Args)
	utils.CheckError(err)

	if (*serverInfo)[0] == "server" {
		startServer((*serverInfo)[1])
	}

	user := utils.User{IP: (*serverInfo)[0], PORT: (*serverInfo)[1], Username: (*serverInfo)[2]}

	conn := connectToServer(user)

	registerWithServer(conn, user.Username)

	go incomingMessages(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		input := strings.Split(text, " ")
		if input[0] == "send" {
			to := input[1]
			from := user.Username
			content := strings.Join(input[2:], " ")
			content = strings.TrimRight(content, "\r\n")
			message := utils.Message{To: to, From: from, Content: content, Register: false}
			go utils.SendMessage(conn, message)
		}
	}
}
