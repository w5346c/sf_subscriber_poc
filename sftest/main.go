package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	sender := NewSalesForceSyncEventSender()
	reader := bufio.NewReader(os.Stdin)

	for {
		reader.ReadString('\n')

		fmt.Println("Sending JoomProTestEvent...")
		res, err := sender.SendJoomProTestEventSync()
		if err != nil {
			fmt.Println("SendJoomProTestEventSync FAILED")
		} else {
			fmt.Println(res)
		}


	}
}
