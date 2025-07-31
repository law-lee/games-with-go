package pointers

import "fmt"

func Run() {
	x := 5
	fmt.Println(x)

	xPtr := &x

	fmt.Println(xPtr)
}
