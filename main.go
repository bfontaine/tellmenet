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

	bind := fmt.Sprintf("%s:%d", *bindAddress, *bindPort)

	l, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("Cannot bind on %s: %v", bind, err)
	}

	defer l.Close()
	log.Printf("Listening on %s", bind)

	for {
		conn, err := l.Accept()
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

	writeField := func(name, value interface{}) {
		buff.Write([]byte(fmt.Sprintf("%s: %v\n", name, value)))
	}

	remote := conn.RemoteAddr()

	writeField("Network", remote.Network())

	if host, port, err := net.SplitHostPort(remote.String()); err == nil {
		ip := net.ParseIP(host)

		writeField("IP", ip)
		writeField("Global unicast", ip.IsGlobalUnicast())
		writeField("Multicast", ip.IsMulticast())
		writeField("Interface-local multicast", ip.IsInterfaceLocalMulticast())
		writeField("Link-local multicast", ip.IsLinkLocalMulticast())
		writeField("Link-local unicast", ip.IsLinkLocalUnicast())
		writeField("Loopback", ip.IsLoopback())
		writeField("Port", port)

		if names, err := net.LookupAddr(host); err == nil && len(names) > 0 {
			if len(names) == 1 {
				writeField("Name", names[0])
			} else {
				buff.Write([]byte("Names:\n"))
				for _, name := range names {
					buff.Write([]byte(fmt.Sprintf("  %s\n", name)))
				}
			}
		}
	}

	return buff.Bytes()
}
