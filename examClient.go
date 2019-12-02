package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type request struct {
	Code string
	Key  int
	Name string
}

var hostname = "localhost:3000"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Insert access key: \n")
	strKey, _ := gin.ReadString('\n')
	key, _ := strconv.Atoi(strings.TrimSpace(strKey))
	remotehost := "localhost:8003"

	go send(key, remotehost)

	server()
}

func send(key int, remotehost string) {

	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	enc := json.NewEncoder(conn)
	msg := request{"ACCESS", key, hostname}
	fmt.Print(msg)
	enc.Encode(msg)
}

func server() {
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		handle(conn)
	}
}

func handle(conn net.Conn) {

	dec := json.NewDecoder(conn)
	var response request
	dec.Decode(&response)
	fmt.Println(response.Code)
}
