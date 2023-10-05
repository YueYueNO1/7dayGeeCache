package geecache

type ByteView struct { //表示缓存值
	b []byte //将存储真实的缓存值，选择byte[]是为了能够支持任意的数据类型
}

// 返回被缓存对象所占的内存大小
func (v ByteView) Len() int {
	return len(v.b)
}

// 返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
