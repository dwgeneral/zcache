// lru 缓存淘汰策略
package lru

import "container/list"

type Cache struct {
	maxBytes int64                    // 允许使用的最大内存
	nbytes   int64                    // 当前已使用的内存
	dl       *list.List               // Go 语言标准库实现的双向链表
	memo     map[string]*list.Element // 记录key对应的链表节点. key:string, value: 双向链表中对应节点的指针
}

// 双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// 值为实现了该接口的任意类型
type Value interface {
	Len() int
}

func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		dl:       list.New(),
		memo:     make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (val Value, ok bool) {
	if ele, ok := c.memo[key]; ok {
		c.dl.MoveToFront(ele)
		node := ele.Value.(*entry)
		return node.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.dl.Back() // 取队首节点, 从链表中删除
	if ele != nil {
		c.dl.Remove(ele)
		node := ele.Value.(*entry)
		delete(c.memo, node.key)
		c.nbytes -= int64(len(node.key)) + int64(node.value.Len())
	}
}

func (c *Cache) Add(key string, val Value) {
	if ele, ok := c.memo[key]; ok {
		c.dl.MoveToFront(ele)
		node := ele.Value.(*entry)
		c.nbytes += int64(val.Len()) - int64(node.value.Len())
		node.value = val
	} else {
		ele := c.dl.PushFront(&entry{key, val})
		c.memo[key] = ele
		c.nbytes += int64(len(key)) + int64(val.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.dl.Len()
}
