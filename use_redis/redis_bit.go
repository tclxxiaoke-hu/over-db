package use_redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

/*
@breif  位图使用举例
*/

//UserSign 签到
func UserSign(uid int) (err error) {
	var offset int = time.Now().Local().Day() - 1
	var keys string = buildSignKey(uid)
	err = redisDb.SetBit(keys, int64(offset), 1).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return
}

//GetSignCount 签到次数
//获取用户签到的次数
func GetSignCount(uid int) (int64, error) {
	var keys string = buildSignKey(uid)
	count := redis.BitCount{Start: 0, End: 31}
	return redisDb.BitCount(keys, &count).Result()
}

//GetSignInfo 获取当月签到情况
//根据需要自己实现返回
func GetSignInfo(uid int) (interface{}, error) {
	var keys string = buildSignKey(uid)
	var day int = time.Now().Local().Day()
	var dddd string = fmt.Sprintf("u%d", day)
	st, _ := redisDb.Do("BITFIELD", keys, "GET", dddd, 0).Result()
	f := st.([]interface{})
	var res []bool = make([]bool, 0)
	var days []string = make([]string, 0)
	var v int64 = f[0].(int64)
	fmt.Println(v)
	for i := day; i > 0; i-- {
		var pos int = (day - i) * -1
		var keys = time.Now().Local().AddDate(0, 0, pos).Format("2006-01-02")
		days = append(days, keys)
		var value = v>>1<<1 != v
		res = append(res, value)
		v >>= 1
	}
	fmt.Println(res)
	fmt.Println(days, len(days))
	return nil, nil
}

func buildSignKey(uid int) string {
	var nowDate = formatDate()
	return fmt.Sprintf("u:sign:%d:%s", uid, nowDate)
}

//获取当前的日期
func formatDate() string {
	return time.Now().Format("2006-01")
}
