// http 用来节点间通信
package basket

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HttpGetter struct {
	baseURL string
}

// Get 实现NodeGetter的Get方法
// 从group和key对应的node地址查找缓存值
func (h *HttpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	// 从node地址查找缓存
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ NodeGetter = (*HttpGetter)(nil)

const (
	defaultBasePath = "/basket/" // 默认基础path
	defaultReplicas = 50         // 每个节点对应的虚拟节点值
)

type HttpPool struct {
	self     string
	basePath string

	mu          sync.Mutex             // 加锁
	nodes       *ConHash               // 缓存服务器节点
	HttpGetters map[string]*HttpGetter // HttpGetter
}

// NewHttpPool 新建一个HttpPool
func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Set 添加缓存服务器节点
func (hp *HttpPool) Set(nodes ...string) {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	// 添加节点
	hp.nodes = NewConHash(defaultReplicas, nil)
	hp.nodes.Add(nodes...)

	// 不同的peer对应不同的HttpGetter
	hp.HttpGetters = make(map[string]*HttpGetter, len(nodes))
	for _, node := range nodes {
		hp.HttpGetters[node] = &HttpGetter{baseURL: node + hp.basePath}
	}

}

// PickNode 找到key对应的缓存节点
func (hp *HttpPool) PickNode(key string) (NodeGetter, bool) {
	hp.mu.Lock()
	defer hp.mu.Unlock()

	// 从hash环上找到key对应的节点
	if node := hp.nodes.Get(key); node != "" && node != hp.self {
		return hp.HttpGetters[node], true
	}
	return nil, false
}

var _ NodePicker = (*HttpPool)(nil)

func (hp *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", hp.self, fmt.Sprintf(format, v...))
}

// ServerHTTP 接管http请求和响应
// 获取请求path对应的group和key：
// /<basepath>/<groupname>/<key>/
func (hp *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// path是否一致
	if !strings.HasPrefix(r.URL.Path, hp.basePath) {
		panic("Unexpected path : " + r.URL.Path)
	}

	hp.Log("%s %s", r.Method, r.URL.Path)

	// 拆分path获取group name 和 key
	parts := strings.SplitN(r.URL.Path[len(hp.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	groupName := parts[0]
	key := parts[1]

	// 获取对应的group
	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 从缓存中查找值
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(value.ByteSlice())

}
