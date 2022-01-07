package server

import (
	"Pchat/internal/server/conf"
	"go.uber.org/zap"
	"gogit.oa.com/March/gopkg/protocol/bypack"
	"net"
	"time"
)

func getIdsByReader(reader *bypack.Reader, cnt int) (ids []int32) {
	defer func() {
		recover()
	}()

	for i := 0; i < cnt; i++ {
		ids = append(ids, reader.Int())
	}
	return ids
}

func logField(conn net.Conn) zap.Field {
	if conn == nil {
		return zap.String("conn", "")
	}
	return zap.String("conn", conn.RemoteAddr().String())
}

func send(conn net.Conn, data []byte) {
	if conn == nil {
		return
	}

	if err := conn.SetWriteDeadline(time.Now().Add(30 * time.Second)); err != nil {
		conf.L.Error(err.Error(), logField(conn))
		return
	}

	for i := 0; i < 3; i++ {
		_, err := conn.Write(data)
		if err != nil {
			conf.L.Warn(err.Error(), logField(conn))
			time.Sleep(1 * time.Second)
		}
	}
}
