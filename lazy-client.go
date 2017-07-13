package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"net"
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

/*
 * Optimization of slow start algorithm using lazy caterer's sequence
 *
 * L(n) = (n^2 + n + 2)/2
 */

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// function to generate a random string of data of length n
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// function to return the number in lazy caterer's sequence
func lazyCaterer(sequenceNum int) int {
	return (int(math.Pow(float64(sequenceNum), float64(2))) + sequenceNum + 2) / 2
}

func main() {

	// connect to this socket
	rAddr, _ := net.ResolveTCPAddr(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	conn, _ := net.DialTCP(CONN_TYPE, nil, rAddr)

	conn.SetWriteBuffer(MAX_MSG_SIZE)

	defer func() {
		fmt.Println("Closing connection")
		conn.Close()
	}()

	//boolean to check if system is in congestion avoidance mode
	var inCongestionAvoidanceMode, maxCapacityReached bool = false, false

	var totalPacketsSent, sequenceNum int = 0, 0

	// generate initial string of length 1
	var text string = RandStringBytes(1)

	// string to read user choice to close connection
	//var choice string

	// Start time
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

				fmt.Println("Increasing packet size by 1")
				text = strings.Join([]string{text, "a"}, "")

			} else {

				fmt.Println("Generating next number in Lazy-Caterer's sequence")
				oldPacketSize := lazyCaterer(sequenceNum)
				sequenceNum = sequenceNum + 1
				newPacketSize := lazyCaterer(sequenceNum)
				text = text + RandStringBytes(newPacketSize-oldPacketSize)

			}
			fmt.Println("Current packet size : ", len(text))
		} else {

			if inCongestionAvoidanceMode {

				fmt.Println("Reached maximum server window capacity")
				maxCapacityReached = true
				text = text[1:]
				fmt.Println("Time = ", time.Now().Format(time.RFC3339))
				fmt.Println("Time taken = ", time.Since(startTime))
				fmt.Println("Total packets sent = ", totalPacketsSent)
				return

			} else {

				fmt.Println("Going into congestion avoidance mode")
				inCongestionAvoidanceMode = true
				oldPacketSize := lazyCaterer(sequenceNum - 1)
				text = text[:oldPacketSize]

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
