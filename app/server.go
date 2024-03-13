package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const crlf = "\r\n"

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
	split := strings.Split(s, crlf)
	startLine := strings.Split(split[0], " ")
	path := startLine[1]

	content := response(http.StatusNotFound)
	if path == "/" {
		content = response(http.StatusOK)
	} else if strings.HasPrefix(path, "/echo/") {
		params := strings.TrimPrefix(path, "/echo/")
		content = responseWithBody(http.StatusOK, params)
	}

	_, err = conn.Write([]byte(content))
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}

func response(status int) string {
	return responseWithBody(status, "")
}

func responseWithBody(status int, body string) string {
	res := fmt.Sprintf("HTTP/1.1 %d %s", status, http.StatusText(status))
	res += crlf
	res += "Content-Type: text/plain"
	res += crlf
	res += "Content-Length: " + fmt.Sprint(len(body))
	res += crlf
	res += crlf
	res += body
	return res
}
