package main

import (
	"fmt"
	"os"
	"github.com/tnsts/design-practice-2/examples"
)

func main() {
	var expression string
	for _, element := range os.Args[1:] {
		expression += element + " "
	}

	res, err := lab1.CalculatePostfix(expression)
	if err == nil {
		fmt.Println(res)
	} else {
		fmt.Println("Error!", err)
	}

}
