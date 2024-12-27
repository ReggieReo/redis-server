package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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
	config := getFlag()
	m := newRedisMap()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go connectionHandler(conn, m, &config)
	}
}

func connectionHandler(conn net.Conn, m *redisMap, config *map[string]string) {
	buf := make([]byte, 256)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
		_, resp := ReadCommand(buf)

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
			switch strings.ToLower(command[0]) {
			case "echo":

				{
					send_back := command[1]
					conn.Write([]byte(formatBulkString(send_back)))
				}
			case "ping":

				{
					conn.Write([]byte("+PONG\r\n"))
				}
			case "get":

				{
					data := m.get(command[1])
					if data == "" {
						conn.Write([]byte("$-1\r\n"))
					} else {
						conn.Write([]byte(formatBulkString(data)))
					}
				}
			case "set":

				{
					if len(command) == 5 && command[3] == "px" {
						px, _ := strconv.Atoi(command[4])
						m.set(command[1], command[2], px)
					} else {
						m.set(command[1], command[2], 0)
					}
					conn.Write([]byte("+OK\r\n"))
				}
			case "config":

				fmt.Println(config)
				{
					switch strings.ToLower(command[1]) {
					case "get":
						reponse := []string{command[2], (*config)[command[2]]}
						conn.Write([]byte(formatArrayString(reponse)))
					}
				}
			}
		default:
		}
	}
}

func getFlag() map[string]string {
	dirPtr := flag.String("dir", "", "the path to the directory where the RDB file is stored")
	dbFilenamePtr := flag.String("dbfilename", "", "the name of the RDB file")
	flag.Parse()
	config := make(map[string]string)
	if *dirPtr != "" {
		config["dir"] = *dirPtr
	}
	if *dbFilenamePtr != "" {
		config["dbfilename"] = *dbFilenamePtr
	}
	return config
}
