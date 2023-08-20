package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	S = "s"
	D = "d"
	B = "b"
	H = "h"
)

var HexChars = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func byteToHexStr(bytes []byte) []string {
	var strBytes []string
	for _, v := range bytes {
		strBytes = append(strBytes, fmt.Sprintf("0x%02x", v))
	}
	return strBytes
}

func byteToBinStr(bytes []byte) []string {
	var strBytes []string
	for _, v := range bytes {
		strBytes = append(strBytes, fmt.Sprintf("%08b", v))
	}
	return strBytes
}

func stringParser(input string) {
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
			fmt.Sprintf("%v", byteToHexStr(bytes[curIdx:curIdx+step])),
		})
		curIdx += step
	}
	table.Render()
}

func decimalNumberParser(input string) {
	number, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("input is not decimal number")
		return
	}
	fmt.Println(number)

	bytes := []byte{}
	for number > 0 {
		bytes = append(bytes, byte(number))
		number >>= 8
	}
	slices.Reverse(bytes)
	table := tablewriter.NewWriter(os.Stdout)
	table.Append(byteToHexStr(bytes))
	table.Append(byteToBinStr(bytes))
	table.Render()
}

func binParser(input string) {
	binChars := map[string]struct{}{"0": {}, "1": {}}

	number := 0
	chars := []byte(input)
	slices.Reverse(chars)
	for idx, val := range chars {
		val := string(val)
		if _, exist := binChars[val]; !exist {
			fmt.Println("input binary format error")
			return
		}
		if val == "1" {
			number += int(math.Pow(2, float64(idx)))
		}
	}

	decimalNumberParser(strconv.Itoa(number))
}

func hexParser(input string) {
	input = strings.ToLower(input)
	if len(input) < 3 || !strings.HasPrefix(input, "0x") {
		fmt.Println("input hex format error")
		return
	}

	var number int
	hexValue := input[2:]
	for _, char := range hexValue {
		idx := slices.Index(HexChars, string(char))
		if idx < 0 {
			fmt.Printf("%c is not hex char\n", char)
			return
		}
		number <<= 4
		number ^= idx
	}

	decimalNumberParser(strconv.Itoa(number))
}

func main() {
	parsers := map[string]func(string){
		S: stringParser,
		D: decimalNumberParser,
		B: binParser,
		H: hexParser,
	}

	var parseType string
	cmd := &cobra.Command{
		Use: "bd",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("input value is not found")
				return
			}
			input := args[0]
			parser := parsers[parseType]
			if parser == nil {
				fmt.Println("error type")
				return
			}
			parser(input)
		},
	}
	cmd.Flags().StringVarP(&parseType, "type", "t", D,
		"parse type\n s:parse utf8 encoded string\n d:parse decimal\n b:parse binary\n h:parse hex\n")

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
