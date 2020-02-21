package basket

import "testing"

func TestGet(t *testing.T) {
	lru := NewLRU(0)
	lru.Put("key1", "123")

	if ele, ok := lru.Get("key1"); !ok || ele.(string) != "123" {
		t.Fatalf("Get key1=123 failed.")
	}
}

func TestRemoveoldest(t *testing.T) {
	lru := NewLRU(2)
	lru.Put("key1", "123")
	lru.Put("key2", "456")
	lru.Put("key3", "789")
	// 容量超出限制，会删除队首的最近最少访问元素
	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("Lru shoudn't contain key1")
	}

}
