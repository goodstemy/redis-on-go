package listener

import (
	"bytes"
	"io"
	"log"
	"net"
)

func Listen(rcv chan bytes.Buffer, snd chan []byte) {
	l, err := net.Listen("tcp", ":4242")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	defer close(rcv)

	var buf bytes.Buffer

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		io.Copy(&buf, conn)

		rcv <- buf

		conn.Write(<-snd)

		conn.Close()

		buf.Reset()
	}
}
