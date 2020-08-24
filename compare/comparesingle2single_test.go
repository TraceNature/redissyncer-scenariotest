package compare

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"math/rand"
	"reflect"
	"testcase/commons"
	"testing"
)

type TT struct {
	A   string
	B   int
	IsA bool
}

func (t *TT) Testtt() {
	var temp interface{}

	if t.IsA {
		temp = t.A
	} else {
		temp = t.B
	}

	fmt.Println(reflect.TypeOf(temp))

}
func TestCompare_CompareDB(t *testing.T) {

	saddr := "114.67.100.239:6379"
	opt := &redis.Options{
		Addr: saddr,
		DB:   0, // use default DB
	}
	opt.Password = "redistest0102"
	client := commons.GetGoRedisClient(opt)

	//compare := CompareSingle2Cluster{
	//	Source: client, Target: client, BatchSize: 10,
	//}

	for i := 0; i < 20; i++ {
		member := &redis.Z{Score: rand.Float64() * float64(rand.Int()), Member: i}
		client.ZAdd("z_aaa", member)
	}
	a := true
	b := false
	c := true
	d := false
	fmt.Println(client.ZRank("z_aaa", "10"))
	fmt.Println(client.ZRank("z_aaa", "100"))
	fmt.Println(client.ZRank("z_aa", "10"))

	fmt.Println(a && b)
	fmt.Println(a && c)
	fmt.Println(b && c)
	fmt.Println(b && d)
	//fmt.Println((a == b) || (b == d))

	fmt.Println(a ^ b)
	fmt.Println(a ^ c)
	fmt.Println(b ^ c)
	fmt.Println(b ^ d)

	//
	//compare.CompareDB()
	//tt := TT{
	//	A:   "abc",
	//	B:   1234,
	//	IsA: false,
	//}
	//
	//tt.Testtt()
	//rdb := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: []string{"114.67.67.7:16379",
	//		" 114.67.67.7:16380",
	//		" 114.67.83.163:16379 ",
	//		" 114.67.83.163:16380 ",
	//		" 114.67.112.67:16379 ",
	//		" 114.67.112.67:16380"},
	//	Password: "testredis0102",
	//})
	////rdb.ClientList()
	//
	//rdb.Set("123", "bbb", 0*time.Second)
	//rdb.Set("456", "bbb", 0*time.Second)
	//fmt.Println(rdb.Get("aaa"))
	//fmt.Println(rdb.Get("123"))
	//fmt.Println(rdb.Get("456"))
	//fmt.Println(rdb.ClusterKeySlot("aaa"))
	//fmt.Println(rdb.ClusterKeySlot("123"))
	//fmt.Println(rdb.ClusterKeySlot("456"))
	//
	//fmt.Println(rdb.ClusterNodes())

}
