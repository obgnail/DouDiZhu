package main

import "fmt"

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
