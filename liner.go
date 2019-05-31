package main

import (
	"net"
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		_, err := net.Dial("tcp4", "192.168.11.141:8999")
		if err != nil {
			fmt.Println("err", err)
			break
		}
	}

	select {
	}
	return
}
