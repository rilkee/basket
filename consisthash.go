// 一致性hash算法实现
package basket

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// ConHash 一致性hash
type ConHash struct {
	hash     Hash           // hash算法
	replicas int            // 虚拟节点数
	nodes    []int          // hash环节点数
	hashMap  map[int]string // 虚拟节点-真实节点
}

// NewConHash
func NewConHash(replicas int, hash Hash) *ConHash {
	m := &ConHash{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add 添加节点到hash环上
func (m *ConHash) Add(nodes ...string) {
	for _, node := range nodes {
		// 将节点值根据指定的虚拟节点数利用hash算法放置到环中
		for i := 0; i < m.replicas; i++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + node)))
			m.nodes = append(m.nodes, h)
			// 映射虚拟节点到真实节点
			m.hashMap[h] = node
		}
	}
	sort.Ints(m.nodes)
}

// Get 从hash环上获取key对应的节点
func (m *ConHash) Get(key string) string {
	if len(m.nodes) == 0 {
		return ""
	}
	// 计算key的hash值
	h := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点
	idx := sort.Search(len(m.nodes), func(i int) bool {
		return m.nodes[i] >= h
	})

	// 从hash环查找
	// 返回hash映射的真实节点
	return m.hashMap[m.nodes[idx%len(m.nodes)]]

}
