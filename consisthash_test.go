package basket

import (
	"strconv"
	"testing"
)

func TestConHash(t *testing.T) {
	h := NewConHash(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	// 添加缓存服务器节点
	// 2, 12, 22, 4, 14, 24, 6, 16, 26
	h.Add("2", "4", "6")

	test := map[string]string{
		"2":  "2",
		"27": "2",
		"13": "4",
		"5":  "6",
	}

	for k, v := range test {
		if h.Get(k) != v {
			t.Fatalf("The key (%s) should belong to (%s), but (%s)", k, v, h.Get(k))
		}
	}

	// 添加新节点的变化8, 18, 28
	// 这时候test中的27会重新分配给节点8
	h.Add("8")
	test["27"] = "8"
	for k, v := range test {
		if h.Get(k) != v {
			t.Fatalf("The key (%s) should belong to (%s), but (%s)", k, v, h.Get(k))
		}
	}

}
