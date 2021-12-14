package main

import (
	"fmt"
	"strings"
)

func main(){
	input := "planetA CIUDAD VALOR\nplanetA CIUDAD VALOR\n"
	lines := strings.Split(input, "\n")
	fmt.Printf("%q\n",lines)
	for _, line := range lines {
		splitLine := strings.Split(line, " ")
		fmt.Println(splitLine[1])
		fmt.Println(splitLine[2])
	}
}