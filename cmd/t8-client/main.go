package main

import (
	"github.com/Daniel-C-R/t8-client-go/internal/timeutil"
)

func main() {
	time, err := timeutil.IsoStringToTimestamp("2023-10-01T12:00:00")
	if err != nil {
		panic(err)
	}
	println(time) // Output: 1696166400
}
