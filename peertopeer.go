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
)

var addrs []int

func main() {

	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Insert your port: ")
	port, _ := gin.ReadString('\n')
	hostname, _ := strconv.Atoi(strings.TrimSpace(port))

	go registerServer(hostname)
	go hotServer(hostname)

	fmt.Print("Insert remote port: ")
	port, _ = gin.ReadString('\n')
	port = strings.TrimSpace(port)

	if port != "" {
		remoteHost, _ := strconv.Atoi(port)
		registerSend(hostname, remoteHost)
	}

	go func() {
		fmt.Print("Ingrese num: ")
		strNum, _ := gin.ReadString('\n')
		if strNum != "" {
			num, _ := strconv.Atoi(strings.TrimSpace(strNum))
			hotSend(num)
		}
	}()
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

	r := bufio.NewReader(conn)
	strnum, _ := r.ReadString('\n')
	strnum = strings.TrimSpace(strnum)
	num, _ := strconv.Atoi(strnum)

	if num == 0 {
		fmt.Print("Perdimos :c")
	} else {
		hotSend(num - 1)
	}

}

func hotSend(num int) {
	randi := rand.Intn(len(addrs))
	fmt.Printf("Enviando %d a %d\n", num, addrs[randi])
	remote := fmt.Sprintf("localhost:%d", addrs[randi]+3)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	fmt.Fprint(conn, num)
}
