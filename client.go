package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"persona/client"
)

var Address = "localhost"

func dialDefaultServer(port int) {
	log.SetFlags(log.Lshortfile)
	addr := fmt.Sprintf("%s:%d", Address, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not contact server because: %s\n", err)
	}

	manager(conn)
}

func DialTLSServer(port int) {
	log.SetFlags(log.Lshortfile)
	certPool, err := client.LoadServerCerts()
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    certPool,
		InsecureSkipVerify: true,
	}

	addr := fmt.Sprintf(":%d", port)
	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("could not start tls client: \n %s\n", err)
	}

	log.Println("tls client successfully connected")
	if err := conn.Handshake(); err != nil {
		log.Fatalf("could not complete handshake:\n  %s\n", err)
	}

	log.Println("tls handshake with server was successfull!")
	manager(conn)
}

func manager(conn net.Conn) {
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
	var secure bool
	flag.IntVar(&port, "port", 4000, "Port to where the server, is listening")
	flag.BoolVar(&secure, "secure", false, "Use TLS or default Server. Default is false")
	flag.Parse()

	if secure {
		log.Println("[INFO] starting server TLS configured!")
		DialTLSServer(port)
	} else {
		log.Println("[WARN] contacting server without TLS configured!")
		dialDefaultServer(port)
	}
}
