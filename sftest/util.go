package main

import "fmt"

func logReq(buf []byte) {
	fmt.Printf("-> %s\n", string(buf))
}

func logRes(buf []byte) {
	fmt.Printf("<- %s\n", string(buf))
}
