package std

import (
	"container/list"
)

type Buffer interface {
	ToArray() []byte
	Write(arr []byte)
	WriteN(arr []byte, len int)
	PeekN(offset, n int) []byte
	ReadN(n int) []byte
	ReadableLen() int
	PopN(n int)
	Clear()
}

type byteBuffer struct {
	data *list.List
}

func NewFIFOByteBuffer() Buffer {
	return &byteBuffer{
		data: list.New(),
	}
}

func (this *byteBuffer) ToArray() []byte {
	ret := make([]byte, 0, this.ReadableLen())
	for ele := this.data.Front(); ele != nil; ele = ele.Next() {
		ret = append(ret, ele.Value.(byte))
	}
	return ret
}

func (this *byteBuffer) Write(arr []byte) {
	for _, v := range arr {
		this.data.PushBack(v)
	}
}

func (this *byteBuffer) WriteN(arr []byte, len int) {
	for i := 0; i < len; i++ {
		this.data.PushBack(arr[i])
	}
}

//checkout ReadableLen before call this
func (this *byteBuffer) PeekN(offset, n int) []byte {
	arr := make([]byte, 0, n)
	var ele = this.data.Front()
	for i := 0; i < offset; ele = ele.Next() {
		if ele == nil {
			return arr
		}
		i++
	}
	for i, it := 0, ele; i < n && it != nil; it = it.Next() {
		arr = append(arr, it.Value.(byte))
		i++
	}
	return arr
}

//checkout ReadableLen before call this
func (this *byteBuffer) ReadN(n int) []byte {
	arr := make([]byte, 0, n)
	dataN := this.data.Len()
	for i := 0; i < n && i < dataN; i++ {
		front := this.data.Front()
		if front == nil {
			break
		}
		arr = append(arr, front.Value.(byte))
		this.data.Remove(front)
	}
	return arr
}

//array len
func (this *byteBuffer) ReadableLen() int {
	return this.data.Len()
}

func (this *byteBuffer) PopN(n int) {
	for i := 0; i < n; i++ {
		front := this.data.Front()
		if front != nil {
			this.data.Remove(this.data.Front())
		}
	}
}

func (this *byteBuffer) Clear() {
	this.PopN(this.ReadableLen())
}
