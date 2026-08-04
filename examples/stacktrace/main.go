package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("stacktrace trimpath demo")
	runJob()
}

func runJob() {
	printCallerFrames()
	loadCustomerConfig()
}

func loadCustomerConfig() {
	panic("intentional panic: compare stack traces with and without -trimpath")
}

func printCallerFrames() {
	fmt.Println()
	fmt.Println("[runtime.Caller]")
	for skip := 0; skip <= 3; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			fmt.Printf("skip=%d unavailable\n", skip)
			continue
		}
		fn := runtime.FuncForPC(pc)
		name := "<unknown>"
		if fn != nil {
			name = fn.Name()
		}
		fmt.Printf("skip=%d func=%s file=%s line=%d\n", skip, name, file, line)
	}
	fmt.Println()
}
