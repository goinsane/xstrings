// +build ignore

package main

import (
	"github.com/goinsane/xstrings"
)

func main() {
	u := xstrings.NewUnmarshaler()
	m := xstrings.NewMarshaler()

	_ = u
	_ = m
}
