package main

import (
	"fmt"
	"log"
	"net/http"
	"zcache"
)

var mockDb = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	zcache.NewGroup("scores", 2<<10, zcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := mockDb[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := zcache.NewHttpPool(addr)
	log.Println("zcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
