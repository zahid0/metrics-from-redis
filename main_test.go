package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConvert(t *testing.T) {
	pool = newPool("redis:6379")
	redisClient := pool.Get()
	defer redisClient.Close()
	_, err := redisClient.Do("SET", "testmetrics:m1:l1=v1,l2=v2", 1)
	if err != nil {
		panic(err)
	}
	_, err = redisClient.Do("SET", "testmetrics:m2:l3=v3,l4=v4", 2)
	if err != nil {
		panic(err)
	}
	prefix = "testmetrics"
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()
	metrics(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "m2{l3=\"v3\",l4=\"v4\"} 2\nm1{l1=\"v1\",l2=\"v2\"} 1\n"
	if string(data) != expected {
		t.Fatalf("%s != %s", data, expected)
	}
}
