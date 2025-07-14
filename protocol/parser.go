package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Parser(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("exprected RESP array")

	}

	line = strings.TrimSuffix(line, "\r\n")
	fmt.Println("after treaming suffix ", line)
	fmt.Println("line is", line)
	if len(line) == 0 || line[0] != '*' {
		fmt.Println("Invalid  resp array")
		return nil, fmt.Errorf("invalid resp array")
	}
	//treaming first character *
	line = line[1:]
	fmt.Println("line after treaming is", line)
	//converting  string to integer
	fmt.Println("length of line is", len(line))
	if len(line) != 1 {
		fmt.Println("Wrong number of arguments")

	}
	size, err := strconv.Atoi(line)
	if err != nil {
		return nil, fmt.Errorf("can't convert  string to integer")
	}
	fmt.Println("size is", size)

	args := make([]string, 0, size)
	// args := make([]string, size)

	for range size {
		fmt.Println("for loop started....")
		line, err := reader.ReadString('\n')
		//checking valid args
		if !strings.HasPrefix(line, "$") || len(line) == 0 {
			return nil, fmt.Errorf("expected bulk string")

		}
		fmt.Println("line after removing backslash n is", line)
		if err != nil {
			return nil, fmt.Errorf("invalid string")
		}
		fmt.Println("trimming of suffix started")
		line = strings.TrimSuffix(line, "\r\n")
		fmt.Println("line after removing r and n is", line)
		//now remove prefix $
		line = line[1:]
		fmt.Println("line after removing prefix", line)
		n, err := strconv.Atoi(line)
		fmt.Println("value of n is:", n)
		if err != nil {
			return nil, fmt.Errorf("can't convert  string to integer")
		}
		buf := make([]byte, n+2)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, fmt.Errorf("failed to read actual values:%v", err)

		}
		actualArg := string(buf[:n])
		fmt.Println("actual argument is", actualArg)
		args = append(args, actualArg)

	}
	return args, nil

}
