package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var sampleNames = []string{
	"Pre",
	"RunIn",
	"First",
	"Pirmas",
	"Antras",
}

var sampleWeights = []string{
	"1.1", "1.2", "1.3", "1.4", "1.5",
}

var sinxPosition int = 1

func main() {

	if len(os.Args) < 2 {
		fmt.Println("need port specified in args, please run executable like 'cube.exe 4321'")
		return
	}

	port := os.Args[1]

	// Listen on TCP port 8080
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer ln.Close()

	fmt.Printf("TCP server listening on port %s", port)

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
	namPattern := regexp.MustCompile(`^\?NAM\s+(\d+)\r\n$`)
	wghPattern := regexp.MustCompile(`^\?WGH\s+(\d+)\r\n$`)

	switch {
	case m == "?STS\r\n":
		conn.Write([]byte("STS READY\n"))
	case m == "STRT\r\n":
		go startSequence(conn)
	case m == "reset_sinx\r\n":
		sinxPosition = 1
	case m == "?SINX\r\n":
		msg := fmt.Sprintf("SINX %d\n", sinxPosition)
		conn.Write([]byte(msg))
		sinxPosition++
	case namPattern.MatchString(m):
		match := namPattern.FindStringSubmatch(m)
		if len(match) > 1 {
			nameNumber := match[1]

			nameNumberInt, err := strconv.Atoi(nameNumber)
			if err != nil {
				fmt.Printf("cannot convert %s to integer!", nameNumber)
			}

			if nameNumberInt > len(sampleNames) {
				conn.Write([]byte(""))
			} else {
				index := nameNumberInt - 1
				sampleName := sampleNames[index]
				conn.Write([]byte(fmt.Sprintf("NAM %d %s\n", nameNumberInt, sampleName)))
			}
		}
	case wghPattern.MatchString(m):
		match := wghPattern.FindStringSubmatch(m)
		if len(match) > 1 {
			nameNumber := match[1]

			nameNumberInt, err := strconv.Atoi(nameNumber)
			if err != nil {
				fmt.Printf("cannot convert %s to integer!", nameNumber)
			}

			if nameNumberInt > len(sampleNames) {
				conn.Write([]byte(""))
			} else {
				index := nameNumberInt - 1
				sampleName := sampleWeights[index]
				println(fmt.Sprintf("WGH %d %s", nameNumberInt, sampleName))
				conn.Write([]byte(fmt.Sprintf("WGH %d %s\n", nameNumberInt, sampleName)))
			}
		}

	case strings.HasPrefix(m, "?PCT"):
		rand.Seed(time.Now().UnixNano())
		randNum := rand.Intn(100)
		pct := float32(randNum) / 100
		conn.Write([]byte(fmt.Sprintf("PCT X C %f\n", pct)))

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
