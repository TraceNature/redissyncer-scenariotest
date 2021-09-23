package generatedata

import (
	"context"
	"github.com/go-redis/redis/v7"
	"runtime"
	"sync"
	"time"

	"testcase/commons"
	"testcase/global"
)

type TargetType int

const (
	TargettypeSingle  TargetType = 0
	TargettypeCluster TargetType = 1
)

func (tt TargetType) String() string {
	switch tt {
	case TargettypeSingle:
		return "single"
	case TargettypeCluster:
		return "cluster"
	default:
		return ""
	}
}

type BigKey struct {
	KeySuffixLen    int   `mapstructure:"keysuffixlen" json:"keysuffixlen" yaml:"keysuffixlen"`
	Length          int   `mapstructure:"length" json:"length" yaml:"length"`
	ValueSize       int   `mapstructure:"valuesize" json:"valuesize" yaml:"valuesize"`
	Expire          int64 `mapstructure:"expire" json:"expire" yaml:"expire"`
	Duration        int64 `mapstructure:"duaration" json:"duaration" yaml:"duaration"`
	DataGenInterval int64 `mapstructure:"datageninterval" json:"datageninterval" yaml:"datageninterval"`
}

type RandKey struct {
	KeySuffixLen int `mapstructure:"keysuffixlen" json:"keysuffixlen" yaml:"keysuffixlen"`
	//ValueSize       int   `mapstructure:"valuesize" json:"valuesize" yaml:"valuesize"`
	Loopstep        int   `mapstructure:"loopstep" json:"loopstep" yaml:"loopstep"`
	Expire          int64 `mapstructure:"expire" json:"expire" yaml:"expire"`
	Duration        int64 `mapstructure:"duaration" json:"duaration" yaml:"duaration"`
	DataGenInterval int64 `mapstructure:"datageninterval" json:"datageninterval" yaml:"datageninterval"`
	Threads         int   `mapstructure:"threads" json:"threads" yaml:"threads"`
}

type GenData struct {
	TargetType TargetType `mapstructure:"type" json:"type" yaml:"type"`
	Addr       []string   `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password   string     `mapstructure:"password" json:"password" yaml:"password"`
	DB         int        `mapstructure:"db" json:"db" yaml:"db"`
	BigKey     *BigKey    `mapstructure:"bigkey" json:"bigkey" yaml:"bigkey"`
	RandKey    *RandKey   `mapstructure:"randkey" json:"randkey" yaml:"randkey"`
}

func (gd *GenData) Exec() {
	global.RSPLog.Sugar().Info("GenData execute")
	if gd.TargetType == TargettypeSingle {
		redisaddr := gd.Addr[0]
		passwd := gd.Password
		redisopt := &redis.Options{
			Addr: redisaddr,
			DB:   gd.DB, // use default DB
		}

		if passwd != "" {
			redisopt.Password = gd.Password
		}

		client := commons.GetGoRedisClient(redisopt)

		_, err := client.Ping().Result()

		if err != nil {
			global.RSPLog.Sugar().Error(err)
			return
		}

		if gd.BigKey != nil {
			d := time.Now().Add(time.Duration(gd.BigKey.Duration) * time.Second)
			ctx, cancel := context.WithDeadline(context.Background(), d)
			defer cancel()

			wg := sync.WaitGroup{}
			keySuffix := commons.RandString(gd.BigKey.KeySuffixLen)
			valuePrefix := commons.RandString(gd.BigKey.ValueSize)
			genBigSingle := GenBigKVSingle{
				RedisConn:       client.Conn(),
				KeySuffix:       keySuffix,
				Length:          gd.BigKey.Length,
				EXPIRE:          time.Duration(gd.BigKey.Expire) * time.Second,
				DB:              gd.DB,
				ValuePrefix:     valuePrefix,
				DataGenInterval: time.Duration(gd.BigKey.DataGenInterval) * time.Millisecond,
			}

			wg.Add(1)
			go func() {
				genBigSingle.KeepGenBigSingle(ctx)
				wg.Done()
			}()
			wg.Wait()
		}

		if gd.RandKey != nil {
			threads := runtime.NumCPU()
			if gd.RandKey.Threads > 0 {
				threads = gd.RandKey.Threads
			}

			d := time.Now().Add(time.Duration(gd.RandKey.Duration) * time.Second)
			ctx, cancel := context.WithDeadline(context.Background(), d)
			defer cancel()

			wg := sync.WaitGroup{}
			keySuffix := commons.RandString(gd.RandKey.KeySuffixLen)
			optSingle := OptSingle{
				RedisConn: client.Conn(),
				KeySuffix: keySuffix,
				Loopstep:  gd.RandKey.Loopstep,
				EXPIRE:    time.Duration(gd.RandKey.Expire) * time.Second,
				DB:        gd.DB,
			}

			for i := 0; i < threads; i++ {
				wg.Add(1)
				go func() {
					optSingle.KeepExecBasicOpt(ctx, time.Duration(gd.RandKey.DataGenInterval)*time.Millisecond, false)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}

	if gd.TargetType == TargettypeCluster {
		opt := &redis.ClusterOptions{
			Addrs: gd.Addr,
		}
		if gd.Password != "" {
			opt.Password = gd.Password
		}
		clusterClient := redis.NewClusterClient(opt)

		if gd.BigKey != nil {
			d := time.Now().Add(time.Duration(gd.BigKey.Duration) * time.Second)
			ctx, cancel := context.WithDeadline(context.Background(), d)
			defer cancel()

			wg := sync.WaitGroup{}
			keySuffix := commons.RandString(gd.BigKey.KeySuffixLen)
			valuePrefix := commons.RandString(gd.BigKey.ValueSize)
			genBigKVCluster := GenBigKVCluster{
				RedisClusterClient: clusterClient,
				KeySuffix:          keySuffix,
				Length:             gd.BigKey.Length,
				EXPIRE:             time.Duration(gd.BigKey.Expire) * time.Second,
				DB:                 gd.DB,
				ValuePrefix:        valuePrefix,
				DataGenInterval:    time.Duration(gd.BigKey.DataGenInterval) * time.Millisecond,
			}

			wg.Add(1)
			go func() {
				genBigKVCluster.KeepGenBigSingle(ctx)
				wg.Done()
			}()
			wg.Wait()
		}

		if gd.RandKey != nil {
			threads := runtime.NumCPU()
			if gd.RandKey.Threads > 0 {
				threads = gd.RandKey.Threads
			}

			d := time.Now().Add(time.Duration(gd.RandKey.Duration) * time.Second)
			ctx, cancel := context.WithDeadline(context.Background(), d)
			defer cancel()

			wg := sync.WaitGroup{}
			keySuffix := commons.RandString(gd.RandKey.KeySuffixLen)
			optCluster := OptCluster{
				ClusterClient: clusterClient,
				KeySuffix:     keySuffix,
				Loopstep:      gd.RandKey.Loopstep,
				EXPIRE:        time.Duration(gd.RandKey.Expire) * time.Second,
				DB:            0,
			}

			for i := 0; i < threads; i++ {
				wg.Add(1)
				go func() {
					optCluster.KeepExecBasicOptCluster(ctx, time.Duration(gd.RandKey.DataGenInterval)*time.Millisecond)
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}
}