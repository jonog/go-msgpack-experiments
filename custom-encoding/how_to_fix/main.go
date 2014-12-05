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
	Random MyType // if this is interface{}, encoder is called. however decoder will not be called
}

type MyType map[string]interface{}

func main() {

	// Setup Msgpack

	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	// Add encoding and decoding functions for MyType
	// i.e. how is MyType converted to byte?
	// for illustration - use json to encode/decode

	customEncodeFunc := func(rv reflect.Value) ([]byte, error) {
		fmt.Println("\ncustomEncodeFunc (custom encoding for MyType):")
		data := rv.Interface().(MyType)
		fmt.Println("===> data: ")
		fmt.Println(data)
		b, err := json.Marshal(data)
		fmt.Println("===> string(b): ")
		fmt.Println(string(b))
		return b, err
	}

	customDecodeFunc := func(rv reflect.Value, bs []byte) error {
		fmt.Println("\ncustomDecodeFunc (custom decoding for MyType):")
		fmt.Println("===> string(bs): ")
		fmt.Println(string(bs))
		var m MyType
		err := json.Unmarshal(bs, &m)
		fmt.Println("===> m: ")
		fmt.Println(m)
		if err == nil {
			rv.Set(reflect.ValueOf(m))
		}
		return err
	}

	// Register the encoding/decoding functions of MyType
	// note: in latest versions of codec is recommended to use SetExt
	// AddExt is officially deprecated but currently implemented via SetExt

	mh.AddExt(reflect.TypeOf(MyType{}), 0, customEncodeFunc, customDecodeFunc)

	// Test Implementation

	var b []byte
	var msgpackHandler = &mh

	w := Widget{Name: "Hi", Id: 1, Random: MyType{"yo1": "hi"}}
	fmt.Println("\n**\n**(1) ORIGINAL: \n**")
	fmt.Println(w)

	enc := codec.NewEncoderBytes(&b, msgpackHandler)
	err := enc.Encode(w)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n**\n**(2) MSGPACK ENCODED: \n**")
	fmt.Println(string(b))

	var decodedWidget Widget
	dec := codec.NewDecoderBytes(b, msgpackHandler)
	err = dec.Decode(&decodedWidget)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n**\n**(3) DECODED VERSION: \n**")
	fmt.Println(decodedWidget)

}
