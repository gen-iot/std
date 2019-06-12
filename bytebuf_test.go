package std

import (
	"testing"
)

func TestByteBuffer(t *testing.T) {
	buffer := NewByteBuffer()
	buffer.Write([]byte{0x01, 0x0e})
	Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	Assert(buffer.ReadInt16() == 0x010e, "bad")
	Assert(buffer.ReadableLen() == 0, "bad")

	buffer.WriteUInt16(0x010e)
	Assert(buffer.PeekInt8(1) == int8(0x0e), "bad")
	Assert(buffer.ReadInt16() == 0x010e, "bad")
	Assert(buffer.ReadableLen() == 0, "bad")

	buffer.WriteInt32(0x010e)
	Assert(buffer.PeekInt8(buffer.ReadableLen()-1) == int8(0x0e), "bad")
	Assert(buffer.ReadInt32() == 0x010e, "bad")
	Assert(buffer.ReadableLen() == 0, "bad")
}
