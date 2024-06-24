package main

import (
	"bytes"
	"fmt"
	"log"

	// Uncomment this block to pass the first stage
	"compress/gzip"
	"io"
	"net"
	"os"
	"strings"
)

const CRLF = "\r\n"
const HTTP = "HTTP/1.1"
const NOT_FOUND = "404 Not Found"
const OK = "200 OK"
const CREATED = "201 Created"
const CONTENT_TYPE_TEXT_PLAIN = "Content-Type: text/plain"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
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
		go handelReq(conn)
	}

}

func handelReq(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from open socket", err.Error())
	}
	req := string(buffer[:n])
	// find req line
	req_line_index := strings.Index(req, CRLF)
	reqline := strings.Split(req[0:req_line_index], " ")
	request_target := reqline[1]
	// fmt.Println("reqline[0]: ", reqline[0])
	request_path := strings.Split(request_target, "/")
	// computing headers
	req_header_index := strings.Index(req, CRLF+CRLF)
	header_list := req[req_line_index+1 : req_header_index]
	header_map := headerContent(header_list)
	headers := prepHeaders(header_map, "")
	response_body := ""

	// body
	// body := req[req_header_index:]
	// response := ""
	status_line := HTTP + " " + OK + CRLF
	if request_target == "/" {

	} else if request_path[1] == "echo" {
		if len(request_path) >= 3 {
			headers = prepHeaders(header_map, request_path[2])
			response_body = request_path[2]
		}
	} else if request_path[1] == "user-agent" {
		agent_val := strings.Split(header_list, "User-Agent: ")
		if len(agent_val) >= 2 {
			response_body = agent_val[1]
			headers = prepHeaders(header_map, response_body)
		}
	} else if request_path[1] == "files" {
		// Grab directory flag from command line args
		home_dir := ""
		if os.Args[1] == "--directory" {
			home_dir = os.Args[2]
		}
		// Read the file
		if len(request_path) >= 3 {
			file_name := fmt.Sprintf("%s%s", home_dir, request_path[2])
			if reqline[0] == "POST" {
				err := os.MkdirAll(home_dir, 0755)
				if err != nil {
					fmt.Println("Error creating directory:", err)
				} else {
					file, err := os.OpenFile(file_name, os.O_WRONLY|os.O_CREATE, 0644)
					if err != nil {
						fmt.Println("Error opening or creating file:", err)
					}
					defer file.Close()
					data := strings.Split(req, CRLF+CRLF)[1]
					_, err = io.WriteString(file, data)
					if err != nil {
						fmt.Println("Error Writing  file:", err)
					}
					status_line = HTTP + " " + CREATED + CRLF
					response_body = ""
				}
			} else {
				file, err := os.Open(file_name)
				if err != nil {
					status_line = HTTP + " " + NOT_FOUND + CRLF
					headers = CRLF
				} else {
					content, err := io.ReadAll(file)
					if err != nil {
						log.Fatal(err)
					}
					header_map["Content-Type"] = "application/octet-stream"
					headers = prepHeaders(header_map, string(content))
					response_body = string(content)

				}
				defer file.Close()
			}
		}
	} else {
		status_line = HTTP + " " + NOT_FOUND + CRLF
	}
	if header_map["Accept-Encoding"] == "gzip" {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		gz.Write([]byte(response_body))
		err = gz.Close()
		response_body = string(buf.Bytes())
		headers = prepHeaders(header_map, string(response_body))
	}

	final_response := status_line + headers + response_body
	byte_data := []byte(final_response)

	conn.Write(byte_data)
	conn.Close()
}
func prepHeaders(headerMap map[string]string, content string) string {
	h1 := fmt.Sprintf("Content-Type: %s", headerMap["Content-Type"])
	headers := h1 + CRLF
	h2 := ""
	if headerMap["Accept-Encoding"] == "gzip" {
		h2 = fmt.Sprintf("Content-Encoding: %s", headerMap["Accept-Encoding"])
		headers += h2 + CRLF
	}
	if content != "" {
		h2 = fmt.Sprintf("Content-Length: %d", len(content))
		headers += h2 + CRLF
	}
	return headers + CRLF
}
func headerContent(headerList string) map[string]string {
	headerList = strings.TrimSpace(headerList)
	all_headers := strings.Split(headerList, CRLF)
	header_map := make(map[string]string)
	for _, element := range all_headers {
		temp := strings.Split(element, ": ")
		if temp[0] == "Accept-Encoding" {
			all_encodings := strings.Split(temp[1], ",")
			for _, encoding := range all_encodings {
				encoding = strings.TrimSpace(encoding)
				if "gzip" == encoding {
					temp[1] = encoding
					header_map[temp[0]] = temp[1]
				}
			}
		} else {
			header_map[temp[0]] = temp[1]
		}
	}
	if _, ok := header_map["Content-Type"]; !ok {
		header_map["Content-Type"] = "text/plain"
	}
	for key, value := range header_map {
		fmt.Printf("%s: %s\n", key, value)
	}
	return header_map
}

// TO DO
func statusContent() string {
	statusLine := ""
	return statusLine
}

// TO DO
func bodyContent() string {
	response_body := ""
	return response_body
}
