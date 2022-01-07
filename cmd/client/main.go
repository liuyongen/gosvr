package main

import (
	"Pchat/internal/client"
	"Pchat/internal/client/conf"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var (
	ConfFile string
	StartAt  = time.Now()
)

func init() {
	rand.Seed(time.Now().Unix())
	flag.StringVar(&ConfFile, "conf", "app.toml", "conf file")
}

func main() {
	flag.Parse()
	fmt.Println("Pchat client go")

	conf.Init(ConfFile)

	client.Run()
	for {
	}
}
