package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
)

const crlf = "\r\n"

var dir = "./"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	dirFlag := flag.String("directory", "", "The directory to serve files from. Defaults to the current directory.")
	flag.Parse()
	if dirFlag != nil && *dirFlag != "" {
		dir = *dirFlag
	}

	fmt.Println("Serving files from", dir)

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
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}
	request := string(buffer[:n])
	parseHttpResponse := parseHttp(request)
	urlPath := parseHttpResponse.Path

	content := response(http.StatusNotFound)
	if urlPath == "/" {
		content = response(http.StatusOK)
	} else if strings.HasPrefix(urlPath, "/echo/") {
		params := strings.TrimPrefix(urlPath, "/echo/")
		content = responseWithBody(http.StatusOK, params)
	} else if strings.HasPrefix(urlPath, "/user-agent") {
		userAgent := parseHttpResponse.Headers["User-Agent"]
		content = responseWithBody(http.StatusOK, userAgent)
	} else if strings.HasPrefix(urlPath, "/files/") {
		filePath := strings.TrimPrefix(urlPath, "/files/")
		localFilePath := path.Join(dir, filePath)
		if parseHttpResponse.Method == "GET" {
			if stat, err := os.Stat(localFilePath); os.IsNotExist(err) || stat.IsDir() {
				content = response(http.StatusNotFound)
			} else {
				file, err := os.ReadFile(localFilePath)
				if err != nil {
					content = responseWithBody(http.StatusInternalServerError, err.Error())
				}
				content = responseWithFile(http.StatusOK, string(file))
			}
		} else if parseHttpResponse.Method == "POST" {
			err := os.WriteFile(localFilePath, []byte(parseHttpResponse.Body), 0644)
			if err != nil {
				content = responseWithBody(http.StatusInternalServerError, err.Error())
			}
			content = response(http.StatusCreated)
		}
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
	Body    string
}

func parseHttp(s string) ParseHttpResponse {
	split := strings.Split(s, crlf)
	startLine := strings.Split(split[0], " ")
	method := startLine[0]
	urlPath := startLine[1]
	headers := make(map[string]string)
	for i := 1; i < len(split); i++ {
		if split[i] == "" {
			break
		}
		header := strings.Split(split[i], ": ")
		headers[header[0]] = header[1]
	}
	if method == "POST" {
		body := split[len(split)-1]
		return ParseHttpResponse{
			Headers: headers,
			Path:    urlPath,
			Method:  method,
			Body:    body,
		}
	}
	return ParseHttpResponse{
		Headers: headers,
		Path:    urlPath,
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

func responseWithFile(status int, body string) string {
	res := fmt.Sprintf("HTTP/1.1 %d %s", status, http.StatusText(status))
	res += crlf
	res += "Content-Type: application/octet-stream"
	res += crlf
	res += "Content-Length: " + fmt.Sprint(len(body))
	res += crlf
	res += crlf
	res += body
	return res
}
