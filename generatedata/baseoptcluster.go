//Package generatedata 用于生成测试过程中的数据

package generatedata

import (
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"context"
)

type OptCluster struct {
	ClusterClient *redis.ClusterClient
	RedisVersion  string
	OptType       OptType
	KeySuffix     string
	Loopstep      int
	EXPIRE        time.Duration
	DB            int
}

func (optcluster *OptCluster) ExecOpt() {

	switch optcluster.OptType.String() {
	case "BO_APPEND":
		optcluster.BO_APPEND()
	case "BO_BITOP":
		optcluster.BO_BITOP()
	case "BO_DECR_DECRBY":
		optcluster.BO_DECR_DECRBY()
	case "BO_INCR_INCRBY_INCRBYFLOAT":
		optcluster.BO_INCR_INCRBY_INCRBYFLOAT()
	case "BO_MSET_MSETNX":
		optcluster.BO_MSET_MSETNX()
	//case "BO_PSETEX_SETEX":
	//	optcluster.BO_PSETEX_SETEX()
	case "BO_PFADD":
		optcluster.BO_PFADD()
	case "BO_SET_SETNX":
		optcluster.BO_SET_SETNX()
	case "BO_SETBIT":
		optcluster.BO_SETBIT()
	case "BO_SETRANGE":
		optcluster.BO_SETRANGE()
	case "BO_HINCRBY_HINCRBYFLOAT":
		optcluster.BO_HINCRBY_HINCRBYFLOAT()
	case "BO_HSET_HMSET_HSETNX":
		optcluster.BO_HSET_HMSET_HSETNX()
	case "BO_LPUSH_LPOP_LPUSHX":
		optcluster.BO_LPUSH_LPOP_LPUSHX()
	case "BO_LREM_LTRIM_LINSERT":
		optcluster.BO_LREM_LTRIM_LINSERT()
	case "BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH":
		optcluster.BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH()
	case "BO_BLPOP_BRPOP_BRPOPLPUSH":
		optcluster.BO_BLPOP_BRPOP_BRPOPLPUSH()
	case "BO_SADD_SMOVE_SPOP_SREM":
		optcluster.BO_SADD_SMOVE_SPOP_SREM()
	case "BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE":
		optcluster.BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE()
	case "BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM":
		optcluster.BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM()
	case "BO_ZPOPMAX_ZPOPMIN":
		optcluster.BO_ZPOPMAX_ZPOPMIN()
	case "BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE":
		optcluster.BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE()
	case "BO__ZUNIONSTORE_ZINTERSTORE":
		optcluster.BO_ZUNIONSTORE_ZINTERSTORE()
	default:
		return
	}

}

// 比较目标库版本是否小于要求版本
func (optcluster OptCluster) VersionLessThan(version string) bool {

	boverarray := strings.Split(optcluster.RedisVersion, ".")
	versionarry := strings.Split(version, ".")

	bover := ""
	ver := ""
	for i := 0; i < 3; i++ {
		if i < len(boverarray) {
			bover = bover + boverarray[i]
		} else {
			bover = bover + "0"
		}

		if i < len(ver) {
			ver = ver + versionarry[i]
		} else {
			ver = ver + "0"
		}
	}

	intbover, _ := strconv.Atoi(bover)
	intver, _ := strconv.Atoi(ver)

	if intbover < intver {
		return true
	}

	return false

}

//SELECT命令
//func (optcluster *OptCluster) BO_SELECT(db int) {
//
//	_, err := optcluster.ClusterClient.Select(db).Result()
//	if err != nil {
//		zaplogger.Sugar().Error(err)
//		return
//	}
//	optcluster.DB = db
//}

//APPEND 命令基本操作
//start version:2.0.0
func (optcluster *OptCluster) BO_APPEND() {
	t1 := time.Now()
	appended := "append_" + optcluster.KeySuffix
	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.Append(appended, strconv.Itoa(i))
	}
	optcluster.ClusterClient.Expire(appended, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "APPEND"), zap.String("key", appended), zap.Int64("time", t2.Sub(t1).Milliseconds()))

}

//BITOP
//start version:2.6.0
func (optcluster *OptCluster) BO_BITOP() {
	t1 := time.Now()
	strarry := []string{}
	opandkey := "opand_" + optcluster.KeySuffix
	oporkey := "opor_" + optcluster.KeySuffix
	opxorkey := "opxor_" + optcluster.KeySuffix
	opnotkey := "opnot_" + optcluster.KeySuffix
	for i := 0; i < optcluster.Loopstep; i++ {
		bitopkey := "bitop_" + optcluster.KeySuffix + strconv.Itoa(i)
		optcluster.ClusterClient.Set(bitopkey, bitopkey, optcluster.EXPIRE)
		strarry = append(strarry, bitopkey)
	}

	optcluster.ClusterClient.BitOpAnd(opandkey, strarry...)
	optcluster.ClusterClient.BitOpOr(oporkey, strarry...)
	optcluster.ClusterClient.BitOpXor(opxorkey, strarry...)
	optcluster.ClusterClient.BitOpNot(opnotkey, strarry[0])
	optcluster.ClusterClient.Expire(opandkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(oporkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(opxorkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(opnotkey, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BITOP"), zap.Any("keys", []string{opandkey, oporkey, opxorkey, opnotkey}), zap.Duration("time", t2.Sub(t1)))
}

//DECR and DECRBY
func (optcluster *OptCluster) BO_DECR_DECRBY() {
	t1 := time.Now()
	desckey := "desc_" + optcluster.KeySuffix
	optcluster.ClusterClient.Set(desckey, optcluster.Loopstep, optcluster.EXPIRE)
	optcluster.ClusterClient.Decr(desckey)
	optcluster.ClusterClient.DecrBy(desckey, rand.Int63n(int64(optcluster.Loopstep)))
	t2 := time.Now()

	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "DECR_DECRBY"), zap.String("key", desckey), zap.Duration("time", t2.Sub(t1)))
}

//INCR and INCRBY and INCRBYFLOAT
func (optcluster *OptCluster) BO_INCR_INCRBY_INCRBYFLOAT() {
	t1 := time.Now()
	incrkey := "incr_" + optcluster.KeySuffix
	optcluster.ClusterClient.Set(incrkey, rand.Intn(optcluster.Loopstep), optcluster.EXPIRE)
	optcluster.ClusterClient.Incr(incrkey)
	optcluster.ClusterClient.IncrBy(incrkey, rand.Int63n(int64(optcluster.Loopstep)))
	optcluster.ClusterClient.IncrByFloat(incrkey, rand.Float64()*float64(rand.Intn(optcluster.Loopstep)))
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "INCR_INCRBY_INCRBYFLOAT"), zap.String("key", incrkey), zap.Duration("time", t2.Sub(t1)))
}

//MSET and MSETNX
func (optcluster *OptCluster) BO_MSET_MSETNX() {
	t1 := time.Now()
	msetarry := []string{}
	msetnxarry := []string{}

	msetkey := "mset_" + optcluster.KeySuffix
	msetnxkey := "msetnx_" + optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		msetarry = append(msetarry, msetkey+strconv.Itoa(i))
		msetarry = append(msetarry, msetkey+strconv.Itoa(i))
		msetnxarry = append(msetnxarry, msetnxkey+strconv.Itoa(i))
		msetnxarry = append(msetnxarry, msetnxkey+strconv.Itoa(i))
	}

	optcluster.ClusterClient.MSetNX(msetnxarry)
	optcluster.ClusterClient.MSet(msetarry)
	optcluster.ClusterClient.MSetNX(msetnxarry)

	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.Expire(msetkey+strconv.Itoa(i), optcluster.EXPIRE)
		optcluster.ClusterClient.Expire(msetnxkey+strconv.Itoa(i), optcluster.EXPIRE)
	}

	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "MGET_MSETNX"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

////PSETEX and SETEX
//func (bo *OptCluster) BO_PSETEX_SETEX() {
//	t1 := time.Now()
//	psetexkey := "psetex_" + bo.KeySuffix
//	setexkey := "setex_" + bo.KeySuffix
//	bo.ClusterClient.SetNX(setexkey, setexkey, bo.EXPIRE)
//	bo.ClusterClient.
//	bo.ClusterClient.Do("SETEX", setexkey, bo.EXPIRE.Seconds(), setexkey)
//	bo.ClusterClient.Do("PSETEX", psetexkey, bo.EXPIRE.Milliseconds(), psetexkey)
//	t2 := time.Now()
//	zaplogger.Info("ExecCMD", zap.String("command", "MGET_MSETNX"), zap.String("KeySuffix", bo.KeySuffix), zap.Duration("time", t2.Sub(t1)))
//
//}

//PFADD
func (optcluster *OptCluster) BO_PFADD() {
	t1 := time.Now()
	pfaddkey := "pfadd_" + optcluster.KeySuffix
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.PFAdd(pfaddkey, rand.Float64()*float64(rand.Int()))
	}
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_PFADD"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//SET and SETNX
func (optcluster *OptCluster) BO_SET_SETNX() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	setkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "set_" + optcluster.KeySuffix
	setnxkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "setnx_" + optcluster.KeySuffix
	optcluster.ClusterClient.Set(setkey, setkey, optcluster.EXPIRE)
	optcluster.ClusterClient.SetNX(setnxkey, setkey, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "SET_SETNX"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//SETBIT
func (optcluster *OptCluster) BO_SETBIT() {
	t1 := time.Now()
	setbitkey := "setbit_" + optcluster.KeySuffix
	optcluster.ClusterClient.SetBit(setbitkey, rand.Int63n(int64(optcluster.Loopstep)), rand.Intn(optcluster.Loopstep))
	optcluster.ClusterClient.Expire(setbitkey, optcluster.EXPIRE)

	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "SETBIT"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//SETRANGE
func (optcluster *OptCluster) BO_SETRANGE() {
	t1 := time.Now()
	setrangekey := "setrange_" + optcluster.KeySuffix
	optcluster.ClusterClient.Set(setrangekey, setrangekey, optcluster.EXPIRE)
	optcluster.ClusterClient.SetRange(setrangekey, rand.Int63n(int64(optcluster.Loopstep)), strconv.Itoa(rand.Intn(optcluster.Loopstep)))
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "SETRANGE"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//HINCRBY and HINCRBYFLOAT
func (optcluster *OptCluster) BO_HINCRBY_HINCRBYFLOAT() {
	t1 := time.Now()
	hincrbykey := "hincrby_" + optcluster.KeySuffix
	hincrbyfloatkey := "hincrbyfloat_" + optcluster.KeySuffix
	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.HIncrBy(hincrbykey, hincrbykey+strconv.Itoa(rand.Intn(optcluster.Loopstep)), int64(rand.Intn(optcluster.Loopstep)))
		optcluster.ClusterClient.HIncrByFloat(hincrbyfloatkey, hincrbyfloatkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)), rand.Float64()*10)
	}
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "HINCRBY_HINCRBYFLOAT"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//HSET HMSET HSETNX HDEL
func (optcluster *OptCluster) BO_HSET_HMSET_HSETNX() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	hsetkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "hset_" + optcluster.KeySuffix
	hmsetkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "hmset_" + optcluster.KeySuffix
	fieldmap := make(map[string]interface{})

	for i := 0; i < optcluster.Loopstep; i++ {
		field := hmsetkey + strconv.Itoa(i)
		fieldmap[field] = field
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.HSet(hsetkey, hsetkey+strconv.Itoa(i), hsetkey+strconv.Itoa(i))
	}

	optcluster.ClusterClient.HMSet(hmsetkey, fieldmap)

	//HSETNX
	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.HSetNX(hmsetkey, hmsetkey+strconv.Itoa(rand.Intn(optcluster.Loopstep*2)), hmsetkey+strconv.Itoa(i))
	}

	//HDEL
	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.HDel(hmsetkey, hmsetkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)))
		}
	}

	optcluster.ClusterClient.Expire(hsetkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(hmsetkey, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "HSET_HMSET_HSETNX_HDEL"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//LPUSH and LPOP and LPUSHX and LSET
func (optcluster *OptCluster) BO_LPUSH_LPOP_LPUSHX() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	lpushkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "lpush_" + optcluster.KeySuffix
	lpushxkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "lpushx_" + optcluster.KeySuffix
	values := make([]interface{}, optcluster.Loopstep)
	for i := 0; i < len(values); i++ {
		values[i] = lpushkey + strconv.Itoa(i)
	}

	optcluster.ClusterClient.LPush(lpushkey, values...)
	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 != 0 {
			optcluster.ClusterClient.LSet(lpushkey, int64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep), lpushkey+strconv.Itoa(i))
		}
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 != 0 {
			optcluster.ClusterClient.LPop(lpushkey)
		}
	}
	optcluster.ClusterClient.LPushX(lpushxkey, values)
	optcluster.ClusterClient.LPushX(lpushkey, values)

	optcluster.ClusterClient.Expire(lpushkey, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "LPUSH_LPOP_LPUSHX_LSET"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))

}

//LREM and LTRIM and LINSERT
func (optcluster *OptCluster) BO_LREM_LTRIM_LINSERT() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	lremkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "lrem_" + optcluster.KeySuffix
	ltrimkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "ltrim_" + optcluster.KeySuffix
	values := make([]interface{}, optcluster.Loopstep)
	for i := 0; i < len(values); i++ {
		values[i] = lremkey + strconv.Itoa(i)
	}
	optcluster.ClusterClient.LPush(lremkey, values...)
	optcluster.ClusterClient.LPush(ltrimkey, values...)

	for i := 0; i < optcluster.Loopstep; i++ {
		op := "BEFORE"
		if i%2 == 0 {
			op = "AFTER"
		} else {
			op = "BEFORE"
		}
		optcluster.ClusterClient.LInsert(lremkey, op, lremkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)), lremkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)))
	}

	optcluster.ClusterClient.LRem(lremkey, int64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep), lremkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)))
	optcluster.ClusterClient.LTrim(ltrimkey, int64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep), int64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep))
	optcluster.ClusterClient.Expire(lremkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(ltrimkey, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "LREM_TRIM_LINSERT"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//RPUSH RPUSHX RPOP RPOPLPUSH
func (optcluster *OptCluster) BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH() {
	t1 := time.Now()
	rpushkey := "rpush_" + optcluster.KeySuffix
	rpushxkey := "rpushx_" + optcluster.KeySuffix
	values := make([]interface{}, optcluster.Loopstep)
	for i := 0; i < len(values); i++ {
		values[i] = rpushkey + strconv.Itoa(i)
	}
	optcluster.ClusterClient.RPush(rpushkey, values...)
	optcluster.ClusterClient.RPushX(rpushxkey, values...)
	optcluster.ClusterClient.RPushX(rpushkey, values...)

	//rpoplpush 操作同一个key相当于将列表逆转
	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.RPopLPush(rpushkey, rpushkey)
		}
	}
	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.RPop(rpushkey)
		}
	}
	optcluster.ClusterClient.Expire(rpushkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(rpushxkey, optcluster.EXPIRE)

	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//BLPOP BRPOP BRPOPLPUSH BRPOPLPUSH
//BRPOPLPUSH 命令在集群模式下key分布在不同节点会报错(error) CROSSSLOT Keys in request don't hash to the same slot
func (optcluster *OptCluster) BO_BLPOP_BRPOP_BRPOPLPUSH() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	blpopkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "blpop_" + optcluster.KeySuffix
	brpopkey := "{" + strconv.Itoa(nodeprefix) + ":}" + "brpop_" + optcluster.KeySuffix

	values := make([]interface{}, optcluster.Loopstep)
	for i := 0; i < len(values); i++ {
		values[i] = blpopkey + strconv.Itoa(i)
	}

	optcluster.ClusterClient.RPush(blpopkey, values...)
	optcluster.ClusterClient.RPush(brpopkey, values...)

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.RPop(blpopkey)
			optcluster.ClusterClient.RPop(brpopkey)
		}
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.BRPopLPush(blpopkey, blpopkey, optcluster.EXPIRE)
		}
	}

	optcluster.ClusterClient.Expire(blpopkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(brpopkey, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_BLPOP_BRPOP_BRPOPLPUSH"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//SADD SMOVE SPOP SREM
func (optcluster *OptCluster) BO_SADD_SMOVE_SPOP_SREM() {
	t1 := time.Now()

	saddkey := "sadd_" + optcluster.KeySuffix
	smovekey := "smove_" + optcluster.KeySuffix
	spopkey := optcluster.KeySuffix
	sremkey := optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.SAdd(saddkey, saddkey+strconv.Itoa(i))
		optcluster.ClusterClient.SAdd(smovekey, smovekey+strconv.Itoa(i))
		optcluster.ClusterClient.SAdd(spopkey, spopkey+strconv.Itoa(i))
		optcluster.ClusterClient.SAdd(sremkey, sremkey+strconv.Itoa(i))
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.SPop(spopkey)
		}
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.SRem(sremkey, sremkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)))
		}
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.SMove(smovekey, smovekey, saddkey+strconv.Itoa(i))
		}
	}

	optcluster.ClusterClient.Del(saddkey)
	optcluster.ClusterClient.Expire(smovekey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(spopkey, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(sremkey, optcluster.EXPIRE)

	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))

}

//SDIFFSTORE SINTERSTORE SUNIONSTORE
func (optcluster *OptCluster) BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	sdiff1 := "{" + strconv.Itoa(nodeprefix) + ":}" + "sdiff1_" + optcluster.KeySuffix
	sdiff2 := "{" + strconv.Itoa(nodeprefix) + ":}" + "sdiff2_" + optcluster.KeySuffix
	sdiffstore := "{" + strconv.Itoa(nodeprefix) + ":}" + "sdiffsotre_" + optcluster.KeySuffix
	sinterstore := "{" + strconv.Itoa(nodeprefix) + ":}" + "sintersotre_" + optcluster.KeySuffix
	sunionstore := "{" + strconv.Itoa(nodeprefix) + ":}" + "sunionstore_" + optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		optcluster.ClusterClient.SAdd(sdiff1, optcluster.KeySuffix+strconv.Itoa(rand.Intn(2*optcluster.Loopstep)))
		optcluster.ClusterClient.SAdd(sdiff2, optcluster.KeySuffix+strconv.Itoa(rand.Intn(2*optcluster.Loopstep)))
	}

	optcluster.ClusterClient.SDiffStore(sdiffstore, sdiff1, sdiff2)
	optcluster.ClusterClient.SInterStore(sinterstore, sdiff1, sdiff2)
	optcluster.ClusterClient.SUnionStore(sunionstore, sdiff1, sdiff2)

	optcluster.ClusterClient.Expire(sdiffstore, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(sinterstore, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(sunionstore, optcluster.EXPIRE)

	optcluster.ClusterClient.Del(sdiff1, sdiff2)

	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//ZADD ZINCRBY ZREM
//start version:1.2.0
func (optcluster OptCluster) BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM() {
	t1 := time.Now()

	zaddkey := "zadd_" + optcluster.KeySuffix
	zincrby := "zincrby_" + optcluster.KeySuffix
	zrem := optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		z := redis.Z{
			Score:  float64(i),
			Member: zaddkey + strconv.Itoa(i),
		}
		optcluster.ClusterClient.ZAdd(zaddkey, &z)
		optcluster.ClusterClient.ZAdd(zincrby, &z)
		optcluster.ClusterClient.ZAdd(zrem, &z)

	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.ZIncrBy(zincrby, float64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep), zaddkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)))
			optcluster.ClusterClient.ZRem(zrem, zaddkey+strconv.Itoa(rand.Intn(optcluster.Loopstep)))
		}
	}

	optcluster.ClusterClient.Expire(zincrby, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(zrem, optcluster.EXPIRE)
	t2 := time.Now()

	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//ZPOPMAX ZPOPMIN
//start version:5.0
func (optcluster OptCluster) BO_ZPOPMAX_ZPOPMIN() {
	t1 := time.Now()

	zpopmax := "zpopmax_" + optcluster.KeySuffix
	zpopmin := "zpopmin_" + optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		z := redis.Z{
			Score:  float64(i),
			Member: zpopmax + strconv.Itoa(i),
		}
		optcluster.ClusterClient.ZAdd(zpopmax, &z)
		optcluster.ClusterClient.ZAdd(zpopmin, &z)
	}

	for i := 0; i < optcluster.Loopstep; i++ {
		if rand.Intn(optcluster.Loopstep)%2 == 0 {
			optcluster.ClusterClient.ZPopMax(zpopmax, 1)
			optcluster.ClusterClient.ZPopMin(zpopmin, 1)
		}
	}

	optcluster.ClusterClient.Expire(zpopmax, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(zpopmin, optcluster.EXPIRE)
	t2 := time.Now()

	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_ZPOPMAX_ZPOPMIN"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//ZREMRANGEBYLEX ZREMRANGEBYRANK ZREMRANGEBYSCORE
//start version:2.8.9
func (optcluster *OptCluster) BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	zremrangebylex := "{" + strconv.Itoa(nodeprefix) + ":}" + "zremrangebylex_" + optcluster.KeySuffix
	zremrangebyrank := "{" + strconv.Itoa(nodeprefix) + ":}" + "zremrangebyrank_" + optcluster.KeySuffix
	zremrangebyscore := "{" + strconv.Itoa(nodeprefix) + ":}" + "zremrangebyscore_" + optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		z := redis.Z{
			Score:  float64(i),
			Member: zremrangebylex + strconv.Itoa(i),
		}
		optcluster.ClusterClient.ZAdd(zremrangebylex, &z)
		optcluster.ClusterClient.ZAdd(zremrangebyrank, &z)
		optcluster.ClusterClient.ZAdd(zremrangebyscore, &z)
	}

	optcluster.ClusterClient.ZRemRangeByLex(zremrangebylex, zremrangebylex+strconv.Itoa(0), zremrangebylex+strconv.Itoa(rand.Intn(optcluster.Loopstep-1)))
	optcluster.ClusterClient.ZRemRangeByRank(zremrangebyrank, int64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep), int64(rand.Intn(2*optcluster.Loopstep)-optcluster.Loopstep))
	optcluster.ClusterClient.ZRemRangeByScore(zremrangebyscore, strconv.Itoa(rand.Intn(optcluster.Loopstep)), strconv.Itoa(rand.Intn(optcluster.Loopstep)))

	optcluster.ClusterClient.Expire(zremrangebylex, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(zremrangebyrank, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(zremrangebyscore, optcluster.EXPIRE)

	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))

}

//ZUNIONSTORE ZINTERSTORE
func (optcluster *OptCluster) BO_ZUNIONSTORE_ZINTERSTORE() {
	t1 := time.Now()
	nodeprefix := rand.Int()
	zset1 := "{" + strconv.Itoa(nodeprefix) + ":}" + "zset1_" + optcluster.KeySuffix
	zset2 := "{" + strconv.Itoa(nodeprefix) + ":}" + "zset2_" + optcluster.KeySuffix
	zset3 := "{" + strconv.Itoa(nodeprefix) + ":}" + "zset3_" + optcluster.KeySuffix
	zinterstore := "{" + strconv.Itoa(nodeprefix) + ":}" + "zinterstore_" + optcluster.KeySuffix
	zunionstore := "{" + strconv.Itoa(nodeprefix) + ":}" + "zunionstore_" + optcluster.KeySuffix

	for i := 0; i < optcluster.Loopstep; i++ {
		z := redis.Z{
			Score:  float64(i),
			Member: zset1 + strconv.Itoa(i),
		}
		optcluster.ClusterClient.ZAdd(zset1, &z)
		optcluster.ClusterClient.ZAdd(zset2, &z)
		optcluster.ClusterClient.ZAdd(zset3, &z)
	}

	zstore := redis.ZStore{
		Keys:    []string{zset1, zset2, zset3},
		Weights: []float64{float64(rand.Intn(optcluster.Loopstep))},
	}

	optcluster.ClusterClient.ZInterStore(zinterstore, &zstore)
	optcluster.ClusterClient.ZUnionStore(zunionstore, &zstore)

	optcluster.ClusterClient.Del(zset1, zset2, zset3)
	optcluster.ClusterClient.Expire(zinterstore, optcluster.EXPIRE)
	optcluster.ClusterClient.Expire(zunionstore, optcluster.EXPIRE)
	t2 := time.Now()
	zaplogger.Info("ExecCMD", zap.Int("db", optcluster.DB), zap.String("command", "BO_ZUNIONSTORE_ZINTERSTORE"), zap.String("KeySuffix", optcluster.KeySuffix), zap.Duration("time", t2.Sub(t1)))
}

//随机执行一个基础操作
func (optcluster *OptCluster) ExecRandOptCluster() {
	optcluster.OptType = BaseOptArray[rand.Intn(len(BaseOptArray))]
	optcluster.ExecOpt()
}

//随机执行一个基础操作半数概率为Set操作
func (optcluster *OptCluster) ExecRandOptHalfIsSetCluster() {
	if rand.Int()%2 == 0 {
		optcluster.OptType = BO_SET_SETNX
		optcluster.ExecOpt()
	} else {
		optcluster.ExecRandOptCluster()
	}
}

//遍历执行所有基本操作
func (optcluster *OptCluster) ExecAllBasicOptCluster() {
	for _, v := range BaseOptArray {
		optcluster.OptType = v
		optcluster.ExecOpt()
	}
}

//持续随机执行基础操作
func (optcluster *OptCluster) KeepExecBasicOptCluster(ctx context.Context, sleeptime time.Duration) {
	i := int64(0)
	keysuffix := optcluster.KeySuffix

	//会引起CROSSSLOT Keys in request don't hash to the same slot错误的命令列表
	tocluster_skip_array := map[OptType]string{
		BO_MSET_MSETNX:                        "BO_MSET_MSETNX",
		BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE: "BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE",
		BO_ZUNIONSTORE_ZINTERSTORE:            "BO_ZUNIONSTORE_ZINTERSTORE",
	}

	for {
		randi := rand.Intn(2 * len(BaseOptArray))
		optcluster.KeySuffix = keysuffix + strconv.FormatInt(i, 10)
		if randi < len(BaseOptArray) {
			optcluster.OptType = BaseOptArray[randi]
		} else {
			optcluster.OptType = BO_SET_SETNX
		}

		if _, ok := tocluster_skip_array[optcluster.OptType]; ok {
			continue
		}

		optcluster.ExecOpt()

		i++
		time.Sleep(sleeptime)
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}

}
