package use_redis

import "fmt"

// ev统计 hyperLogLog

func LogAdd(key, v string) (err error) {
	err = redisDb.PFAdd(key, v).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func LogGet(key string) (res int64, err error) {
	res, err = redisDb.PFCount(key).Result()
	if err != nil {
		fmt.Println(err)
		return res, err
	}

	return res, nil
}
