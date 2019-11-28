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

type Info struct {
	Tipo     string
	NodeNum  int
	NodeAddr int
}

type MyInfo struct {
	cont     int
	first    bool
	nextNum  int
	nextAddr int
}

var chMyInfo chan MyInfo
var readyToStart chan bool

var addrs []int
var myNum int
var hostname int

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	myNum = rand.Intn(int(1e6))
	var n int
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Insert your port: ")
	port, _ := gin.ReadString('\n')
	hostname, _ = strconv.Atoi(strings.TrimSpace(port))

	fmt.Print("Ingrese la cantidad de nodos: ")
	fmt.Scanf("%d\n", &n)
	addrs = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Ingrese nodo %d: ", i+1)
		fmt.Scanf("%d\n", &(addrs[i]))
	}
	readyToStart = make(chan bool)

	go func() {
		chMyInfo = make(chan MyInfo)
		chMyInfo <- MyInfo{0, true, int(1e7), -1}
	}()
	go func() {
		gin := bufio.NewReader(os.Stdin)
		fmt.Print("Presione enter para iniciar...")
		gin.ReadString('\n')
		info := Info{"SENDNUM", myNum, hostname}
		for _, addr := range addrs {
			sendFunc(addr, info)
		}
	}()
	server()
}

func server() {
	host := fmt.Sprintf("localhost:%d", hostname+1)
	ln, _ := net.Listen("tcp", host)

	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')
	var info Info
	json.Unmarshal([]byte(msg), &info)

	switch info.Tipo {
	case "SENDNUM":
		myInfo := <-chMyInfo
		myInfo.cont++
		if info.NodeNum < myNum {
			myInfo.first = false
		} else if info.NodeNum < myInfo.nextNum {
			myInfo.nextNum = info.NodeNum
			myInfo.nextAddr = info.NodeAddr
		}

		go func() {
			chMyInfo <- myInfo
		}()

		if myInfo.cont == len(addrs) {
			if myInfo.first {
				fmt.Println("Soy el primer!! :D")
				criticalSection()
			} else {
				readyToStart <- true
			}
		}

	case "START":
		<-readyToStart
		criticalSection()
	}
}

func sendFunc(remoteAddr int, info Info) {
	remote := fmt.Sprintf("localhost:%d", remoteAddr+1)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	bytesMsg, _ := json.Marshal(info)
	fmt.Fprint(conn, string(bytesMsg))
}

func criticalSection() {
	fmt.Println("Ha llegado mi turno!! :)")
	myInfo := <-chMyInfo
	if myInfo.nextAddr == -1 {
		fmt.Println("I was the last one! :(")
	} else {
		info := Info{Tipo: "START"}
		fmt.Println(myInfo, info)
		sendFunc(myInfo.nextAddr, info)
	}
}
