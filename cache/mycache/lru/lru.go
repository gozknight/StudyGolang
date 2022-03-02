package lru

import "container/list"

// Cache 一个不是并发安全的cache.
type Cache struct {
	// 允许使用的最大内存
	maxBytes int64
	// 当前已使用的内存
	nowBytes int64
	// golang内置的双向链表
	linkedList *list.List
	// cache记录
	cache map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为 nil
	OnEvicted func(key string, value Value)
}

// 结点信息
type entry struct {
	key   string
	value Value
}

// Value 计算value使用的内存
type Value interface {
	Len() int
}

// Len 双向链表的长度，获取添加了多少条数据
func (c *Cache) Len() int {
	return c.linkedList.Len()
}

// New Cache构造器，创建一个Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		linkedList: list.New(),
		cache:      make(map[string]*list.Element),
		OnEvicted:  onEvicted,
	}
}

// Get 1、从字典中找到对应的双向链表的节点。2、将该节点移动到队尾
// 双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.linkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	// 队尾元素，即双向链表的最后一个结点
	ele := c.linkedList.Back()
	if ele != nil {
		c.linkedList.Remove(ele)
		kv := ele.Value.(*entry)
		// 删除cache中的映射
		delete(c.cache, kv.key)
		// 更新内存状态
		c.nowBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 增加新kv到缓存中或者修改v.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 修改
		c.linkedList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 新增
		ele := c.linkedList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超出最大内存，淘汰最近未使用的
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}
