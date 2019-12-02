package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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

var myKeys []int
var addrs []int
var cont int
var myPort int
var remoteAddr string
var found bool

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Insert host port number: \n")
	strPort, _ := gin.ReadString('\n')
	myPort, _ = strconv.Atoi(strings.TrimSpace(strPort))

	go registerServer(myPort)
	go hotServer(myPort)

	for i := 0; i < 10; i++ {
		myKeys = append(myKeys, rand.Intn(99)+1)
	}
	fmt.Println(myKeys)

	gin = bufio.NewReader(os.Stdin)
	fmt.Print("Insert remote host number: \n")
	strRemote, _ := gin.ReadString('\n')

	if strings.TrimSpace(strRemote) != "" {
		remotePort, _ := strconv.Atoi(strings.TrimSpace(strRemote))
		registerSend(myPort, remotePort)
	}

	go func() {
		fmt.Println("Waiting for any access request...")
	}()

	notifyServer(myPort)

}

func notifyServer(port int) {
	serverPort := port + 2
	serverHost := fmt.Sprintf("localhost:%d", serverPort)
	ln, _ := net.Listen("tcp", serverHost)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		handleNotify(conn)
	}
}

func handleNotify(conn net.Conn) {
	defer conn.Close()

	// Recibimos addr del nuevo nodo
	r := bufio.NewReader(conn)
	remoteIP, _ := r.ReadString('\n')
	remoteIP = strings.TrimSpace(remoteIP)
	remoteHost, _ := strconv.Atoi(remoteIP)

	// Agregamos nuevo nodo a la lista de direcciones
	for _, addr := range addrs {
		if addr == remoteHost {
			return
		}
	}
	addrs = append(addrs, remoteHost)
	fmt.Println(addrs)
}

func notifySend(addr, remotePort int) {
	remote := fmt.Sprintf("localhost:%d", addr+2)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	fmt.Fprint(conn, remotePort)
}

func registerServer(port int) {
	serverPort := port + 1
	serverHost := fmt.Sprintf("localhost:%d", serverPort)
	ln, _ := net.Listen("tcp", serverHost)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		handleRegister(conn)
	}
}

func handleRegister(conn net.Conn) {
	defer conn.Close()

	//Recibimos el port del nuevo nodo
	r := bufio.NewReader(conn)
	remoteIP, _ := r.ReadString('\n')
	remoteIP = strings.TrimSpace(remoteIP)
	remotePort, _ := strconv.Atoi(remoteIP)

	// respondemos enviando lista de direcciones de nodos actuales
	byteAddrs, _ := json.Marshal(addrs)
	fmt.Fprintf(conn, "%s\n", string(byteAddrs))

	// notificar a nodos actuales de llegada de nuevo nodo
	for _, addr := range addrs {
		notifySend(addr, remotePort)
	}

	// Agregamos nuevo nodo a la lista de direcciones
	for _, addr := range addrs {
		if addr == remotePort {
			return
		}
	}
	addrs = append(addrs, remotePort)
	fmt.Println(addrs)
}

func registerSend(hostport, remoteport int) {

	remoteHost := fmt.Sprintf("localhost:%d", remoteport+1)
	conn, _ := net.Dial("tcp", remoteHost)
	defer conn.Close()

	// Enviar direccion
	fmt.Fprintln(conn, hostport)

	// Recibir lista de direcciones
	r := bufio.NewReader(conn)
	strAddrs, _ := r.ReadString('\n')
	var respAddrs []int
	json.Unmarshal([]byte(strAddrs), &respAddrs)

	// agregamos direcciones de nodos a propia libreta
	for _, addr := range respAddrs {
		if addr == remoteport {
			return
		}
	}
	addrs = append(respAddrs, remoteport)
	fmt.Println(addrs)
}

func hotServer(port int) {
	serverPort := port + 3
	serverHost := fmt.Sprintf("localhost:%d", serverPort)
	ln, _ := net.Listen("tcp", serverHost)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		handleHot(conn)
	}
}

func handleHot(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(conn)
	var msg request
	name := fmt.Sprintf("localhost:%d", myPort+3)

	if err := dec.Decode(&msg); err != nil {
		log.Println("Can't decode from ", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		key := msg.Key
		switch msg.Code {
		case "ACCESS":
			remoteAddr = msg.Name
			cont = 0
			for _, num := range myKeys {
				if num == key {
					fmt.Println("Key Found! Giving access")
					response := request{"APROVED", num, name}
					hotResponse(msg.Name, response)
					return
				}
			}

			fmt.Println("Key not found! Searching in other nodes...")
			newMsg := request{"SEARCH", msg.Key, name}
			hotSendAll(newMsg)
		case "SEARCH":
			for _, num := range myKeys {
				if num == msg.Key {
					fmt.Println("Key Found!")
					response := request{"FOUND", num, name}
					hotResponse(msg.Name, response)
					return
				}
			}
			fmt.Println("Key not found!")
			response := request{"NOTFOUND", key, name}
			hotResponse(msg.Name, response)

		case "NOTFOUND":
			cont++
			if cont == len(addrs) && !found {
				response := request{"Bye", key, ""}
				hotResponse(remoteAddr, response)
				cont = 0
			}
		case "FOUND":
			if !found {
				fmt.Println("Key Found! Giving access")
				response := request{"Welcome", 0, ""}
				hotResponse(remoteAddr, response)
				found = true
			}
		}
	}
}

func hotSendAll(msg request) {
	for _, addr := range addrs {
		remote := fmt.Sprintf("localhost:%d", addr+3)
		conn, _ := net.Dial("tcp", remote)
		defer conn.Close()
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}

func hotResponse(remote string, msg request) {
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(msg)
}
