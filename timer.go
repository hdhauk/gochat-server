package main

import (
	"fmt"
	"time"
)

func printRunTime() {
	start := time.Now()
	fmt.Println("GoChat-Server\n------------------")
	for {
		fmt.Printf("Server runtime: %.F seconds\r", time.Since(start).Seconds())
		time.Sleep(1 * time.Second)
	}
}
