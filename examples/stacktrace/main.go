package main

import "fmt"

func main() {
	fmt.Println("stacktrace trimpath demo")
	runJob()
}

func runJob() {
	loadCustomerConfig()
}

func loadCustomerConfig() {
	panic("intentional panic: compare stack traces with and without -trimpath")
}
