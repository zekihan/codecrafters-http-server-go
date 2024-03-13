package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

const (
	okContent       = "HTTP/1.1 200 OK\r\n\r\n"
	notFoundContent = "HTTP/1.1 404 Not Found\r\n\r\n"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	read := make([]byte, 1024)
	_, err := conn.Read(read)
	if err != nil {
		return
	}
	s := string(read)
	split := strings.Split(s, "\r\n")
	startLine := strings.Split(split[0], " ")
	path := startLine[1]

	content := notFoundContent
	if path == "/" {
		content = okContent
	}

	_, err = conn.Write([]byte(content))
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
