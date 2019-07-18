//demo_23.go
package main

import "fmt"

func main() {
	fmt.Println("begin")

	for i := 1; i <= 10; i++ {
		if i == 6 {
			goto END
		}
		fmt.Println("i =", i)
	}

	END :
		fmt.Println("end")
}
