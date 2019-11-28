package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var remoteHost string
var hostname string
var num int

func main() {
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Enter port: ")
	port, _ := gin.ReadString('\n')
	port = strings.TrimSpace(port)
	hostname = fmt.Sprintf("localhost:%s", port)

	fmt.Print("Enter remote port: ")
	port, _ = gin.ReadString('\n')
	port = strings.TrimSpace(port)
	remoteHost = fmt.Sprintf("localhost:%s", port)

	fmt.Print("Enter a number: ")
	port, _ = gin.ReadString('\n')
	port = strings.TrimSpace(port)
	num, _ = strconv.Atoi(strings.TrimSpace(port))

	if hostname == "localhost:8000" {
		send(num)
	}

	server()
}

func server() {
	ln, _ := net.Listen("tcp", hostname)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {

	defer conn.Close()
	r := bufio.NewReader(conn)
	str, _ := r.ReadString('\n')
	n, _ := strconv.Atoi(strings.TrimSpace(str))
	fmt.Printf("Nos ha llegado el %d\n", n)

	if n > 1000 {
		fmt.Printf("Se ha llegado al limite, numero: %d\n", n)
	} else {
		send(n * 2)
	}

}

func send(n int) {
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()
	fmt.Fprintf(conn, "%d\n", n)
}
