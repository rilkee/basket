/*
LRU：最近最少使用缓存淘汰算法
优先移除最近最久未使用的数据。
过程：访问过的放到队尾，新加的放到队尾，超出容量删除队首数据。

*/
package basket

import "container/list"

// Lru is the cache
type Lru struct {
	capacity int                           // 容量
	node     *list.List                    // 双向链表
	hash     map[interface{}]*list.Element // 字典：存放缓存节点

}

type entry struct {
	key   string
	value interface{}
}

// NewLRU 新建Cache
func NewLRU(capacity int) *Lru {
	return &Lru{
		capacity: capacity,
		node:     list.New(),
		hash:     make(map[interface{}]*list.Element),
	}
}

// Get 获取缓存节点值
func (l *Lru) Get(key string) (value interface{}, ok bool) {
	// 查找hash中是否有该key，有的话返回
	// 同时移到末尾
	if el, ok := l.hash[key]; ok {
		l.node.MoveToBack(el)
		return el.Value.(*entry).value, true
	}

	return

}

// Put 添加缓存节点
func (l *Lru) Put(key string, value interface{}) {
	// 如果缓存中已经有
	// 移到末尾，更新缓存值
	if el, ok := l.hash[key]; ok {
		l.node.MoveToBack(el)
		el.Value.(*entry).value = value
		return
	}

	// 缓存中没有，添加到后面
	el := l.node.PushBack(&entry{key: key, value: value})
	l.hash[key] = el

	// 如果容量超出了
	if l.capacity != 0 && l.node.Len() > l.capacity {
		ele := l.node.Front() // 取队首元素
		if ele != nil {
			l.node.Remove(ele) // 从链表中删除
			kv := ele.Value.(*entry)
			delete(l.hash, kv.key) // 从字典中删除
		}

	}
}
