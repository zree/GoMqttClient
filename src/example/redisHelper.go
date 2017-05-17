package example

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)

func RedisInsert(key string,value string){
	c,err :=redis.Dial("tcp","localhost:6389")
	if err!=nil{
		fmt.Println("Connect to redis error",err)
		return
	}
	c.Do("APPEND",key,value)
	c.Close()
}