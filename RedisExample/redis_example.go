package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/go-redis/redis"
)

type OnlineServers struct {
	Servers []string // `json:"option_A"`
}

func main() { //docker run -d -p 6379:6379 redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	// pong, err := client.Ping().Result()

	err := client.Set("Server1", "192.168.0.9", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("Server1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Server1", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}

	var testsvrs []string
	testsvrs = append(testsvrs, "localhost")
	osvrs := &OnlineServers{
		Servers: testsvrs,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(osvrs)
	client.Set("ServerList", b.String(), 0)
	val3, err := client.Get("ServerList").Result()

	dec := json.NewDecoder(strings.NewReader(val3))
	for {
		var m OnlineServers
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
	}

}
