package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func byteToStr(bytes []byte) []string {
	var strBytes []string
	for _, v := range bytes {
		strBytes = append(strBytes, fmt.Sprintf("0x%x", v))
	}
	return strBytes
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("input error")
		return
	}

	input := os.Args[1]
	fmt.Println(input)
	chars, bytes := []rune(input), []byte(input)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Char", "Unicode", "UTF8"})
	curIdx := 0
	for _, char := range chars {
		step := 0
		if char <= 0x7F {
			step = 1
		} else if char <= 0x7FF {
			step = 2
		} else if char <= 0xFFFF {
			step = 3
		} else if char <= 0x01FFFF {
			step = 4
		} else {
			panic("overflow!")
		}
		table.Append([]string{
			fmt.Sprintf("%s", string(char)),
			fmt.Sprintf("\\u%x (%d)", char, char),
			fmt.Sprintf("%v", byteToStr(bytes[curIdx:curIdx+step])),
		})
		curIdx += step
	}
	table.Render()
}
