package model

import (
	"Pchat/internal/server/conf"
	"Pchat/pkg/util"
	ctx "context"
	"fmt"
	"github.com/buger/jsonparser"
	"go.uber.org/zap"
	"net"
	"time"
)

type User struct {
	Conn     net.Conn
	Mid      int32
	LastTime time.Time
}

func NewUser(conn net.Conn, mid int32) *User {
	conf.L.Info("new user", zap.Int32("mid", mid))
	user := &User{
		Conn:     conn,
		Mid:      mid,
		LastTime: time.Now(),
	}
	AddUser(user)
	return user
}

func (u *User) LogField() zap.Field {
	return zap.String("mid", fmt.Sprintf("%d[%s]", u.Mid, u.Conn.RemoteAddr().String()))
}

func (u *User) Send(data []byte) {
	if u.Conn == nil {
		return
	}

	if err := u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		conf.L.Warn(err.Error(), u.LogField())
		return
	}

	if _, err := u.Conn.Write(data); err != nil {
		conf.L.Warn(err.Error(), u.LogField())
		return
	}
}

func (u *User) Heartbeat() {
	now := time.Now()
	u.LastTime = now
}

//禁言
func (u *User) Ban(t int8) bool {

	key := GetBanKey(u.Mid)
	ban := conf.Redis.Get(ctx.Background(), key).Val()
	if ban == "" {
		return false
	}
	banArr, err := util.Unserialize(ban)
	if err != nil {
		conf.L.Warn(err.Error(), u.LogField())
	}
	return util.InArray(t, banArr)

}

//倒计时
func (u *User) CountDown() int32 {

	key := GetCDKey(u.Mid)
	exist := conf.Redis.Exists(ctx.Background(), key).Val()
	if exist == 0 {
		conf.Redis.SetEX(ctx.Background(), key, 1, 60*time.Second)
	}

	time := conf.Redis.TTL(ctx.Background(), key).Val().Seconds()

	return int32(time)
}

//次数
func (u *User) RemainTimes() int32 {

	key := GetTimesKey(u.Mid)
	useTimes := conf.Redis.Incr(ctx.Background(), key).Val()
	//一天过期
	now := util.Time()
	yesDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	yesTime, _ := util.Strtotime("2006-01-02", yesDate)
	duration := yesTime - now
	conf.Redis.Expire(ctx.Background(), key, time.Duration(duration)*time.Second)
	//vip等级
	vipKey := GetVipKey(u.Mid)
	vip, _ := conf.Redis.Get(ctx.Background(), vipKey).Int()
	//规则
	ruleKey := GetRuleKey()
	rule := conf.Redis.Get(ctx.Background(), ruleKey).Val()

	idx := fmt.Sprintf("[%d]", vip)
	vipTimes, err := jsonparser.GetInt([]byte(rule), "vip_times", idx)
	if err != nil {
		return -1
	}
	diff := vipTimes - useTimes
	if diff > 0 {
		return int32(diff)
	}

	return 0

}
