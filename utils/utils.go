package utils

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	To       string
	From     string
	Content  string
	Register bool
	Exit     bool
}

type User struct {
	IP       string
	PORT     string
	Username string
}

// Consolidated repeated error checks into a single function
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}

// Consolidated repeated message sends into a single function
func SendMessage(conn net.Conn, message Message) {
	encoder := gob.NewEncoder(conn)
	encoder.Encode(message)
}

// Prints the message
func PrintMessage(message Message) {
	fmt.Printf("%s: %s\n", message.From, message.Content)
}

// Reads input and sends each line with stripped \r and \n to the channel
func ReadCommands(ch chan<- string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		CheckError(err)
		text = strings.TrimRight(text, "\r\n")
		ch <- text
	}
}
