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

var addrs []int
var myhost string
var myKeys []int
var spaceShip string
var cont int
var found bool

type request struct {
	Code string
	Key  int
	Name string
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Insert your port: ")
	port, _ := gin.ReadString('\n')
	hostname, _ := strconv.Atoi(strings.TrimSpace(port))
	myhost = fmt.Sprintf("localhost:%d", hostname+3)

	go registerServer(hostname)
	go hotServer(hostname)

	for i := 0; i < 10; i++ {
		myKeys = append(myKeys, rand.Intn(100)+1)
	}
	fmt.Println(myKeys)

	fmt.Print("Insert remote port: ")
	port, _ = gin.ReadString('\n')
	port = strings.TrimSpace(port)

	if port != "" {
		remoteHost, _ := strconv.Atoi(port)
		registerSend(hostname, remoteHost)
	}

	notifyServer(hostname)
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
	dec.Decode(&msg)

	switch msg.Code {
	case "ACCESS":
		spaceShip = msg.Name
		cont = 0
		found = false
		for _, num := range myKeys {
			if num == msg.Key {
				fmt.Println("Key found!!")
				response := request{"Welcome", num, myhost}
				hotResponse(response, spaceShip)
				return
			}
		}
		fmt.Println("Key not found, checking on other nodes...")
		response := request{"SEARCH", msg.Key, myhost}
		hotSendAll(response)
	case "SEARCH":
		for _, num := range myKeys {
			if num == msg.Key {
				fmt.Println("Key found!!")
				response := request{"FOUND", num, myhost}
				hotResponse(response, msg.Name)
				return
			}
		}
		fmt.Println("Key not found")
		response := request{"NOTFOUND", msg.Key, myhost}
		hotResponse(response, msg.Name)
	case "FOUND":
		if !found {
			fmt.Println("Key found on other node !")
			response := request{"Welcome", msg.Key, myhost}
			hotResponse(response, spaceShip)
		}

	case "NOTFOUND":
		cont++
		if cont == len(addrs) {
			response := request{"Bye", msg.Key, myhost}
			hotResponse(response, spaceShip)
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

func hotResponse(msg request, remote string) {
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	enc := json.NewEncoder(conn)
	enc.Encode(msg)
}
