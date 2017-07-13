package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

/*
 * address and port to which client will connect i.e. server address
 */
const (
	CONN_HOST    = "localhost"
	CONN_PORT    = "8888"
	CONN_TYPE    = "tcp"
	MAX_MSG_SIZE = 512
)

func main() {

	// connect to this socket
	rAddr, _ := net.ResolveTCPAddr(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	conn, _ := net.DialTCP(CONN_TYPE, nil, rAddr)

	// create buffer for outgoing data
	conn.SetWriteBuffer(MAX_MSG_SIZE)

	// close connection safely before returning from main
	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	//boolean to check if system is in congestion avoidance mode
	var inCongestionAvoidanceMode, maxCapacityReached bool = false, false

	// variable to check throughput
	var totalPacketsSent int = 0

	// read in input from stdin
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to send: ")
	text, _ := reader.ReadString('\n')

	// string to read user choice to close connection
	//var choice string

	// record start time to calculate end-to-end delay
	startTime := time.Now()
	fmt.Println("Initiate connection at t = ", startTime.Format(time.RFC3339))

	for {
		// send to socket
		conn.Write([]byte(text))
		totalPacketsSent += len(text)

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)

		// decide course of action based on server reply
		if strings.Compare(message, "ACK\n") == 0 {

			if inCongestionAvoidanceMode && !maxCapacityReached {
				// already in congestion avoidance mode. So increasing packet size by one.
				fmt.Println("Increasing packet size by 1")
				text = strings.Join([]string{text, "a"}, "")

			} else {
				// ACK is received and NOT in congestion avoidance mode. So doubling packet size.
				fmt.Println("Doubling packet size")
				text = strings.Join([]string{text, text}, "")

			}
			fmt.Println("Current packet size : ", len(text))

		} else {

			if inCongestionAvoidanceMode {
				// NAK received when already in congestion avoidance mode
				fmt.Println("Reached maximum server window capacity")
				maxCapacityReached = true
				text = text[1:]

				// Calculate throughput and end-to-end delay
				fmt.Println("Time = ", time.Now().Format(time.RFC3339))
				fmt.Println("Time taken = ", time.Since(startTime))
				fmt.Println("Total packets sent = ", totalPacketsSent)
				fmt.Println("Throughput = ", float64(totalPacketsSent)/float64(time.Since(startTime)))
				return

			} else {
				// NAK received for the first time
				fmt.Println("Going into congestion avoidance mode")
				inCongestionAvoidanceMode = true

				// reduce data length to half of current length
				text = text[:len(text)/2]

			}
		}

		//fmt.Println("Send packet? (y/n)")
		//choice, _ = reader.ReadString('\n')
		//choice = strings.TrimSpace(choice)
		//if strings.Compare(choice, "y") != 0 {
		//	break
		//}

		fmt.Printf("\n")
	}
}
