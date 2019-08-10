package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hwhaocool/peas/logger"
)

// TCPNetwork manages all server and client connections
type TCPNetwork struct {
	listener        net.Listener
	connIdForServer int
	connIdForClient int
	connsForServer  map[int]*net.Conn
	connsForClient  map[int]*net.Conn
	shutdownFlag    int32
	readTimeoutSec  int
}

// NewTCPNetwork creates a TCPNetwork object
func NewTCPNetwork() *TCPNetwork {
	s := &TCPNetwork{}
	s.connsForServer = make(map[int]*net.Conn)
	s.connsForClient = make(map[int]*net.Conn)
	s.shutdownFlag = 0
	//	default config
	return s
}

// func (t *TCPNetwork) RunApi() error {

// }

// Listen an address to accept client connection
func (t *TCPNetwork) Listen(port string) error {
	addr := "0.0.0.0:" + port
	ls, err := net.Listen("tcp", addr)
	if nil != err {

		logger.LogError("listen error :%s", err.Error())
		return err
	}

	//	accept
	t.listener = ls
	go t.acceptRoutine()
	return nil
}

// Connect the remote server
func (t *TCPNetwork) Connect(addr string) (*Connection, error) {
	conn, err := net.Dial("tcp", addr)
	if nil != err {
		return nil, err
	}

	connection := t.createConn(conn)
	connection.from = 1
	connection.run()
	connection.init()

	return connection, nil
}

func (t *TCPNetwork) addConn(conn net.Conn) {
	connID := 0
	t.connsForClient[connID] = &conn
}

//Run run a server
func Run(port string) {

	fmt.Println("Starting server...")

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)
	fmt.Printf("Listening on %s.\n", src)

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		handleMessage(scanner.Text(), conn)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println("> " + message)

	if len(message) > 0 && message[0] == '/' {
		switch {
		case message == "/time":
			resp := "It is " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

		case message == "/quit":
			fmt.Println("Quitting.")
			conn.Write([]byte("I'm shutting down now.\n"))
			fmt.Println("< " + "%quit%")
			conn.Write([]byte("%quit%\n"))
			os.Exit(0)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}

func (t *TCPNetwork) acceptRoutine() {
	// After accept temporary failure, enter sleep and try again
	var tempDelay time.Duration

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			// Check if the error is an temporary error
			if acceptErr, ok := err.(net.Error); ok && acceptErr.Temporary() {
				if 0 == tempDelay {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Println("Accept error %s , retry after %d ms", acceptErr.Error(), tempDelay)

				time.Sleep(tempDelay)
				continue
			}

			logger.LogError("accept routine quit.error:%s", err.Error())
			t.listener = nil
			return
		}

		//send a message
		conn.Write([]byte("ping"))

		// Process conn event
		t.addConn(conn)
	}
}
