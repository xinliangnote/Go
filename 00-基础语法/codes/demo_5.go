package main

import (
	"fmt"
)

func main() {
	var arr =  [5] int {1, 2, 3, 4, 5}
	modifyArr(arr)
	fmt.Println(arr)
}

func modifyArr(a [5] int){
	a[1] = 20
}
