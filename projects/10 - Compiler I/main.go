package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	//Parsing input .jack file
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide project directory!")
		return
	}

	dir, err := os.ReadDir(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	for i := range dir {
		if strings.Contains(dir[i].Name(), ".jack") {
			file, _ := os.Open(dir[i].Name())

			scanner := bufio.NewScanner(file)
			lines := make([]string, 0)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			input := tokenize(lines)

			//parse
			fmt.Println(input)

			file.Close()
		}
	}
}

func tokenType(token string) string {
	//INT_CONST, STRING_CONST, SYMBOL, KEYWORD
	if strings.Contains(token, "\"") {
		return "STRING"
	}

	_, err := strconv.Atoi(token)
	if err == nil {
		return "INT"
	}

	if len(token) == 1 {
		return "SYMBOL"
	} else if slices.Contains(keywords, token) {
		return "KEYWORD"
	} else {
		return "IDENTIFIER"
	}
}

var keywords = []string{
	"class",
	"constructor",
	"function",
	"method",
	"field",
	"static",
	"var",
	"int",
	"char",
	"boolean",
	"void",
	"true",
	"false",
	"null",
	"this",
	"let",
	"do",
	"if",
	"else",
	"while",
	"return",
}

func tokenize(input []string) []string {
	ignore := false
	res := make([]string, 0)
	for index := 0; index < len(input); index++ {
		if ignore && !strings.Contains(input[index], "*/") {
			continue //inside a comment
		}
		if input[index] == "" || input[index] == "\n" {
			//empty line, nothing
		} else if strings.Contains(input[index], "//") {
			//comment, nothing
		} else if strings.Contains(input[index], "/**") || strings.Contains(input[index], "/*") || strings.Contains(input[index], "*/") {
			ignore = !ignore
		} else {
			str := strings.Trim(input[index], " ")
			res = append(res, strings.Split(str, " ")...) //append all tokens of valid line
		}
	}
	return res
}
