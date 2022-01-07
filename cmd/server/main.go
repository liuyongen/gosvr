package main

import (
	"Pchat/internal/server"
	"Pchat/internal/server/conf"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	ConfFile string
	Version  string
)

func init() {
	rand.Seed(time.Now().Unix())
	flag.StringVar(&ConfFile, "conf", "app.toml", "conf file")
	flag.StringVar(&Version, "v", "0.0.1", "version") // 为了命令行启动参数
}

func main() {
	flag.Parse()

	conf.Init(ConfFile)
	fmt.Println(conf.Version)

	done := make(chan struct{}, 1)
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sign := <-signs
		conf.L.Info("catch os signal", zap.String("value", sign.String()))

		server.Close()
		close(done)
	}()

	go server.Run()

	<-done
	conf.L.Info("exiting...")
}
