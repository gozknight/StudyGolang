package cache

// A ByteView 只读数据结构
type ByteView struct {
	b []byte
}

// Len 返回view的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回只读的数据
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 返回view的string
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
