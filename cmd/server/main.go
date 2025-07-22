package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"redis/command"
	"redis/internal/db"
	"redis/protocol"
	"strings"
)

var sharedStore = &db.Store{
	Mpp: make(map[string]string),
}

func appendAOF(args []string, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Write command in RESP format for replay
	line := "*" + fmt.Sprintf("%d", len(args)) + "\r\n"
	for _, arg := range args {
		line += "$" + fmt.Sprintf("%d", len(arg)) + "\r\n" + arg + "\r\n"
	}
	_, err = f.WriteString(line)
	return err
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("hello world ....")
	reader := bufio.NewReader(conn)
	fmt.Println("shared store is", sharedStore.Mpp)
	var inTransaction bool
	var queueCommands [][]string

	for {
		args, err := protocol.Parser(reader)
		fmt.Println("args is", args, "and error is", err)
		if err != nil || len(args) == 0 {
			conn.Write([]byte("-ERR invalid command\r\n"))
			continue
		}
		cmd := strings.ToUpper(args[0])

		if inTransaction && cmd != "EXEC" && cmd != "DISCARD" {
			queueCommands = append(queueCommands, args)
			conn.Write([]byte("+QUEUED\r\n"))
			continue
		}
		switch cmd {
		case "SET":
			if len(args) < 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			command.SetCommand(args[1], args[2], sharedStore)
			appendAOF(args, "appendonly.aof")
			conn.Write([]byte("+OK\r\n"))
		case "PING":
			conn.Write([]byte("PONG\r\n"))
		case "GET":
			if len(args) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			value := command.GetCommand(args[1], sharedStore)
			conn.Write([]byte(value))
			conn.Write([]byte("+OK\r\n"))
		case "DEL":
			if len(args) < 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
				continue
			}
			count := command.Deletecommand(args[1:], sharedStore)
			appendAOF(args, "appendonly.aof")
			conn.Write([]byte(fmt.Sprintf("(integer) %d\r\n", count)))
		case "INCR":
			if len(args) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'incr' command\r\n"))
				continue
			}
			count, err := command.INCRcommand(args[1], sharedStore)
			if err != nil {
				conn.Write([]byte("(error) ERR value is not an integer or out of range\r\n"))
			}
			appendAOF(args, "appendonly.aof")
			conn.Write([]byte(fmt.Sprintf("(integer) %d\r\n", count)))
		case "INCRBY":
			if len(args) != 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'incrby' command\r\n"))
				continue
			}
			count, err := command.INCRBYcommand(args[1], args[2], sharedStore)
			if err != nil {
				conn.Write([]byte("-" + err.Error() + "\r\n"))
			} else {
				appendAOF(args, "appendonly.aof")
				conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
			}
		case "MULTI":
			conn.Write([]byte("+OK\r\n"))
			inTransaction = true
			queueCommands = [][]string{}
		case "EXEC":
			for _, cmdArgs := range queueCommands {
				first_arg := strings.ToUpper(cmdArgs[0])
				switch first_arg {
				case "SET":
					if len(cmdArgs) < 3 {
						conn.Write([]byte("-Err wrong  number of arguments for set command\r\n"))
						continue
					}
					command.SetCommand(cmdArgs[1], cmdArgs[2], sharedStore)
					appendAOF(cmdArgs, "appendonly.aof")
					conn.Write([]byte("+OK\r\n"))
				case "DEL":
					if len(cmdArgs) < 2 {
						conn.Write([]byte("-ERR wrong number of arguments for 'del' command\r\n"))
						continue
					}
					command.Deletecommand(cmdArgs[1:], sharedStore)
					appendAOF(cmdArgs, "appendonly.aof")
					conn.Write([]byte("+OK\r\n"))
				case "INCR":
					if len(cmdArgs) != 2 {
						conn.Write([]byte("-ERR wrong number of arguments for 'incr' command\r\n"))
						continue
					}
					command.INCRcommand(cmdArgs[1], sharedStore)
					appendAOF(cmdArgs, "appendonly.aof")
					conn.Write([]byte("+OK\r\n"))
				case "INCRBY":
					if len(cmdArgs) != 3 {
						conn.Write([]byte("-ERR wrong number of arguments for 'incrby' command\r\n"))
						continue
					}
					command.INCRBYcommand(cmdArgs[1], cmdArgs[2], sharedStore)
					appendAOF(cmdArgs, "appendonly.aof")
					conn.Write([]byte("+OK\r\n"))
				}
			}
			inTransaction = false
			queueCommands = nil
		case "DISCARD":
			inTransaction = false
			queueCommands = nil
			conn.Write([]byte("+OK\r\n"))
		}
	}
}

func loadAOF() {
	file, err := os.Open("appendonly.aof")
	if err != nil {
		log.Printf("No existing AOF file:%v", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		args, err := protocol.Parser(reader)
		if err != nil || len(args) == 0 {
			break
		}
		cmd := strings.ToUpper(args[0])
		switch cmd {
		case "SET":
			if len(args) >= 3 {
				command.SetCommand(args[1], args[2], sharedStore)
			}
		case "DEL":
			if len(args) >= 2 {
				command.Deletecommand(args[1:], sharedStore)
			}
		case "INCR":
			if len(args) >= 2 {
				command.INCRcommand(args[1], sharedStore)
			}
		case "INCRBY":
			if len(args) >= 3 {
				command.INCRBYcommand(args[1], args[2], sharedStore)
			}
		}
	}
}
func main() {
	fmt.Println("starting redis......")

	//load append only  file

	loadAOF()

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
