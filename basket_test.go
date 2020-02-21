package basket

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed.")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGroup_Get(t *testing.T) {
	// 统计key调用用户定义Getter次数
	loadcounts := make(map[string]int, len(db))

	gee := NewGroup("test", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			// 从db查找
			if v, ok := db[key]; ok {
				if _, ok := loadcounts[key]; !ok {
					loadcounts[key] = 0
				}
				loadcounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		// 第一次查找缓存，没有，从db中查找，并写入缓存
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of: " + k)
		}
		// 第二次查找缓存，已经写入缓存中
		if _, err := gee.Get(k); err != nil || loadcounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}

}
