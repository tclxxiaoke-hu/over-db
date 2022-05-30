package use_redis

import (
	"time"
)

const (
	incrId         = "shop:incr_id"
	countBits      = 32
	beginTimestamp = 1652112000
)

//GetIncrId redis 自增id生成器
func GetIncrId() (int64, error) {
	// 当前时间与开始时间的秒数
	timestamp := time.Now().Unix() - beginTimestamp

	// 生成序号
	count, err := Incr(incrId)
	if err != nil {
		return 0, err
	}

	// 左移32位
	return timestamp<<countBits | count, nil
}

func Incr(key string) (v int64, err error) {
	return redisDb.Incr(key).Result()
}
