package std

import "reflect"

type CodecId int

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode(interface{}, []byte) error
	Id() CodecId
}

type TypeBinder struct {
	regType reflect.Type
}

func NewTypeBinder(src interface{}) *TypeBinder {
	binder := &TypeBinder{}
	binder.regType = reflect.TypeOf(src)

	return binder
}

// create new object
func (this *TypeBinder) NewObject() interface{} {
	return reflect.New(this.regType).Interface()
}

func (this *TypeBinder) Decode(codec Codec, data []byte) (interface{}, error) {
	obj := this.NewObject()
	if err := codec.Decode(obj, data); err != nil {
		return nil, err
	}
	return obj, nil
}
