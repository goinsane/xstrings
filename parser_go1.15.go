// +build go1.15

package xstrings

import (
	"strconv"
)

func init() {
	parseComplex = strconv.ParseComplex
}
