package main

import (
	"fmt"
	"net"
)

const serverIP = "147.75.85.99"
const localIP = "127.0.0.1"

// we set this to 127.0.0.1 in /etc/hosts
const serverHost = "maze.liveoverflow.com"
const minPort = 1337
const maxPort = 1357

// proxy http requests to maze.liveoverflow.com
// set hostname and port to something on this machine
// send all requests to that to maze.liveoverflow.com

// MessageHandler will receive all message payloads.
type MessageHandler func(origin string, msg []byte) [][]byte

type udpProxy struct {
	clientConn *net.UDPConn
	serverConn *net.UDPConn
	clientChan chan []byte
	serverChan chan []byte
	clientAddr *net.UDPAddr
	msgHandler MessageHandler
	activePort int
}

func newUDPProxy() {
	up := &udpProxy{
		serverChan: make(chan []byte),
		clientChan: make(chan []byte),
		msgHandler: msgInterceptor,
	}
	for port := minPort; port <= maxPort; port++ {
		maybeConn, err := net.ListenUDP("udp", addr(localIP, port))
		check(err)
		go func(port int) {
			buffer := make([]byte, 4096)
			// fmt.Println("Listening on", port)
			n, clientAddr, err := maybeConn.ReadFromUDP(buffer)
			check(err)
			if up.clientAddr != nil {
				fmt.Printf("== WARN == Switched to port %d\n", port)
			}
			msg := buffer[:n]
			fmt.Printf("STARTED COMMS ON PORT %d\n", port)
			fmt.Printf("CLIENT>0x%x\n", decode(msg))
			up.clientAddr = clientAddr
			up.clientConn = maybeConn
			up.serverChan <- msg
			up.serverConn, err = net.DialUDP("udp", nil, addr(serverIP, port))
			if err != nil {
				fmt.Println(err)
			}
			check(err)
			up.activePort = port
		}(port)
		// wait for first message
	}
	done := make(chan bool, 1)
	msg := <-up.serverChan
	fmt.Println("Starting UDP proxy")
	go up.recvClient()
	go up.recvServer()
	go up.sendClient()
	go up.sendServer()
	up.serverChan <- msg
	<-done
}

func (up *udpProxy) sendClient() {
	for msg := range up.clientChan {
		_, err := up.clientConn.WriteToUDP(msg, up.clientAddr)
		if err != nil {
			fmt.Println("CLIENT SEND", err)
		}
		// print("CLIENT SEND", msg)
	}
}

func (up *udpProxy) ensureServerConn() {
	if up.serverConn == nil {
		serverConn, err := net.DialUDP("udp", nil, addr(serverIP, up.activePort))
		if err != nil {
			fmt.Println(err)
		}
		up.serverConn = serverConn
	}
}

func (up *udpProxy) sendServer() {
	for msg := range up.serverChan {
		up.ensureServerConn()
		_, _, err := up.serverConn.WriteMsgUDP(msg, nil, nil)
		if err != nil {
			fmt.Println("SERVER SEND", err)
		}
		// print("SERVER SEND", msg)
	}
}

func (up *udpProxy) recvClient() {
	buffer := make([]byte, 4096)
	for {
		n, clientAddr, err := up.clientConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		up.clientAddr = clientAddr
		msg := make([]byte, n)
		copy(msg, buffer[:n])
		msgs := up.msgHandler("CLIENT", msg)
		for _, out := range msgs {
			if len(out) > 0 {
				up.serverChan <- out
			}
		}
	}
}

func (up *udpProxy) recvServer() error {
	buffer := make([]byte, 4096)
	for {
		up.ensureServerConn()
		n, _, err := up.serverConn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		msg := make([]byte, n)
		copy(msg, buffer[:n])
		msgs := up.msgHandler("SERVER", msg)
		for _, out := range msgs {
			if len(out) > 0 {
				up.clientChan <- out
			}
		}
	}
}

func addr(ip string, port int) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	check(err)
	return addr
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
