package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nyaxt/otaru/go-fuzz/filewritecache"
)

func main() {
	blob, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	out := filewritecache.Fuzz(blob)
	fmt.Printf("Fuzz out: %d\n", out)
}
