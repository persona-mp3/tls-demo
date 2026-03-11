package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var Address = "localhost"

func dialDefaultServer(port int) {
	log.SetFlags(log.Lshortfile)
	addr := fmt.Sprintf("%s:%d", Address, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not contact server because: %s\n", err)
	}

	stdinCh := make(chan string)
	responseCh := make(chan string)
	go collectInput(stdinCh)
	go fromServer(conn, responseCh)
	for {
		select {
		case msg, open := <-stdinCh:
			if !open {
				return
			}
			if _, err := fmt.Fprintf(conn, "%s", msg); err != nil {
				log.Printf("error occured writing message to server %s\n", err)
				return
			}

		case serverMsg, ok := <-responseCh:
			if !ok {
				return
			}
			fmt.Printf("\t %s\n", serverMsg)
		}
	}
}

func fromServer(conn net.Conn, out chan<- string) {
	buffer := make([]byte, 1024)
	defer close(out)
	for {
		n, err := conn.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			log.Printf("server disconnected!\n")
			return
		} else if err != nil {
			log.Printf("an unexpected error occured: %s\n", err)
			return
		}

		fmt.Println("New message from server!")
		out <- string(buffer[:n])
	}
}
func collectInput(out chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print(" [*] ")
		out <- scanner.Text()
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	var port int
	flag.IntVar(&port, "port", 4000, "Port to where the server, is listening")
	flag.Parse()

	log.Println("[WARN] contacting server without TLS configured!")
	dialDefaultServer(port)
}
