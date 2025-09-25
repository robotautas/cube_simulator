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

var sercon = false
var calibration = false

var serconStatus = 1 //idle
var taskList = "25|" +
	"c9bbbd59-eb3f-4da4-8077-28ca68ed2ec6|SampleAnalysis|Test|1|NCHS_O-30s-20min|Blank|" +
	"de7ad592-8c4a-4668-8f47-422fb327c21b|SampleAnalysis|Test|1|NCHS_O-30s-20min|Blank|" +
	"d57e1e48-6db0-470b-87ff-f34bcf1dd4b1|SampleAnalysis|Test|1|NCHS_O-30s-30min|Test|" +
	"499b98f9-10fa-495c-9793-f1cfa90664d8|SampleAnalysis|Test|1|NCHS_O-30s-30min|Test|" +
	"490b4400-7908-47eb-989b-e7953e0c2c53|SampleAnalysis|Test|1|NCHS_O-30s-30min|Test|" +
	"4c67d33f-d283-4b50-ac28-bae2371b8cd7|SampleAnalysis|RunIn|1|NCHS_O-30s-20min|RunIN|" +
	"48cad7fa-542c-4a85-886d-69f344636142|SampleAnalysis|RunIn|1|NCHS_O-30s-20min|RunIN|" +
	"d37943f1-8b48-43af-bcea-9fa425686e06|SampleAnalysis|Sulf|2.2|NCHS_O-30s-20min|Sulfanilamide20250408|" +
	"91cf89c6-a432-44c3-8c17-280a07d5bfbb|SampleAnalysis|Sulf|2.1|NCHS_O-30s-20min|Sulfanilamide2025040802|" +
	"e4afe9bf-6a49-4e34-a8c4-59c4d23f065f|SampleAnalysis|Sulf|2.1|NCHS_O-30s-20min|Sulfanilamide2025040803|" +
	"9cd56363-3bda-4f9d-8243-4e3622c3e3c2|SampleAnalysis|RunIn|1|NOS SPEEDRUN|RunIn|" +
	"0ad44484-4e9b-4936-a7da-b19d6ff9b4e3|SampleAnalysis|RunIn|1|NOS SPEEDRUN|RunIn|" +
	"2880c45d-9187-4851-94d5-f1514f509fdf|SampleAnalysis|OX|5.54|NOS SPEEDRUN|OX|" +
	"92643f96-95fb-4331-be4f-c20a5b079d79|SampleAnalysis|OX|6|NOS SPEEDRUN|OX2024040902|" +
	"f700de81-e34e-4fe9-a65e-a9575a8b09a5|SampleAnalysis|OX|5.07|NOS SPEEDRUN|OX2025040903|" +
	"ed925789-48b0-43e0-a08e-74ddcedd6d69|SampleAnalysis|OX|4.81|NOS SPEEDRUN|OX2025090404|" +
	"4eea58fa-1b31-46e3-ae72-98d968f51c9c|SampleAnalysis|OX|5.91|NOS SPEEDRUN|OX2025040905|" +
	"59d3b605-6878-41f9-9a5c-53a2a0905598|SampleAnalysis|RunIn|1|NOS SPEEDRUN|RunIn|" +
	"99999999-9999-9999-9999-999999999999|SampleAnalysis|Paskutinis|1|NOS SPEEDRUN|RunIn2|" +
	"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa|SampleAnalysis|Pre|100|NOS SPEEDRUN|RunIn3|" +
	"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb|SampleAnalysis|RunIn|100|NOS SPEEDRUN|RunIn3|" +
	"cccccccc-cccc-cccc-cccc-cccccccccccc|SampleAnalysis|First|100|NOS SPEEDRUN|wasd|" +
	"dddddddd-dddd-dddd-dddd-dddddddddddd|SampleAnalysis|Pirmas|100|NOS SPEEDRUN|wasd|" +
	"eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee|SampleAnalysis|Antras|100|NOS SPEEDRUN|wasd|" // +
	// "71cf9633-3afc-4974-8ad5-4bd5108e974c|Configuration|whatever|Method has not been set|wasd\n"

var serconSequence = []string{
	"99999999-9999-9999-9999-999999999999|SampleAnalysis|Paskutinis|1|NOS SPEEDRUN|RunIn2",
	"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa|SampleAnalysis|Pre|100|NOS SPEEDRUN|RunIn3",
	"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb|SampleAnalysis|RunIn|100|NOS SPEEDRUN|RunIn3",
	"cccccccc-cccc-cccc-cccc-cccccccccccc|SampleAnalysis|First|100|NOS SPEEDRUN|wasd",
	"dddddddd-dddd-dddd-dddd-dddddddddddd|SampleAnalysis|Pirmas|123|NOS SPEEDRUN|wasd",
	"eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee|SampleAnalysis|Antras|234|NOS SPEEDRUN|wasd",
}

var serconLastSample = 0
var serconCurrentSample = 0
var serconNextSample = 1 // index in sercon sequence

func main() {

	for _, arg := range os.Args {
		if strings.ToLower(arg) == "earth" {
			sercon = true
			break

		}
	}

	for _, arg := range os.Args {
		if strings.ToLower(arg) == "cal" {
			calibration = true
			break

		}
	}

	if calibration {
		taskList = "5|" +
			"99999999-9999-9999-9999-999999999999|SampleAnalysis|Paskutinis|1|NOS SPEEDRUN|RunIn2" +
			"c9bbbd59-eb3f-4da4-8077-28ca68ed2ec1|SampleAnalysis|CalStart|1|NCHS_O-30s-20min|Blank|" +
			"de7ad592-8c4a-4668-8f47-422fb327c212|SampleAnalysis|Nr_1|1|NCHS_O-30s-20min|Blank|" +
			"f700de81-e34e-4fe9-a65e-a9575a8b0915|SampleAnalysis|Nr_14|5.07|NOS SPEEDRUN|OX2025040903|" +
			"ed925789-48b0-43e0-a08e-74ddcedd6d16|SampleAnalysis|CalEnd|4.81|NOS SPEEDRUN|OX2025090404|"

		serconSequence = []string{
			"99999999-9999-9999-9999-999999999999|SampleAnalysis|Paskutinis|1|NOS SPEEDRUN|RunIn2",
			"c9bbbd59-eb3f-4da4-8077-28ca68ed2ec1|SampleAnalysis|CalStart|1|NCHS_O-30s-20min|Blank",
			"de7ad592-8c4a-4668-8f47-422fb327c212|SampleAnalysis|Nr_1|1|NCHS_O-30s-20min|Blank",
			"f700de81-e34e-4fe9-a65e-a9575a8b0915|SampleAnalysis|Nr_14|5.07|NOS SPEEDRUN|OX2025040903",
			"ed925789-48b0-43e0-a08e-74ddcedd6d16|SampleAnalysis|CalEnd|4.81|NOS SPEEDRUN|OX2025090404",
		}

		// serconNextSample = 0
	}

	if len(os.Args) < 2 {
		fmt.Println("need port specified in args, please run executable like 'cube.exe 4321' or 'cube.exe 4444 earth' for sercon version")
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

	fmt.Printf("TCP server listening on port %s\n", port)
	if sercon {
		fmt.Printf("Emulating Sercon Earth analyzer ")
	} else {
		fmt.Printf("Emulating Elementar analyzer")
	}

	if calibration {
		fmt.Printf("in calibration mode")
	}

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
		if sercon {
			earthRespond(message, conn)
		} else {

			respond(message, conn)
		}

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

func earthRespond(m string, conn net.Conn) {

	switch {
	case m == "?TASKS\n":
		conn.Write([]byte("7\n"))
	case m == "?SAMPLES\n":
		conn.Write([]byte("6\n"))
	case m == "?TASKCOUNT\n":
		conn.Write([]byte("25\n"))
	case m == "?CURRENT\n":
		fmt.Fprintf(conn, "%s\n", serconSequence[serconCurrentSample])
	case m == "?LAST\n":
		fmt.Fprintf(conn, "%s\n", serconSequence[serconLastSample])
	case m == "?NEXT\n":
		fmt.Fprintf(conn, "%s\n", serconSequence[serconNextSample])
	case m == "?TASKLIST\n":
		fmt.Fprint(conn, taskList)
	case m == "?STATUS\n":
		fmt.Fprintf(conn, "%d\n", serconStatus)
	case m == "?RTR\n":
		fmt.Fprint(conn, "1\n")
	case strings.HasPrefix(m, "?PEAK"):
		fmt.Fprint(conn, "1.2345\n")
	case m == "START W\n" || m == "START WAIT\n":
		serconStatus = 4
		fmt.Fprint(conn, "00")
	case m == "STARTNEXT\n":
		serconStatus = 1
		go startNext(conn)
	case m == "RESETSIM\n":
		serconLastSample = 0
		serconCurrentSample = 0
		serconNextSample = 1 // index in sercon sequence

	case m == "?LASTYIELD 1 1 1\n":
		fmt.Fprint(conn, "1.2\n")

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

func startNext(conn net.Conn) {
	fmt.Fprint(conn, "00")
	serconNextSample++
	serconCurrentSample++
	serconLastSample++
	serconStatus = 1

	// cia vyksta visas deginimas.
	for i := 0; i < 30; i++ {
		fmt.Printf("Processing sample %v ... \n", serconSequence[serconCurrentSample])
		time.Sleep(1 * time.Second)
	}

	if serconCurrentSample != len(serconSequence)-1 {

		serconStatus = 4
	} else {
		serconStatus = 0
	}

	fmt.Printf("Burning sample %d complete, sercon status is now %d", serconCurrentSample, serconStatus)

}
