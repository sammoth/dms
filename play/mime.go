// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/kostin88/dms/dlna/dms"
)

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		fmt.Println(dms.MimeTypeByPath(arg))
	}
}
