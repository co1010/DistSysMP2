package main

import (
	"MP2/utils"
	"encoding/gob"
	"net"
)

func startServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	utils.CheckError(err)

	connections := make(map[string]*net.Conn)

	for {
		conn, err := ln.Accept()
		utils.CheckError(err)

		go handleConnection(conn, &connections)
	}
}

func handleConnection(conn net.Conn, connections *map[string]*net.Conn) {
	for {
		var message utils.Message

		decoder := gob.NewDecoder(conn)
		decoder.Decode(&message)

		if message.Register {
			(*connections)[message.From] = &conn
		} else {
			to := (*connections)[message.To]
			utils.SendMessage(*to, message)
		}
	}
}
