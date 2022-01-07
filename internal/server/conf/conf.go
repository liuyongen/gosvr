package conf

import (
	"path/filepath"
	"time"

	"gogit.oa.com/March/gopkg/build"

	"gogit.oa.com/March/gopkg/metric"

	"go.uber.org/zap"
	"gogit.oa.com/March/gopkg/logger"

	"github.com/BurntSushi/toml"

	"context"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v8"
	"gogit.oa.com/March/gopkg/util"
)

var (
	Version string
	StartAt time.Time
	Redis   *redis.Client
	Mem     *memcache.Client
	L       *zap.Logger
)

var Conf = struct {
	App struct {
		Name     string `toml:"name"`
		Env      string `toml:"env"`
		HttpAddr string `toml:"http_addr"`
	}
	Log struct {
		Path   string
		Size   int
		Backup int
		Age    int
	}
	Redis struct {
		Addr string
	}
	Memcached struct {
		Addr string
	}
	Server struct {
		TcpAddr string `toml:"tcp_addr"`
		UdpAddr string `toml:"udp_addr"`
	}
	Keys struct {
		Tcp0x888 string `toml:"tcp0x888"`
		Udp0x10e string `toml:"udp0x10e"`
		Udp0x10f string `toml:"udp0x10f"`
	}
}{}

func Init(filename string) {
	StartAt = time.Now()

	//解析配置
	pth, err := filepath.Abs(filename)
	util.MustNil(err)

	_, err = toml.DecodeFile(pth, &Conf)
	util.MustNil(err)

	//启动信息
	initVersion()

	//Redis缓存
	Redis = redis.NewClient(&redis.Options{Addr: Conf.Redis.Addr})
	util.MustNil(Redis.Ping(context.Background()).Err())

	//Mem缓存
	Mem := memcache.New(Conf.Memcached.Addr)
	err = Mem.Ping()
	util.MustNil(err)

	//文件日志
	p, err := filepath.Abs(Conf.Log.Path)
	util.MustNil(err)
	L = logger.NewFileLogger(p, Conf.App.Name, Conf.Log.Size, Conf.Log.Backup, Conf.Log.Age)
	L.Info("config init", zap.Any("value", Conf))

	metric.InitMetricsWithCode(Conf.App.Name, []string{"startedCounter", "handledCounter", "handledHistogram"})
}

func initVersion() {
	Version = build.Print(Conf.App.Name)
}
