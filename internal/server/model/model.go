package model

import (
	"Pchat/internal/server/conf"
	"time"

	"go.uber.org/zap"
)

func DataInfo() map[string]int32 {
	return map[string]int32{
		"conn_count": CountConn(),
		"mid_count":  CountMid(),
		"user_count": CountUser(),
	}
}

func RunTicker(done chan struct{}) {
	var tk = time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-tk.C:
			var Data = struct {
				UserCnt    int
				UserDelCnt int
			}{}

			var now = time.Now()

			mapUser.Range(func(key, value interface{}) bool {
				Data.UserCnt++
				user := value.(*User)
				if user.LastTime.Add(10 * time.Minute).Before(now) { // 10分钟没有更新
					Data.UserDelCnt++
					DelUser(user.Mid)
				}
				return true
			})
			conf.L.Info("online info",
				zap.Int("user_cnt", Data.UserCnt),
				zap.Int("user_del_cnt", Data.UserDelCnt),
			)

		case <-done:
			conf.L.Info("ticker stop")
			tk.Stop()
			return
		}
	}
}
