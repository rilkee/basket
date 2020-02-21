package main

import (
	"basket"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "230",
	"Mack": "444",
	"Fib":  "222",
}

// createGroup 创建缓存组
func createGroup() *basket.Group {
	return basket.NewGroup("test", 2<<10, basket.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// startCacheServer 开启缓存服务器
func startCacheServer(addr string, addrs []string, ba *basket.Group) {
	nodes := basket.NewHttpPool(addr)
	nodes.Set(addrs...)
	ba.RegisterNodes(nodes)
	log.Println("basket cache server is running: ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nodes))

}

func startAPIServer(addr string, ba *basket.Group) {
	http.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// 获取url中的key
			key := r.URL.Query().Get("key")
			// 从缓存中查找
			v, err := ba.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(v.ByteSlice())

		}))
	log.Println("API server is running: ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8888, "basket cache server port")
	flag.BoolVar(&api, "api", false, "start the api server")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	nodesAddr := map[int]string{
		8888: "http://localhost:8888",
		8887: "http://localhost:8887",
		8886: "http://localhost:8886",
	}

	addrs := make([]string, 3)

	for _, v := range nodesAddr {
		addrs = append(addrs, v)
	}

	b := createGroup()

	if api {
		go startAPIServer(apiAddr, b)
	}

	startCacheServer(nodesAddr[port], []string(addrs), b)

}
