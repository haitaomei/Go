package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func main() { //docker run -d -p 6379:6379 redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	// pong, err := client.Ping().Result()

	err := client.Set("Key1", "value1", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("Key1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Key1", val)

	str := ""
	for i := 0; i < 2000; i++ {
		str += strconv.Itoa(i)
	}
	fmt.Println("value size:", binary.Size([]byte(str))/1024, "KB")
	start := time.Now()
	client.Set("key2", str, 0).Err()
	end := time.Now()
	fmt.Println("Write response time:", end.Sub(start))
	start = time.Now()
	val2, err := client.Get("key2").Result()
	end = time.Now()
	fmt.Println("Read response time:", end.Sub(start))
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		sizeInBytes := binary.Size([]byte(val2)) / 1024
		fmt.Println("key2's value's size:", sizeInBytes, "KB")
	}

	client.Do("SADD", "Key_for_set", "member1").Result()
	client.Do("SADD", "Key_for_set", "member2").Result()
	client.Do("SADD", "Key_for_set", "member3").Result()
	client.Do("SADD", "Key_for_set", "member1").Result()
	members := client.Do("SMEMBERS", "Key_for_set").Val()
	s := reflect.ValueOf(members)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}
	for i := 0; i < s.Len(); i++ {
		vs := s.Index(i).Interface().(string)
		fmt.Println(vs)
	}
	client.Do("SPOP", "Key_for_set").Result()
}
