package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	var bindAddress = flag.String("bind", "0.0.0.0", "bind address")
	var bindPort = flag.Int("port", 23, "port")

	flag.Parse()

	bind, err := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%d", *bindAddress, *bindPort))
	if err != nil {
		log.Fatalf("Cannot resolve %s: %v", bind, err)
	}

	l, err := net.ListenTCP("tcp", bind)
	if err != nil {
		log.Fatalf("Cannot bind on %s: %v", bind, err)
	}

	defer l.Close()
	log.Printf("Listening on %s", bind)

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Fatalf("Can't accept: %v", err)
		}

		log.Printf("Someone's connected: %v", conn.RemoteAddr())

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	conn.Write(connInfos(conn))
	conn.Close()
}

func connInfos(conn net.Conn) []byte {
	var buff bytes.Buffer

	writeField := func(name, value string) {
		buff.Write([]byte(fmt.Sprintf("%s: %s\n", name, value)))
	}

	remote := conn.RemoteAddr()
	if host, port, err := net.SplitHostPort(remote.String()); err == nil {
		writeField("Remote IP", host)
		writeField("Remote Port", port)

		if names, err := net.LookupAddr(host); err == nil && len(names) > 0 {
			if len(names) == 1 {
				writeField("Remote Name", names[0])
			}
		}
	}

	return buff.Bytes()
}
