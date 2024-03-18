package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	//"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// declare version varaible
var version = "0.0.1"

func main() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No Playdate consoles have been found. Maybe plug it in and try again?")
		return
	}

	found := false
	name := ""
	for _, port := range ports {
		if port.VID == "1331" && port.PID == "5740" {
			fmt.Println("Found Playdate device: " + port.Name)
			found = true
			name = port.Name
			break
		}
	}

	if !found {
		fmt.Println("No Playdate consoles have been found. Maybe plug it in and try again?")
		os.Exit(1)
	}

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(name, mode)
	if err != nil {
		log.Fatal(err)
	}
	port.Write([]byte("echo off\n"))
	buf := make([]byte, 512)

	for {
		n, err := port.Read(buf)
		if err == io.EOF {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		for _, line := range strings.Split(string(buf[:n]), "\n") {
			if strings.HasPrefix(line, "playlink") {
				if len(strings.Split(line, "|")) >= 3 && strings.Split(line, "|")[2] == version {
					action := strings.Split(line, "|")[1] // PlayLink action
					switch action {
					case "init":
						fmt.Println("Sending initialized message...")
						port.Write([]byte("msg playlink|init_ok|0.0.1\n"))
					case "ping":
						port.Write([]byte("msg playlink|pong|0.0.1\n"))
					case "connected":
						fmt.Println("Playdate successfully connected!")
					case "request":
						method := strings.Split(line, "|")[3]
						url := strings.Split(line, "|")[4]
						identifier := strings.Split(line, "|")[5]
						fmt.Println("Received request for " + url + " with method " + method)
						if method == "GET" {
							client := resty.New()

							resp, _ := client.R().Get(url)

							switch contentType := resp.Header().Get("Content-Type"); contentType {
							case "application/json", "application/geo+json", "application/json; charset=utf-8":
								var data interface{}
								err := json.Unmarshal(resp.Body(), &data)
								if err != nil {
									log.Fatal(err)
								}
								bodyJSON, err := json.Marshal(data)
								if err != nil {
									log.Fatal(err)
								}
								//headers := resp.Header()

								//for key, values := range headers {
								//	fmt.Println("Key:", key, "Values:", values)
								//}

								port.Write([]byte("msg playlink|response|" + strconv.Itoa(resp.StatusCode()) + "|" + base64.StdEncoding.EncodeToString([]byte(identifier)) + "\n"))
								
								port.Write([]byte("msg " + "startresp|" + "\n"))
								for _, section := range splitString(string(bodyJSON), 85) {
									//fmt.Println("part: " + base64.StdEncoding.EncodeToString([]byte(section)))
									port.Write([]byte("msg " + base64.StdEncoding.EncodeToString([]byte(section))))
									port.Write([]byte("\n"))
									time.Sleep(10 * time.Millisecond)
								}
								port.Write([]byte("msg " + "endresp|" + "\n"))
							}

						}

					}
				} else if strings.Split(line, "|")[2] != version {
					port.Write([]byte("msg playlink|err_ver|" + version + "\n")) // send message to serial saying wrong version, along with the expected version
				}
			}
			if len(os.Args) > 1 && os.Args[1] == "--debug" {
				if strings.HasPrefix(line, "recv") {
					fmt.Println(line)
				}
			}
		}
	}
}

func splitString(s string, chunkSize int) []string {
	var chunks []string

	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}

	return chunks
}
