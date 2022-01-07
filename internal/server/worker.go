package server

import (
	"Pchat/internal/server/conf"
	"Pchat/internal/server/model"
	"fmt"
	"go.uber.org/zap"
	"net"

	"gogit.oa.com/March/gopkg/util"

	"gogit.oa.com/March/gopkg/protocol/bypack"
)

const (
	OFFLINE_MAX = 100 //最大离线数
)

//本地后台命令字
var LocalCmds = []string{
	"TCP0x881",
	"TCP0x888",
	"UDP0x103",
	"UDP0x104",
	"TCP0x109",
	"UDP0x10E",
	"UDP0x10F",
	"TCP0x110",
	"TCP0x887",
	"TCP0x886",
}

//客户端命令字
var LoginCmds = []string{
	"TCP0x2",   //Heartbeat
	"TCP0x103", //Chat
	"TCP0x301", //Group chat
}

type Worker struct {
	conn   net.Conn
	addr   net.Addr
	reader *bypack.Reader
	user   *model.User
}

func NewWorker(addr net.Addr, reader *bypack.Reader) *Worker {
	return &Worker{
		addr:   addr,
		reader: reader,
	}
}

func NewWorkerWithConn(conn net.Conn, reader *bypack.Reader) *Worker {
	return &Worker{
		conn:   conn,
		reader: reader,
		addr:   conn.RemoteAddr(),
	}
}

func (w *Worker) PreAction(method string) error {
	for _, v := range LocalCmds {
		if v == method {
			if !w.IsLocal() {
				return fmt.Errorf("only for local")
			}
		}
	}

	for _, v := range LoginCmds {
		if v == method {
			if !w.IsLogin() {
				return fmt.Errorf("need login before")
			}
		}
	}

	return nil
}

func (w *Worker) IsLogin() (ok bool) {
	w.user, ok = model.GetUserWithConn(w.conn)
	return ok
}

func (w *Worker) IsLocal() bool {
	if w.addr == nil {
		return false
	}
	return util.IsLocal(w.addr)
}

func (w *Worker) send(data []byte) {
	send(w.conn, data)
}

// 接收客户端发心跳 10s
func (w *Worker) TCP0x2() error {
	w.user.Heartbeat()
	w.send(model.BufferHeartbeat)
	return nil
}

// 登录
func (w *Worker) TCP0x101() error {
	var mid = w.reader.Int()
	if mid <= 0 {
		return fmt.Errorf("invalid mid %d login", mid)
	}
	conf.L.Info("user login", logField(w.conn), zap.Int32("mid", mid))

	//删除老的映射关系
	midTmp, ok := model.GetMid(w.conn)
	if ok && mid != midTmp {
		model.DelMid(w.conn)
		return fmt.Errorf("conn is aleady login")
	}

	model.AddMid(w.conn, mid)

	user, ok := model.GetUser(mid)
	if !ok {
		user = model.NewUser(w.conn, mid)
	}

	user.Conn = w.conn
	user.Mid = mid

	w.send(model.BufferLoginSuccess)

	//发送离线数据(最多一百条)
	num := 0
	for num < OFFLINE_MAX {
		offBuf := model.GetOffBuffer(mid)
		if offBuf != nil {
			w.send(offBuf)
		}
		num++
	}

	conf.L.Info("login success", logField(w.conn), zap.Int32("mid", mid))
	return nil
}

// TCP私聊
func (w *Worker) TCP0x103() error {
	return w.x103()
}

// UDP私聊
func (w *Worker) UDP0x103() error {
	return w.x103()
}

func (w *Worker) x103() error {
	fMid := w.reader.Int()
	toMid := w.reader.Int()

	buff := w.reader.RawBuffer
	my, ok := model.GetUser(fMid)
	my.Send(buff)

	user, ok := model.GetUser(toMid)
	if !ok || user == nil {
		//存储离线
		model.SetOffBuffer(toMid, w.reader.RawBuffer)
		return nil
		// return fmt.Errorf("user %d not found", toMid)
	}

	conf.L.Debug(fmt.Sprintf("fmid %d send to mid: %d data: %s", fMid, toMid, w.reader.String()), logField(w.conn))
	user.Send(buff)

	return nil
}

//TCP群聊
func (w *Worker) TCP0x301() error {
	return w.x301()
}

//UDP群聊
func (w *Worker) UDP0x301() error {
	return w.x301()
}

func (w *Worker) x301() error {

	mid := w.reader.Int()
	msg := w.reader.String()
	sendTime := w.reader.Int()
	name := w.reader.String()
	avatar := w.reader.String()
	msgType := w.reader.Int()

	user, ok := model.GetUser(mid)
	if !ok {
		return nil
	}

	conf.L.Info("param", zap.Any("mid", mid),
		zap.Any("msg", msg),
		zap.Any("sendTime", sendTime),
		zap.Any("name", name),
		zap.Any("avatar", avatar),
		zap.Any("msgType", msgType))

	wb := bypack.NewWriter(0x301)
	wb.Int(mid)
	wb.String(msg)
	wb.Int(sendTime)
	wb.String(name)
	wb.String(avatar)
	wb.Int(msgType)

	// banned
	isBan := user.Ban(5)
	conf.L.Info("Ban", zap.Any("mid", mid), zap.Any("isBan", isBan))
	if isBan == true {
		wb.Int(0)
		wb.Int(-1)
		wb.Int(-1)
		wb.End()
		user.Send(wb.GetBuffer())
		return nil
	}

	//limited
	remainTime := user.RemainTimes()
	conf.L.Info("Remain Times", zap.Any("mid", mid), zap.Any("remainTime", remainTime))
	if remainTime == 0 {
		wb.Int(0)
		wb.Int(remainTime)
		wb.Int(-2)
		wb.End()
		user.Send(wb.GetBuffer())
		return nil
	}

	//Countdown
	countDown := user.CountDown()
	conf.L.Info("Countdown", zap.Any("mid", mid), zap.Any("countDown", countDown))
	if countDown > 0 && countDown < 60 {
		wb.Int(countDown)
		wb.Int(-1)
		wb.Int(-3)
		wb.End()
		user.Send(wb.GetBuffer())
		return nil
	}

	//self
	wb.Int(countDown)
	wb.Int(remainTime)
	wb.Int(0)
	wb.End()
	user.Send(wb.GetBuffer())

	//broadcast
	buff := w.reader.RawBuffer
	users := model.GetUsers()
	for _, v := range users {
		//filter
		if v.Mid == mid {
			continue
		}
		go v.Send(buff)
	}

	return nil
}
