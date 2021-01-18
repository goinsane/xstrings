// +build ignore

package main

import (
	"fmt"
	"reflect"

	"github.com/goinsane/xstrings"
)

func main() {
	p := xstrings.Unmarshaler{}

	str := "10"

	var x int
	err := p.UnmarshalByValue(str, reflect.ValueOf(&x))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(x)

	var y int
	err = p.Unmarshal(str, &y)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(y)

	z := 10
	ifc, err := p.Parse(str, reflect.TypeOf(&z))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*ifc.(*int))

	var s struct {
		X int
		Y int
	}
	err = p.Unmarshal(`{"X": 5, "Y": 6}`, &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)

	var a *int
	//a = new(int)
	err = p.Unmarshal(`10`, a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
}
