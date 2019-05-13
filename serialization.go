package std

import (
	"bytes"
	"encoding/json"
	"github.com/vmihailenco/msgpack"
)

type Serialization interface {
	Marshal(v interface{}) ([]byte, error)
	UnMarshal(data []byte, v interface{}) error
}

type binderProxy struct {
	marshalFn   func(v interface{}) ([]byte, error)
	unmarshalFn func(data []byte, v interface{}) error
}

func (this *binderProxy) Marshal(v interface{}) ([]byte, error) {
	return this.marshalFn(v)
}

func (this *binderProxy) UnMarshal(data []byte, v interface{}) error {
	return this.unmarshalFn(data, v)
}

var JsonSerialization Serialization = &binderProxy{
	marshalFn:   json.Marshal,
	unmarshalFn: json.Unmarshal,
}

var MsgPackSerialization Serialization = &binderProxy{
	marshalFn:   MsgpackMarshal,
	unmarshalFn: MsgpackUnmarshal,
}

func MsgpackMarshal(o interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := msgpack.NewEncoder(buffer)
	encoder.UseJSONTag(true)
	err := encoder.Encode(o)
	return buffer.Bytes(), err
}

func MsgpackUnmarshal(data []byte, o interface{}) error {
	reader := bytes.NewReader(data)
	decoder := msgpack.NewDecoder(reader)
	decoder.UseJSONTag(true)
	err := decoder.Decode(o)
	return err
}
