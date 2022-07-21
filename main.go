package main

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"strings"
)

var prefix string
var pool *redis.Pool

func metrics(w http.ResponseWriter, r *http.Request) {
	redisClient := pool.Get()
	defer redisClient.Close()
	keys, err := redis.Strings(redisClient.Do("KEYS", prefix+":*"))
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		value, err := redisClient.Do("GET", key)
		if err != nil {
			panic(err)
		}
		v := strings.SplitN(key, ":", 3)
		fmt.Fprintf(w, "%s{", v[1])
		for i, label := range strings.Split(v[2], ",") {
			if i > 0 {
				fmt.Fprintf(w, ",")
			}
			label_kv := strings.SplitN(label, "=", 2)
			fmt.Fprintf(w, "%s=\"%s\"", label_kv[0], label_kv[1])
		}
		fmt.Fprintf(w, "} %s\n", value)
	}
}

func main() {
	var redisAddress string
	flag.StringVar(&prefix, "p", "metrics", "Metrics key prefix")
	flag.StringVar(&redisAddress, "c", "redis:6379", "Redis connection string")
	flag.Parse()
	pool = newPool(redisAddress)
	http.HandleFunc("/metrics", metrics)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
	}
}

func newPool(redisAddress string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddress)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
