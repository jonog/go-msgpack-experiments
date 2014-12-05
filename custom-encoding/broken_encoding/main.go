package main

import (
	"fmt"
	"reflect"

	"github.com/ugorji/go/codec"
)

type Widget struct {
	Name   string
	Id     int
	Random interface{}
}

func main() {

	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	var (
		b []byte
		h = &mh
	)

	w := Widget{
		Name: "Hi",
		Id:   1,
		Random: map[string]string{
			"yo1": "hi",
			"yo2": "there",
		},
	}
	fmt.Println("===> w: ")
	fmt.Println(w)

	enc := codec.NewEncoderBytes(&b, h)
	err := enc.Encode(w)
	if err != nil {
		panic(err)
	}

	fmt.Println("===> b: ")
	fmt.Println(string(b))

	var decodedWidget Widget
	dec := codec.NewDecoderBytes(b, h)
	err = dec.Decode(&decodedWidget)
	if err != nil {
		panic(err)
	}

	fmt.Println("===> decodedWidget: ")
	fmt.Println(decodedWidget)

}
