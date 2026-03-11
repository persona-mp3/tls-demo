package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"persona/server"
)

func startDefaultServer(port int) {
	log.SetFlags(log.Lshortfile)
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("could not start default server", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("could not accpet connection: ", err)
			continue
		}

		go handleConnection(conn)
	}
}

func startTLSServer(port int) {
	log.SetFlags(log.Lshortfile)
	config, err := server.LoadTlsConfig()

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("could not start default server", err)
	}

	tlsListener := tls.NewListener(listener, config)
	log.Printf("[INFO] tls server running on port %d\n", port)
	for {
		conn, err := tlsListener.Accept()
		if err != nil {
			log.Println("could not accpet connection: ", err)
			continue
		}

		passedHandshake := server.InitHandshake(conn)
		if !passedHandshake {
			conn.Close()
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.SetFlags(log.Lshortfile)
	defer conn.Close()
	if _, err := conn.Write([]byte("Hello there client!")); err != nil {
		log.Printf("error writing to connection %s\n", err)
		return
	}

	buffer := make([]byte, 1024)
	for {
		bytesRead, err := conn.Read(buffer)
		if err != nil && errors.Is(err, io.EOF) {
			log.Printf("client unexpectedly disconnected\n")
			return
		}

		content := buffer[:bytesRead]
		log.Printf("content from client: %s\n", content)
		log.Println("echoing back to client")
		fmt.Fprintf(conn, "ISERVER: %s\n", content)
	}

}

func main() {
	log.SetFlags(log.Lshortfile)
	var port int
	var secure bool
	flag.IntVar(&port, "port", 4000, "Port to listen on")
	flag.BoolVar(&secure, "secure", false, "Use TLS or default Server. Default is false")
	flag.Parse()

	if secure {
		log.Println("[INFO] starting server TLS configured!")
		startTLSServer(port)
	} else {
		log.Println("[WARN] starting server without TLS configured!")
		startDefaultServer(port)
	}
}
