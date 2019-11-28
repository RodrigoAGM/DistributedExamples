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
var n, min int
var contCh chan int

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

	fmt.Print("Enter N: ")
	port, _ = gin.ReadString('\n')
	port = strings.TrimSpace(port)
	n, _ = strconv.Atoi(strings.TrimSpace(port))
	contCh = make(chan int, 1)
	contCh <- 0

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
	num, _ := strconv.Atoi(strings.TrimSpace(str))
	fmt.Printf("Nos ha llegado el %d\n", num)

	cont := <-contCh
	if cont == 0 {
		min = num
	} else if num < min {
		send(min)
		min = num
	} else {
		send(num)
	}
	cont++
	if cont == n {
		fmt.Printf("NUMERO FINAL: %d!!\n", min)
		cont = 0
	}

	contCh <- cont
}

func send(num int) {
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()
	fmt.Fprintf(conn, "%d\n", num)
}
