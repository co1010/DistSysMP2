package utils

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type Message struct {
	To       string
	From     string
	Content  string
	Register bool
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

func SendMessage(conn net.Conn, message Message) {
	encoder := gob.NewEncoder(conn)
	encoder.Encode(message)
}

func PrintMessage(message Message) {
	fmt.Printf("%s: %s\n", message.From, message.Content)
}
