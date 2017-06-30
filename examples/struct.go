package main

import (
	"fmt"

	"github.com/panux/encoding-sh"
)

func main() {
	dat, err := sh.Encode(struct {
		Number        int
		List          []string
		String        string
		ComplexString string
	}{
		Number:        1,
		List:          []string{"apples", "oranges", "pears"},
		String:        "Hello",
		ComplexString: "Hello world!\nGoodbye!\n",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dat))
}
