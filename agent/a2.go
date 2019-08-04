package agent

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/hwhaocool/peas/logger"
)

func Run(serverHost string, serverPort string) {
	dest := serverHost + ":" + serverPort

	logger.LogInfo("Connecting to %s...\n", dest)

	conn, err := net.Dial("tcp", dest)

	if err != nil {
		if _, t := err.(*net.OpError); t {
			logger.LogError("Some problem connecting.")
		} else {
			logger.LogError("connection error :%s", err.Error())
		}
		os.Exit(1)
	}

	go readConnection(conn)

	go keepHealth(conn)
}

func keepHealth(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')

		conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error writing to stream.")
			break
		}
	}
}

func readConnection(conn net.Conn) {
	for {
		msg, err := readAll(conn)

		if err != nil {
			continue
		}

		logger.LogInfo("head buf is %b", msg)
		logger.LogInfo("head buf is %s", string(msg))
	}
}

func readAll(conn net.Conn) ([]byte, error) {

	buf := make([]byte, 0, 4096)
	len := 0

	for {
		n, err := conn.Read(buf[len:])
		if n > 0 {
			len += n
		}
		if err != nil {
			if err != io.EOF {
				//Error Handler
				logger.LogError("read error :%s", err.Error())
			}

			return nil, err
		}
	}

	return buf[:len], nil
}

func receiveTask() {

}

func handleCommands(text string) bool {
	r, err := regexp.Compile("^%.*%$")
	if err == nil {
		if r.MatchString(text) {

			switch {
			case text == "%quit%":
				fmt.Println("\b\bServer is leaving. Hanging up.")
				os.Exit(0)
			}

			return true
		}
	}

	return false
}
