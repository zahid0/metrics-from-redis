package main

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"sort"
	"strings"
)

var prefix string
var pool *redis.Pool

func metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	redisClient := pool.Get()
	defer redisClient.Close()
	keys, err := redis.Strings(redisClient.Do("KEYS", prefix+":*"))
	if err != nil {
		panic(err)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value, err := redisClient.Do("GET", key)
		if err != nil {
			panic(err)
		}
		metric := strings.TrimPrefix(key, prefix+":")
		v := strings.SplitN(metric, ":", 2)
		fmt.Fprintf(w, "%s{", v[0])
		for i, label := range strings.Split(v[1], ",") {
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
	var redisDb int
	flag.StringVar(&prefix, "p", "metrics", "Metrics key prefix")
	flag.StringVar(&redisAddress, "c", "redis:6379", "Redis connection string")
	flag.IntVar(&redisDb, "d", 0, "Redis db to connect")
	flag.Parse()
	pool = newPool(redisAddress, redisDb)
	http.HandleFunc("/metrics", metrics)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error Starting the HTTP Server : ", err)
	}
}

func newPool(redisAddress string, redisDb int) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddress, redis.DialDatabase(redisDb))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
