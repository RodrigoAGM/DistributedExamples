package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const localAddr = "localhost:8003"
const (
	cnum = iota
	opa
	opb
)

type tmsg struct {
	Code int
	Addr string
	Op   int
}

var addrsC = []string{"localhost:8001", "localhost:8000", "localhost:8002"}

var chInfo chan map[string]int

func main() {
	chInfo = make(chan map[string]int)
	go func() {
		chInfo <- map[string]int{}
	}()
	go serverC()

	time.Sleep(time.Millisecond * 100)
	var op int
	for {
		fmt.Print("Your option: ")
		fmt.Scanf("%d\n", &op)
		msg := tmsg{cnum, localAddr, op}
		for _, addr := range addrsC {
			sendC(addr, msg)
		}
	}

}

func serverC() {
	ln, _ := net.Listen("tcp", localAddr)
	defer ln.Close()

	for {
		conn, _ := ln.Accept()
		go handleC(conn)
	}
}

func handleC(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg tmsg

	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from ", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case cnum:
			consensusC(conn, msg)
		}
	}
}

func consensusC(conn net.Conn, msg tmsg) {
	info := <-chInfo
	info[msg.Addr] = msg.Op

	if len(info) == len(addrsC) {
		ca, cb := 0, 0
		for _, op := range info {
			if op == opa {
				ca++
			} else {
				cb++
			}
		}
		if ca > cb {
			fmt.Println("GO A !")
		} else {
			fmt.Println("GO B !")
		}
		info = map[string]int{}
	}
	go func() { chInfo <- info }()
}

func sendC(remote string, msg tmsg) {
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(msg)
}
