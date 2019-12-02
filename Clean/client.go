package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var myhost = "localhost:3000"

type request struct {
	Code string
	Key  int
	Name string
}

func main() {
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Insert your key: ")
	strKey, _ := gin.ReadString('\n')
	key, _ := strconv.Atoi(strings.TrimSpace(strKey))
	remotehost := "localhost:8003"

	msg := request{"ACCESS", key, myhost}
	go send(msg, remotehost)
	server()
}

func server() {
	ln, _ := net.Listen("tcp", myhost)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	var response request
	dec := json.NewDecoder(conn)
	dec.Decode(&response)

	fmt.Println(response.Code)
}

func send(msg request, remotehost string) {
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()

	enc := json.NewEncoder(conn)
	enc.Encode(msg)
}
