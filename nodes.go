package basket

// NodePicker 节点选择interface
type NodePicker interface {
	// 根据key选择对应节点的NodeGetter
	PickNode(key string) (NodeGetter, bool)
}

// NodeGetter 不同的node对应的NodeGetter
type NodeGetter interface {
	// 从group和key对应的node地址查找缓存值
	Get(group string, key string) ([]byte, error)
}
