package server

import (
	"Pchat/internal/server/conf"
	"Pchat/internal/server/model"
	"encoding/json"
	"fmt"
)

// 监控测试
func (w *Worker) TCP0x881() error {
	w.send(model.BufferMonitorTest)
	return nil
}

// 重启 FIXME 不实现
func (w *Worker) TCP0x882() error {
	return nil
}

// 系统信息
func (w *Worker) TCP0x888() error {
	key := w.reader.String()
	// fmt.Println("key:", key)
	// fmt.Println("key1:", conf.Conf.Keys.Tcp0x888)
	if key != conf.Conf.Keys.Tcp0x888 {
		_ = w.conn.Close()
		return fmt.Errorf("invalid key %s", key)
	}

	var data = map[string]interface{}{
		"tableInfo": model.DataInfo(),
		"version":   conf.Version,
		"startAt":   conf.StartAt.Format("2006-01-02 15:04:05"),
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.send(model.GetBufferAdmin(string(b)))

	return nil
}
