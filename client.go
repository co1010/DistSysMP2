package main

import (
	"MP2/utils"
	"encoding/gob"
	"net"
)

func connectToServer(user utils.User) net.Conn {
	// Set address as ip:port
	address := user.IP + ":" + user.PORT

	// Attempt to connect to IP:PORT of destination
	conn, err := net.Dial("tcp", address)
	utils.CheckError(err)

	return conn
}

func registerWithServer(conn net.Conn, username string) {
	message := utils.Message{From: username, Register: true}
	utils.SendMessage(conn, message)
}

func incomingMessages(conn net.Conn) {
	for {
		var message utils.Message
		decoder := gob.NewDecoder(conn)
		decoder.Decode(&message)
		utils.PrintMessage(message)
	}
}
