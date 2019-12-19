package std

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestMd5Stream(t *testing.T) {
	selfFile := os.Args[0]
	stream, err := os.Open(selfFile)
	AssertError(err, "open file")
	defer CloseIgnoreErr(stream)
	md5Str1, err := Md5Stream(stream)
	AssertError(err, "Md5Stream")
	t.Log(md5Str1)
	_, err = stream.Seek(0, io.SeekStart)
	AssertError(err, "seek err")
	md5Str2, err := Md5Stream(stream)
	AssertError(err, "Md5Stream")
	t.Log(md5Str2)
	Assert(strings.Compare(md5Str1, md5Str2) == 0, "mismatch")
}
