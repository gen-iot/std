package std

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const KJsonIndent = "    "
const KJsonIndentPrefix = ""

type JsonObject map[string]interface{}

type JsonObjectLessFun func(l JsonObject, r JsonObject) bool

func NewJsonObjectFromBytes(bytes []byte) (JsonObject, error) {
	out := NewJsonObject()
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func NewJsonObject() JsonObject {
	return make(JsonObject)
}

func (this JsonObject) Get(k string) interface{} {
	return this[k]
}

func (this JsonObject) GetOr(k string, def interface{}) interface{} {
	if v, ok := this[k]; ok {
		return v
	}
	return def
}

func (this JsonObject) HasKey(k string) bool {
	_, ok := this[k]
	return ok
}

func any2String(v interface{}, def string) string {
	if rv, ok := v.(string); ok {
		return rv
	}
	if rv, ok := v.(*string); ok {
		return *rv
	}
	return def
}

func (this JsonObject) GetStringOr(k string, def string) string {
	v, ok := this[k]
	if !ok {
		return def
	}
	return any2String(v, def)
}

func (this JsonObject) GetString(k string) string {
	return this.GetStringOr(k, "")
}

func any2Int(v interface{}, def int) int {
	value, err := strconv.Atoi(fmt.Sprintf("%v", v))
	if err != nil {
		return def
	}
	return value
}

func (this JsonObject) GetIntOr(k string, def int) int {
	v, ok := this[k]
	if !ok {
		return def
	}
	return any2Int(v, def)
}

func (this JsonObject) GetBool(k string) bool {
	v, ok := this[k]
	if !ok {
		return false
	}
	if vv, ok := v.(bool); ok {
		return vv
	}
	if vp, ok := v.(*bool); ok {
		return *vp
	}
	return false
}

func (this JsonObject) GetInt(k string) int {
	return this.GetIntOr(k, 0)
}

func (this JsonObject) GetJsonObject(k string) JsonObject {
	v, ok := this[k]
	if !ok {
		return nil
	}
	if rv, ok := v.(map[string]interface{}); ok {
		return rv
	}
	return nil
}

func (this JsonObject) GetJsonArray(k string) JsonArray {
	v, ok := this[k]
	if !ok {
		return nil
	}
	if rv, ok := v.([]interface{}); ok {
		return rv
	}
	return nil
}

func (this JsonObject) ToJsonString() (string, error) {
	return MarshalObject(this)
}

// 注意,在函数之间传递时,必须使用指针进行传递
type JsonArray []interface{}

const DefaultJsonArrayCap = 4

func NewJsonArray() JsonArray {
	return NewJsonArrayWithCap(DefaultJsonArrayCap)
}

func NewJsonArrayWithCap(cap int) JsonArray {
	ret := make(JsonArray, 0, cap)
	return ret
}

func (this *JsonArray) Add(o interface{}) *JsonArray {
	*this = append(*this, o)
	return this
}

func (this JsonArray) ToJsonString() (string, error) {
	return MarshalObject(this)
}

//将数据转换为json
func MarshalObject(e interface{}) (string, error) {
	if bs, err := json.Marshal(e); err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}
