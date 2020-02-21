# BASKET

> "学习的最好方式是制造轮子。" —— 鲁迅

basket是一个分布式缓存服务，模仿go中的groupcache实现。

其他轮子：

- [nan](https://github.com/shiniao/nan/) —— Nan编程语言实现
- [mid](https://github.com/shiniao/mid/) —— markdown编译器
- [gaga](https://github.com/shiniao/gaga/) —— go语言web框架（模仿gin）
  



## 实现
- 利用lru算法实现缓存淘汰
- 支持并发存取
- 一致性hash算法实现节点分布式
- 缓存节点利用HTTP通信