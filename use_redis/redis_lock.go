package use_redis

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/exp/rand"
)

// 分布式锁使用
var randNum int64

/*
Lock
@brief 分布式锁
key
delayTimeSeconds 过期时间 单位 秒
*/
func Lock(key string, delayTimeSeconds uint) (isLock bool, err error) {
	rand.Seed(uint64(time.Now().Unix()))
	atomic.AddInt64(&randNum, int64(rand.Intn(int(time.Now().Unix()))))
	err = redisDb.Do("set", key, randNum, "ex", delayTimeSeconds, "nx").Err()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		return
	}

	isLock = true
	return
}

func UnLock(key string) (err error) {
	tempV := atomic.LoadInt64(&randNum)
	script := lockScript()
	sha, err := script.Load(redisDb).Result()
	if err != nil {
		fmt.Println("script load err: ", err)
		return fmt.Errorf("err: %v", err)
	}
	ret := redisDb.EvalSha(sha, []string{key}, tempV)
	if result, err := ret.Result(); err != nil {
		fmt.Println("ret load err: ", err)
		return fmt.Errorf("Execute Redis fail: %v\n", err.Error())
	} else if result.(int64) == 1 {
		fmt.Println("成功了")
		fmt.Printf("userid: %s, result: %d", key, result)
	}

	return
}

func lockScript() *redis.Script {
	return redis.NewScript(` 
		local flag
		-- key 
		local limitKey    = tostring(KEYS[1])
		-- 输入的值
		local inNum       = tonumber(ARGV[1])

		--1.记录行为 value和score都是用时间戳
		local randNum = redis.call("get",limitKey)
		if inNum == randNum then
			local flag = redis.call("del",limitKey)
			if flag == false then
		  		return 0
			end
		end

		return 1
    `)
}
