package main

import (
	"fmt"
	"strings"
	"strconv"
)

func main(){
	input := []byte("planetA CIUDAD VALOR\nplanetA CIUDAD2 OTROVALOR\n")
	lines := strings.Split(string(input), "\n")
	fmt.Printf("%q\n",lines)
	for i, line := range lines {
		if line != "" {
			splitLine := strings.Split(line, " ")
			fmt.Println(splitLine[1])
			fmt.Println(splitLine[2])
			lines[i] = splitLine[0] + " " + splitLine[1] + " " + "NuevoValor" +strconv.Itoa(i)
		}
	}
	fmt.Printf("%q\n",lines)
	output := strings.Join(lines, "\n")
	fmt.Printf("%q\n",output)
	fmt.Printf("%q\n",[]byte(output))
	//text := "UpdateName " + "in.Planet" +" "+ "in.City" + " " + "in.Value" +"\n"
	//fmt.Print(text)
}