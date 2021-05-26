package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hash     Hash           // hash函数
	hashMap  map[int]string // 虚拟节点与真实节点的映射表
}

// create Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // 默认算法
	}
	return m
}

/*
Add 函数允许传入 0 或 多个真实节点的名称。
对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，虚拟节点的名称是：strconv.Itoa(i) + key，
即通过添加编号的方式区分不同虚拟节点。
使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上。
在 hashMap 中增加虚拟节点和真实节点的映射关系。
最后一步，环上的哈希值排序。
*/
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hashValue := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hashValue)
			m.hashMap[hashValue] = key
		}
	}
	sort.Ints(m.keys)
}

/*
选择节点就非常简单了，第一步，计算 key 的哈希值。
第二步，顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。如果 idx == len(m.keys)，
说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
第三步，通过 hashMap 映射得到真实的节点。
*/
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hashValue := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hashValue
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
