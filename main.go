package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var names = []string{
	"sample 1",
	"sample 2",
	"sample 3",
	"sample 4",
	"sample 5",
	"sample 6",
	"sample 7",
	"sample 8",
	"sample 9",
	"sample 10",
}

func main() {
	// Listen on TCP port 8080
	ln, err := net.Listen("tcp", ":1984")
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer ln.Close()

	fmt.Println("TCP server listening on port 1984")

	for {
		// Accept a connection
		println("\nWaiting for connection...")
		conn, err := ln.Accept()
		println(conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle the connection in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading:", err.Error())
			} else {
				fmt.Println("Client disconnected")
				return
			}
			break
		}

		fmt.Printf("Message received: '%s'", message)
		// Add any additional handling for the received message here

		// Optionally, send a response back to the client
		respond(message, conn)

	}
}

func respond(m string, conn net.Conn) {
	switch m {
	case "?STS\r\n":
		conn.Write([]byte("STS READY"))
	case "START\r\n":
		go startSequence(conn)
	default:
		conn.Write([]byte(fmt.Sprintf("Response message to %s\n", m)))
	}
}

func startSequence(conn net.Conn) {
	conn.Write([]byte("STARTED SEQUENCE!\n"))
	file, err := os.Open("sequence.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var lastTime time.Time
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 2)
		println(line)

		if len(parts) < 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}

		currentTime, err := time.Parse("15:04:05", parts[0])
		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			continue
		}

		if !lastTime.IsZero() {
			waitDuration := currentTime.Sub(lastTime)
			if waitDuration > 0 {
				time.Sleep(waitDuration)
			}
		}

		lastTime = currentTime

		fmt.Println(parts[1]) // Print the line excluding the timestamp
		conn.Write([]byte(fmt.Sprintf("%s\r\n", parts[1])))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}
