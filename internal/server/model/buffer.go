package model

import (
	"Pchat/internal/server/conf"
	ctx "context"
	"fmt"

	"go.uber.org/zap"
	"gogit.oa.com/March/gopkg/protocol/bypack"
)

var (
	BufferLoginSuccess []byte
	BufferMonitorTest  []byte
	BufferHeartbeat    []byte
)

func init() {
	setBufferLoginSuccess()
	setBufferMonitorTest()
	setBufferHeartbeat()
}

func setBufferLoginSuccess() {
	w := bypack.NewWriter(0x201)
	w.End()
	BufferLoginSuccess = w.GetBuffer()
}

func setBufferMonitorTest() {
	w := bypack.NewWriter(0x881)
	w.End()
	BufferMonitorTest = w.GetBuffer()
}

func setBufferHeartbeat() {
	w := bypack.NewWriter(0x1)
	w.End()
	BufferHeartbeat = w.GetBuffer()
}

func GetBufferAdmin(data string) []byte {
	return getStringBuffer(0x888, data)
}

func GetBufferOnlineNum() []byte {
	w := bypack.NewWriter(0x109)
	w.Int(CountMid())
	w.End()
	return w.GetBuffer()
}

func GetBufferJS(data string) []byte {
	return getStringBuffer(0x10E, data)
}

func GetBuffer887(data string) []byte {
	return getStringBuffer(0x887, data)
}

func GetBuffer886(data string) []byte {
	return getStringBuffer(0x886, data)
}

func getStringBuffer(cmd uint16, data string) []byte {
	w := bypack.NewWriter(cmd)
	w.String(data)
	w.End()
	return w.GetBuffer()
}

func GetOffBuffer(mid int32) []byte {
	key := GetMsgOffKey(mid)
	data := conf.Redis.RPop(ctx.Background(), key).Val()
	conf.L.Info("Get off", zap.Any("mid", mid), zap.Any("data", data))
	if data != "" {
		buff := []byte(data)
		fmt.Println("buff>>>>>>>>>>>>>>>>", buff)
		conf.L.Info("Get off", zap.Any("mid", mid), zap.Any("buff", buff))
		return buff
	}
	return nil

}

func SetOffBuffer(mid int32, data []byte) {
	key := GetMsgOffKey(mid)
	str := string(data)
	conf.L.Info("Set off", zap.Any("mid", mid), zap.Any("str", str))
	conf.Redis.LPush(ctx.Background(), key, str)
}

func GetGroupBuffer(cmd uint16, args ...interface{}) []byte {

	w := bypack.NewWriter(cmd)
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			w.Int(int32(v))
		case string:
			w.String(string(v))
		}
	}
	w.End()
	return w.GetBuffer()
}
