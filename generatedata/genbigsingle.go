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

type GenBigKVSingle struct {
	RedisConn       *redis.Conn
	KeySuffix       string
	Length          int //set、hash等容器型key的长度
	EXPIRE          time.Duration
	DB              int
	ValuePrefix     string
	DataGenInterval time.Duration
}

func (gbkv *GenBigKVSingle) GenBigHash() string {
	t1 := time.Now()
	key := "BigHash_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		gbkv.RedisConn.HSet(key, key+strconv.Itoa(i), gbkv.ValuePrefix+strconv.Itoa(i))
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisConn.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "hash"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))

	return key
}

func (gbkv *GenBigKVSingle) GenBigList() string {
	t1 := time.Now()
	key := "BigList_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		gbkv.RedisConn.LPush(key, gbkv.ValuePrefix+strconv.Itoa(i))
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisConn.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "list"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))
	return key
}

func (gbkv *GenBigKVSingle) GenBigSet() string {
	t1 := time.Now()
	key := "BigSet_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		gbkv.RedisConn.SAdd(key, gbkv.ValuePrefix+strconv.Itoa(i))
		time.Sleep(gbkv.DataGenInterval)
	}

	gbkv.RedisConn.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "set"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))
	return key
}
func (gbkv *GenBigKVSingle) GenBigZset() string {
	t1 := time.Now()
	key := "BigZset_" + gbkv.KeySuffix
	for i := 0; i < gbkv.Length; i++ {
		member := &redis.Z{
			Score:  rand.Float64(),
			Member: gbkv.ValuePrefix + strconv.Itoa(i),
		}
		gbkv.RedisConn.ZAdd(key, member)
		time.Sleep(gbkv.DataGenInterval)
	}
	gbkv.RedisConn.Expire(key, gbkv.EXPIRE)
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "zset"), zap.String("key", key), zap.String("duration", t2.Sub(t1).String()))
	return key
}

func (gbkv *GenBigKVSingle) GenBigString() {
	t1 := time.Now()
	key := "BigString_" + gbkv.KeySuffix
	for i := 0; i < 100; i++ {
		gbkv.RedisConn.Set(key+strconv.Itoa(i), gbkv.ValuePrefix+strconv.Itoa(i), gbkv.EXPIRE)
		time.Sleep(gbkv.DataGenInterval)
	}
	t2 := time.Now()
	global.RSPLog.Info("GenBigKV", zap.Int("db", gbkv.DB), zap.String("keytype", "string"), zap.String("keyprefix", key), zap.String("duration", t2.Sub(t1).String()))
}

// GenBigSingleExec 顺序执行所有命令，生成各个类型的大key
func (gbkv *GenBigKVSingle) GenBigSingleExec() {
	gbkv.GenBigHash()
	gbkv.GenBigList()
	gbkv.GenBigSet()
	gbkv.GenBigZset()
	gbkv.GenBigString()
}

func (gbkv *GenBigKVSingle) KeepGenBigSingle(ctx context.Context) {
	i:=int64(0)
	keySuffix:=gbkv.KeySuffix
	for {
		gbkv.KeySuffix=keySuffix+"_"+strconv.FormatInt(i,10)
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

func (gbkv *GenBigKVSingle) GenerateBaseDataParallel(client *redis.Client) map[string]string {

	global.RSPLog.Sugar().Info("Generate Base data Beging...")
	bigkvmap := make(map[string]string)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		newgbkv := new(GenBigKV)
		newgbkv.RedisConn = client.Conn()
		newgbkv.KeySuffix = gbkv.KeySuffix
		newgbkv.Length = gbkv.Length
		newgbkv.EXPIRE = gbkv.EXPIRE
		newgbkv.DB = gbkv.DB
		newgbkv.ValuePrefix = gbkv.ValuePrefix
		bigkvmap[newgbkv.GenBigHash()] = "Hash"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		newgbkv := new(GenBigKV)
		newgbkv.RedisConn = client.Conn()
		newgbkv.KeySuffix = gbkv.KeySuffix
		newgbkv.Length = gbkv.Length
		newgbkv.EXPIRE = gbkv.EXPIRE
		newgbkv.ValuePrefix = gbkv.ValuePrefix
		newgbkv.DB = gbkv.DB
		bigkvmap[newgbkv.GenBigList()] = "List"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		newgbkv := new(GenBigKV)
		newgbkv.RedisConn = client.Conn()
		newgbkv.KeySuffix = gbkv.KeySuffix
		newgbkv.Length = gbkv.Length
		newgbkv.EXPIRE = gbkv.EXPIRE
		newgbkv.ValuePrefix = gbkv.ValuePrefix
		newgbkv.DB = gbkv.DB
		bigkvmap[newgbkv.GenBigSet()] = "Set"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		newgbkv := new(GenBigKV)
		newgbkv.RedisConn = client.Conn()
		newgbkv.KeySuffix = gbkv.KeySuffix
		newgbkv.Length = gbkv.Length
		newgbkv.EXPIRE = gbkv.EXPIRE
		newgbkv.ValuePrefix = gbkv.ValuePrefix
		newgbkv.DB = gbkv.DB
		bigkvmap[newgbkv.GenBigZset()] = "Zset"
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		newgbkv := new(GenBigKV)
		newgbkv.RedisConn = client.Conn()
		newgbkv.KeySuffix = gbkv.KeySuffix
		newgbkv.Length = gbkv.Length
		newgbkv.EXPIRE = gbkv.EXPIRE
		newgbkv.ValuePrefix = gbkv.ValuePrefix
		newgbkv.DB = gbkv.DB
		newgbkv.GenBigString()
		wg.Done()
	}()
	wg.Wait()
	return bigkvmap
}
