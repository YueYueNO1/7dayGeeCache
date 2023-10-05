package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 一致性哈希，是GeeCache从单节点走向分布式节点的重要环节
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            //虚拟节点倍数
	keys     []int          // Sorted  哈希环
	hashMap  map[int]string //虚拟节点的映射表，键是虚拟节点的哈希值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas, //虚拟节点倍数
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE //保证数据识别和标识
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //计算虚拟节点的哈希值
			m.keys = append(m.keys, hash)                      //添加到环上（切片类型）
			m.hashMap[hash] = key                              //增加虚拟节点和真实节点的映射关系
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash //实现思想是二分查找。idx、是为了找到第一个大于hash的数值
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
