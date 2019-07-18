//demo_22.go
package main

import "fmt"

func main() {

	for i := 1; i <= 10; i++ {
		if i == 6 {
			continue
		}
		fmt.Println("i =", i)
	}
}
