/*

                            是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/

package basket

import (
	"fmt"
	"log"
	"sync"
)

// Getter 表示回调用户自定义数据源
type Getter interface {
	Get(key string) ([]byte, error)
}

// 利用回调函数实现源数据的获取
type GetterFunc func(key string) ([]byte, error)

// Get 实现了Getter接口
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 代表缓存命名组，用来对缓存逻辑上分类
type Group struct {
	name      string     // group name
	getter    Getter     // the getter callback
	mainCache cache      // 缓存
	nodes     NodePicker // 选择不同节点对应的Getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) // 全局group组
)

// NewGroup 新建一个group
func NewGroup(name string, capacity int, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{capacity: capacity},
	}

	groups[name] = g
	return g
}

// GetGroup 获取某个缓存组信息
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// RegisterNodes 注册缓存服务器节点
func (g *Group) RegisterNodes(nodes NodePicker) {
	if g.nodes != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.nodes = nodes
}

// Get 从缓存中获取key对应的值
func (g *Group) Get(key string) (ByteView, error) {
	// 如果key为空，error
	if key == "" {
		return ByteView{}, fmt.Errorf("Key is required")
	}

	// 先从缓存中找
	if v, ok := g.mainCache.Get(key); ok {
		return v, nil
	}

	// 缓存中没有的话，其他节点查找（用户本地或者远程节点）
	// 并加入缓存
	return g.Load(key)

}

// Load 从不同来源加载源数据
func (g *Group) Load(key string) (ByteView, error) {
	// 从远程节点
	if g.nodes != nil {
		// 查找key对应的缓存服务器节点
		if node, ok := g.nodes.PickNode(key); ok {
			var err error
			// 从该缓存服务器节点获取key的缓存值
			if value, err := g.getFromNode(node, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from node", err)
		}
	}

	return g.getLocally(key)
}

// getLocally 从用户实现的Getter本地数据源查找
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	// 加入缓存
	g.mainCache.Put(key, value)

	return value, nil

}

// getFromNode 从node节点获取缓存
func (g *Group) getFromNode(node NodeGetter, key string) (ByteView, error) {
	bytes, err := node.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
