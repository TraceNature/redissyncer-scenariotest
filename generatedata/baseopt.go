package generatedata

type OptType int32

const (
	BO_APPEND OptType = iota
	BO_BITOP
	BO_DECR_DECRBY
	BO_INCR_INCRBY_INCRBYFLOAT
	BO_MSET_MSETNX
	//BO_PSETEX_SETEX
	BO_PFADD
	BO_PFMERGE
	BO_SET_SETNX
	BO_SETBIT
	BO_SETRANGE
	BO_HINCRBY_HINCRBYFLOAT
	BO_HSET_HMSET_HSETNX
	BO_LPUSH_LPOP_LPUSHX
	BO_LREM_LTRIM_LINSERT
	BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH
	BO_BLPOP_BRPOP_BRPOPLPUSH
	BO_SADD_SMOVE_SPOP_SREM
	BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE
	BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM
	BO_ZPOPMAX_ZPOPMIN
	BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE
	BO_ZUNIONSTORE_ZINTERSTORE
)

var BaseOptArray = []OptType{
	BO_DECR_DECRBY,
	BO_INCR_INCRBY_INCRBYFLOAT,
	BO_MSET_MSETNX,
	//BO_PSETEX_SETEX,
	BO_PFADD,
	BO_PFMERGE,
	BO_SET_SETNX,
	BO_SETBIT,
	BO_SETRANGE,
	BO_HINCRBY_HINCRBYFLOAT,
	BO_HSET_HMSET_HSETNX,
	BO_LPUSH_LPOP_LPUSHX,
	BO_LREM_LTRIM_LINSERT,
	BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH,
	BO_BLPOP_BRPOP_BRPOPLPUSH,
	BO_SADD_SMOVE_SPOP_SREM,
	BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE,
	BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM,
	BO_ZPOPMAX_ZPOPMIN,
	BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE,
	BO_ZUNIONSTORE_ZINTERSTORE,
}

func (ot OptType) String() string {
	switch ot {
	case BO_APPEND:
		return "BO_APPEND"
	case BO_BITOP:
		return "BO_BITOP"
	case BO_DECR_DECRBY:
		return "BO_DECR_DECRBY"
	case BO_INCR_INCRBY_INCRBYFLOAT:
		return "BO_INCR_INCRBY_INCRBYFLOAT"
	case BO_MSET_MSETNX:
		return "BO_MSET_MSETNX"
	//case BO_PSETEX_SETEX:
	//	return "BO_PSETEX_SETEX"
	case BO_PFADD:
		return "BO_PFADD"
	case BO_PFMERGE:
		return "BO_PFMERGE"
	case BO_SET_SETNX:
		return "BO_SET_SETNX"
	case BO_SETBIT:
		return "BO_SETBIT"
	case BO_SETRANGE:
		return "BO_SETRANGE"
	case BO_HINCRBY_HINCRBYFLOAT:
		return "BO_HINCRBY_HINCRBYFLOAT"
	case BO_HSET_HMSET_HSETNX:
		return "BO_HSET_HMSET_HSETNX"
	case BO_LPUSH_LPOP_LPUSHX:
		return "BO_LPUSH_LPOP_LPUSHX"
	case BO_LREM_LTRIM_LINSERT:
		return "BO_LREM_LTRIM_LINSERT"
	case BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH:
		return "BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH"
	case BO_BLPOP_BRPOP_BRPOPLPUSH:
		return "BO_BLPOP_BRPOP_BRPOPLPUSH"
	case BO_SADD_SMOVE_SPOP_SREM:
		return "BO_SADD_SMOVE_SPOP_SREM"
	case BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE:
		return "BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE"
	case BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM:
		return "BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM"
	case BO_ZPOPMAX_ZPOPMIN:
		return "BO_ZPOPMAX_ZPOPMIN"
	case BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE:
		return "BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE"
	case BO_ZUNIONSTORE_ZINTERSTORE:
		return "BO_ZUNIONSTORE_ZINTERSTORE"
	default:
		return ""
	}
}

type BaseOpt interface {
	BO_APPEND()
	BO_BITOP()
	BO_DECR_DECRBY()
	BO_INCR_INCRBY_INCRBYFLOAT()
	BO_MSET_MSETNX()
	BO_SET_SETNX()
	BO_SETBIT()
	BO_SETRANGE()
	BO_HINCRBY_HINCRBYFLOAT()
	BO_PFADD()
	//BO_PFMERGE()
	BO_HSET_HMSET_HSETNX()
	BO_LPUSH_LPOP_LPUSHX()
	BO_LREM_LTRIM_LINSERT()
	BO_RPUSH_RPUSHX_RPOP_RPOPLPUSH()
	BO_BLPOP_BRPOP_BRPOPLPUSH()
	BO_SADD_SMOVE_SPOP_SREM()
	BO_SDIFFSTORE_SINTERSTORE_SUNIONSTORE()
	BO_ZADD_ZINCRBY_ZPOPMAX_ZPOPMIN_ZREM()
	BO_ZPOPMAX_ZPOPMIN()
	BO_ZREMRANGEBYLEX_ZREMRANGEBYRANK_ZREMRANGEBYSCORE()
	BO_ZUNIONSTORE_ZINTERSTORE()
}
