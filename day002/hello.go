package main

import "fmt"

func Hello(name string) string {
	if name == "" {
		name = "Golang"
	}
	return "Hello, " + name
}

func main() {
	fmt.Println(Hello("Yassine"))
}
