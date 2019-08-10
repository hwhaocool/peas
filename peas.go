package main

import (
	"fmt"
	"log"
	"sync"

	"os"
	"syscall"

	"os/signal"

	"sync/atomic"

	"github.com/hwhaocool/peas/agent"
	"github.com/hwhaocool/peas/server"
)

var (
	kServerAddress  = "localhost:14444"
	serverConnected int32
	stopFlag        int32
)

func getEnv(name string) (string, bool) {
	value := os.Getenv("string")
	if value == "" {
		return value, false
	}
	return value, true
}

func main() {

	workType, err := getEnv("type")
	if err {
		log.Panic("env type is required, it server or agent")
	}

	port, err := getEnv("port")
	if err {
		log.Panic("server listen tcp port is required")
	}

	if port == "8080" {
		log.Panic("do not set tcp port to 8080, change another one")
	}

	switch {
	case workType == "server":
		serverIns := server.NewTCPNetwork()

		err := serverIns.Listen(port)
		if nil != err {
			log.Println(err)
			return
		}

		err := serverIns.RunApi()
		if nil != err {
			log.Println(err)
			return
		}

	case workType == "agent":
		serverHost, err := getEnv("serverHost")
		if err {
			log.Panic("server host is required, it can be ip or host ")
		}

		agent.Run(serverHost, port)

		
		fmt.Println("Quitting.")
		// conn.Write([]byte("I'm shutting down now.\n"))
		// fmt.Println("< " + "%quit%")
		// conn.Write([]byte("%quit%\n"))
		// os.Exit(0)

	default:
		// conn.Write([]byte("Unrecognized command.\n"))
		fmt.Println("222.")
	}

	stopCh := make(chan struct{})

	// process event
	var wg sync.WaitGroup
	wg.Add(1)
	go routineEchoServer(server, &wg, stopCh)

	wg.Add(1)
	go routineEchoClient(client, &wg, stopCh)

	// input event
	wg.Add(1)
	// go routineInput(&wg, clientConn)
	go routineInput(&wg, clientConn)

	// server.Send()

	// wait
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

MAINLOOP:
	for {
		select {
		case <-sc:
			{
				//	app cancelled by user , do clean up work
				log.Println("Terminating ...")
				break MAINLOOP
			}
		}
	}

	atomic.StoreInt32(&stopFlag, 1)
	log.Println("Press enter to exit")
	close(stopCh)
	wg.Wait()
	clientConn.Close()
	server.Shutdown()
}
