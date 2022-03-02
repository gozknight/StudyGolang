package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash 字节哈希映射到uint32
type Hash func(data []byte) uint32

// Map 一致性哈希散列键
type Map struct {
	// Hash 函数
	hash Hash
	// 虚拟节点个数
	replicas int
	// 哈希环 Sorted
	keys []int
	// 虚拟节点与真实节点的映射表
	hashMap map[int]string
}

// New 新建一个Map实例
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE //不实现自定义Hash，使用默认方法
	}
	return m
}

// Add 添加真实节点到Map中
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// Get 获取散列中最近的项目到提供的密钥
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}