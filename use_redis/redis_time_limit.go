package use_redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// 滑动时间窗口
func newScript() *redis.Script {
	return redis.NewScript(` 
		local flag
		-- key 
		local limitKey    = tostring(KEYS[1])
		-- 当前时间戳
		local mowTs       = tonumber(ARGV[1])
		--时间窗口以秒为周期
		local period      = tonumber(ARGV[2])

		--1.记录行为 value和score都是用时间戳
		local flag = redis.call("zadd",limitKey, mowTs, mowTs)
		if flag == false then
		  return 0
		end
		
		--2.移除时间窗口之前的行为记录
		local flag = redis.call("zremrangebyscore",limitKey, 0, mowTs - period)
		if flag == false then
		  return 0
		end
		
		--3.设置过期时间
		local flag = redis.call("expireat",limitKey,mowTs + period + 1)
		if flag == false then
		  return 0
		end
		
		--4.获取窗口内的行为数量
		local behaviorNum = redis.call("zcard",limitKey) 
		
		return behaviorNum
    `)
}

func ActionAllowed(userId, actionKey string, period int64) (err error) {
	key := fmt.Sprintf("hist:%s:%s", userId, actionKey)
	nowTs := time.Now().Unix()

	script := newScript()
	sha, err := script.Load(redisDb).Result()
	if err != nil {
		fmt.Println("script load err: ", err)
		return fmt.Errorf("err: %v", err)
	}
	ret := redisDb.EvalSha(sha, []string{key}, nowTs, period)
	if result, err := ret.Result(); err != nil {
		fmt.Println("ret load err: ", err)
		return fmt.Errorf("Execute Redis fail: %v\n", err.Error())
	} else if result.(int64) < 2 {
		fmt.Println("")
		fmt.Printf("userid: %s, result: %d", userId, result)
	}

	return
}
