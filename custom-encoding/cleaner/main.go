package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ugorji/go/codec"
)

type Widget struct {
	Id     int
	Name   string
	Random MyType
}

type MyType map[string]interface{}

func main() {

	w := Widget{
		Id:     1,
		Name:   "Hi",
		Random: MyType{"yo1": "hi"},
	}

	mh := NewWidgetMsgpackHandler()
	wx := &WidgetX{msgpack: &mh}

	fmt.Println("\n**\n**(1) ORIGINAL: \n**")
	fmt.Println(w)

	b, err := wx.ToMsgPack(w)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n**\n**(2) MSGPACK ENCODED: \n**")
	fmt.Println(string(b))

	newWidget, err := wx.FromMsgPack(b)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n**\n**(3) DECODED VERSION: \n**")
	fmt.Println(newWidget)

}

// helpers to move between layers

func NewWidgetMsgpackHandler() codec.MsgpackHandle {

	var mh codec.MsgpackHandle

	customEncodeFunc := func(rv reflect.Value) ([]byte, error) {
		data := rv.Interface().(MyType)
		// insert appropriate implementation here
		b, err := json.Marshal(data)
		return b, err
	}

	customDecodeFunc := func(rv reflect.Value, bs []byte) error {
		var m MyType
		// insert appropriate implementation here
		err := json.Unmarshal(bs, &m)
		if err == nil {
			rv.Set(reflect.ValueOf(m))
		}
		return err
	}

	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
	mh.AddExt(reflect.TypeOf(MyType{}), 0, customEncodeFunc, customDecodeFunc)

	return mh
}

type WidgetX struct {
	msgpack *codec.MsgpackHandle
}

func (x *WidgetX) ToMsgPack(w Widget) ([]byte, error) {
	var b []byte
	enc := codec.NewEncoderBytes(&b, x.msgpack)
	err := enc.Encode(w)
	return b, err
}

func (x *WidgetX) FromMsgPack(b []byte) (Widget, error) {
	var w Widget
	dec := codec.NewDecoderBytes(b, x.msgpack)
	err := dec.Decode(&w)
	return w, err
}
