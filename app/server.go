package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go connectionHandler(conn)
	}
}

func connectionHandler(conn net.Conn) {
	buf := make([]byte, 256)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			return
		}
		_, resp := ReadCommand(buf)
		fmt.Println("type of command ", string(resp.Type))
		fmt.Println("command raw ", string(resp.Raw))

		switch resp.Type {
		case String:
			switch string(resp.Data) {
			case "PONG":
				conn.Write([]byte("+PONG\r\n"))
			}
		case Array:
			var i int
			var command []string
			temp_resp := resp.Data
			for ; i < resp.Count; i++ {
				rl, rresp := ReadCommand(temp_resp)
				command = append(command, string(rresp.Data))
				temp_resp = temp_resp[rl:]
			}
			switch command[0] {
			case "ECHO":
				{
					send_back := command[1]
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(send_back), send_back)
					conn.Write([]byte(response))
					break
				}
			case "PING":
				{
					conn.Write([]byte("+PONG\r\n"))
				}
			}
		default:

		}
	}
}
