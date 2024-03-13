package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const crlf = "\r\n"
const lf = "\n"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "127.0.0.1:4221")
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
	parseHttpResponse := parseHttp(string(read))
	path := parseHttpResponse.Path

	content := response(http.StatusNotFound)
	if path == "/" {
		content = response(http.StatusOK)
	} else if strings.HasPrefix(path, "/echo/") {
		params := strings.TrimPrefix(path, "/echo/")
		content = responseWithBody(http.StatusOK, params)
	} else if strings.HasPrefix(path, "/user-agent") {
		userAgent := parseHttpResponse.Headers["User-Agent"]
		content = responseWithBody(http.StatusOK, userAgent)
	}

	_, err = conn.Write([]byte(content))
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}

type ParseHttpResponse struct {
	Headers map[string]string
	Path    string
	Method  string
}

func parseHttp(s string) ParseHttpResponse {
	split := strings.Split(s, crlf)
	startLine := strings.Split(split[0], " ")
	method := startLine[0]
	path := startLine[1]
	headers := make(map[string]string)
	for i := 1; i < len(split); i++ {
		if split[i] == "" {
			break
		}
		header := strings.Split(split[i], ": ")
		headers[header[0]] = header[1]
	}
	return ParseHttpResponse{
		Headers: headers,
		Path:    path,
		Method:  method,
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
