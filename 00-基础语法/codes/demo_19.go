//demo_19.go
package main

import (
	"fmt"
)

func main() {

	person := [] string {"Tom", "Aaron", "John"}
	fmt.Printf("len=%d cap=%d slice=%v\n",len(person),cap(person),person)

	fmt.Println("")

	//循环
	for k, v := range person {
		fmt.Printf("person[%d]: %s\n", k, v)
	}

	fmt.Println("")

	for i := range person {
		fmt.Printf("person[%d]: %s\n", i, person[i])
	}

	fmt.Println("")

	for i := 0; i < len(person); i++ {
		fmt.Printf("person[%d]: %s\n", i, person[i])
	}

	fmt.Println("")

	//使用空白符
	for _, name := range person {
		fmt.Println("name :", name)
	}
}
