package main

import (
	"fmt"
	"net"
	"time"
)

/*
 * constants for server host address, port number and protocol
 */
const (
	CONN_HOST    = "localhost"
	CONN_PORT    = "8888"
	CONN_TYPE    = "tcp"
	MAX_MSG_SIZE = 256
)

/*
 * function to handle client request
 * i.e. check if incoming data has exceeded server window size or not
 */
func handleConnection(conn *net.TCPConn) {
	fmt.Println("Handling new connection...")

	// Close particular client connection when this function ends
	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	processingTime := time.Millisecond * 100
	message := make([]byte, 2*MAX_MSG_SIZE)
	for {
		// Read data sent by client delimited by newline
		n, _ := conn.Read(message[:])

		// if no data is sent, return from function
		if n == 0 {
			return
		}

		// array to store server response
		var response []byte
		fmt.Println("Data length : ", n)

		if n <= MAX_MSG_SIZE {
			response = []byte("ACK\n")
		} else {
			response = []byte("NAK\n")
		}
		time.Sleep(processingTime)

		// send data to client
		conn.Write(response)
	}
}

func main() {
	// Start listening to port 8888 for TCP connection
	rAddr, _ := net.ResolveTCPAddr(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	listener, err := net.ListenTCP(CONN_TYPE, rAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server listening for incoming connections")

	// close connection safely before shutting down server
	defer func() {
		listener.Close()
		fmt.Println("Listener closed")
	}()

	// listen to port forever until closed explicitly
	for {
		// Get net.TCPConn object when client requests
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			break
		}

		// handle client request in a separate thread
		go handleConnection(conn)
	}
}
