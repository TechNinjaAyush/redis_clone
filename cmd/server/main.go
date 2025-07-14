package main

import (
	"bufio"
	"fmt"
	"net"
	"redis/command"
	"redis/internal/db"
	"redis/protocol"
	"strings"
)

var sharedStore = &db.Store{
	Mpp: make(map[string]string),
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// we use buffer to read data in chunks
	fmt.Println("hello world ....")
	reader := bufio.NewReader(conn)
	fmt.Println("shared store is", sharedStore.Mpp)

	for {
		args, err := protocol.Parser(reader)
		fmt.Println("args is", args, "and error is", err)
		cmd := strings.ToUpper(args[0])

		switch cmd {

		case "SET":
			command.SetCommand(args[1], args[2], sharedStore)
			conn.Write([]byte("+OK\r\n"))

		case "PING":
			conn.Write([]byte("PONG\r\n"))

		case "GET":
			value := command.GetCommand(args[1], sharedStore)
			conn.Write([]byte(value))
			conn.Write([]byte("+OK\r\n"))

		case "DEL":
			count := command.Deletecommand(args[1:], sharedStore)
			conn.Write([]byte(fmt.Sprintf("(integer) %d\r\n", count)))

		case "INCR":
			count, err := command.INCRcommand(args[1], sharedStore)
			if err != nil {
				conn.Write([]byte("(error) ERR value is not an integer or out of range\r\n"))
			}
			conn.Write([]byte(fmt.Sprintf("(integer) %d\r\n", count)))

		case "INCRBY":
			count, err := command.INCRBYcommand(args[1], args[2], sharedStore)
			if err != nil {
				conn.Write([]byte("-" + err.Error() + "\r\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
			}

		case "MULTI":
			conn.Write([]byte("+OK\r\n"))

		}

	}
}

func main() {
	fmt.Println("starting redis......")

	// create a listener to listen to new  incoming connections  on port 8080

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error listening", err)

	}
	defer listener.Close()
	fmt.Println("Server is  listening on port 6379")

	// handle incoming connections
	// we use  infinite for loop so that we  have to run  redis server  infinite so that it can behave like an   online server
	for {

		fmt.Println("conneting.....")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error in handling incoming connections ", err.Error())
			continue
		}
		fmt.Println("Client connected succesfully")
		// handle connection
		go handleConnection(conn)

	}

}
