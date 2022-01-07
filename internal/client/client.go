package client

import (
	"Pchat/internal/client/conf"
	"fmt"
	"gogit.oa.com/March/gopkg/protocol/bypack"
	"io"
	"net"

	"gogit.oa.com/March/gopkg/util"
)

func Run() {
	reader := SendAndRecv(buffer888())
	conf.L.Info(reader.String())
}

func Send(data []byte) {
	_, err := Conn().Write(data)
	util.MustNil(err)
}

func SendAndRecv(data []byte) *bypack.Reader {
	conn := Conn()
	n, err := conn.Write(data)
	fmt.Println("n1:", n)
	util.MustNil(err)

	hb := make([]byte, bypack.HeaderSize)
	n, err1 := io.ReadFull(conn, hb)
	fmt.Println("n:", n)
	util.MustNil(err1)

	header, err2 := bypack.NewHeader(hb)
	util.MustNil(err2)

	body := make([]byte, header.GetSize())
	n, err3 := io.ReadFull(conn, body)
	util.MustNil(err3)

	return bypack.NewReader(header.GetCmd(), body[:n])
}

func Conn() net.Conn {
	conn, err := net.Dial("tcp", conf.Conf.Client.Addr)
	util.MustNil(err)
	return conn
}
