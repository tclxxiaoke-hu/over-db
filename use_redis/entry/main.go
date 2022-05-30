package main

import (
	"fmt"
	"over-db/use_redis"
)

func main() {
	use_redis.InitClient()

	v, err := use_redis.GetIncrId()
	if err != nil {
		fmt.Errorf("%v", err)
	}

	fmt.Println(v)
}
