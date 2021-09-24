package generatedata

import (
	"context"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"sync"
	"testcase/global"
	"time"
)

type GenBigKVCluster struct {
	RedisClusterClient *redis.ClusterClient
	KeySuffix          string
	Length             int
	EXPIRE             time.Duration
	DB                 int
	ValuePrefix        string
	DataGenInterval    time.Duration
}

func (gbkv *GenBigKVCluster) GenBigHash() string {
	t1 := time.Now()
	key := "BigHash_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		gbkv.RedisClusterClient.HSet(key, key+strconv.Itoa(i), gbkv.ValuePrefix+strconv.Itoa(i))
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisClusterClient.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()

	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "hash"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))
	return key
}

func (gbkv *GenBigKVCluster) GenBigList() string {
	t1 := time.Now()
	key := "BigList_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		gbkv.RedisClusterClient.LPush(key, gbkv.ValuePrefix+strconv.Itoa(i))
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisClusterClient.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "list"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))

	return key
}

func (gbkv *GenBigKVCluster) GenBigSet() string {
	t1 := time.Now()
	key := "BigSet_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		gbkv.RedisClusterClient.SAdd(key, gbkv.ValuePrefix+strconv.Itoa(i))
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisClusterClient.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "set"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))
	return key
}

func (gbkv *GenBigKVCluster) GenBigZset() string {
	t1 := time.Now()
	key := "BigZset_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		member := &redis.Z{
			Score:  rand.Float64(),
			Member: gbkv.ValuePrefix + strconv.Itoa(i),
		}
		gbkv.RedisClusterClient.ZAdd(key, member)
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisClusterClient.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "zset"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))
	return key
}

func (gbkv *GenBigKVCluster) GenBigString() {
	t1 := time.Now()
	key := "BigString_" + gbkv.KeySuffix
	for i := 0; i < 100; i++ {
		gbkv.RedisClusterClient.Set(key+strconv.Itoa(i), gbkv.ValuePrefix+strconv.Itoa(i), gbkv.EXPIRE)
		time.Sleep(gbkv.DataGenInterval)
	}
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "string"), zap.String("keyprefix", key), zap.String("duration", t2.Sub(t1).String()))
}

func (gbkv *GenBigKVCluster) GenBigSingleExec() {
	gbkv.GenBigHash()
	gbkv.GenBigList()
	gbkv.GenBigSet()
	gbkv.GenBigZset()
	gbkv.GenBigString()
}

func (gbkv *GenBigKVCluster) KeepGenBigSingle(ctx context.Context) {
	i := int64(0)
	keySuffix := gbkv.KeySuffix
	for {
		gbkv.KeySuffix = keySuffix + "_" + strconv.FormatInt(i, 10)
		gbkv.GenBigSingleExec()
		i++
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (gbkv *GenBigKVCluster) GenerateBaseDataParallelCluster() map[string]string {

	zaplogger.Sugar().Info("Generate Base data Beging...")
	bigkvmap := make(map[string]string)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		bigkvmap[gbkv.GenBigHash()] = "Hash"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		bigkvmap[gbkv.GenBigList()] = "List"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		bigkvmap[gbkv.GenBigSet()] = "Set"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		bigkvmap[gbkv.GenBigZset()] = "Zset"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		gbkv.GenBigString()
		wg.Done()
	}()
	wg.Wait()
	return bigkvmap
}
