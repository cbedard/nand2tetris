package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var memAddress = map[string]string{
	"SP": "SP",
}

var output = make([]string, 0)

func main() {
	//Parsing input .asm file
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	input := removeCommentsAndWhitespace(lines)

	output = append(output, "@256")
	output = append(output, "D=A")
	output = append(output, "@SP")
	output = append(output, "M=D")

	for currLine := 0; currLine < len(input); currLine++ {
		commandType := getCommandType(input[currLine])

		if commandType == "C_ARITHMETIC" {
			//writeArithmetic(input[currLine])
		} else if commandType == "C_PUSH" {
			writePush(input[currLine])
		}
	}
}

func removeCommentsAndWhitespace(input []string) []string {
	res := make([]string, 0)
	for index := 0; index < len(input); index++ {
		if input[index] == "" || input[index] == "\n" {
			//empty line, nothing
		} else if strings.Contains(input[index], "//") {
			//comment, nothing
		} else {
			str := strings.Trim(input[index], " ")
			res = append(res, str)
		}
	}
	return res
}

func getCommandType(line string) string {
	args := strings.Split(line, " ")
	if len(args) < 2 {
		return "C_ARITHMETIC"
	} else if args[0] == "PUSH" {
		return "C_PUSH"
	} else {
		return "C_POP"
	}
}

func writePush(command string) {
	output = append(output, "@SP")
	output = append(output, "A=M")
	output = append(output, "M=D")
	output = append(output, "@SP")
	output = append(output, "M=M+1")
}

// SP Ops
func incrementSP() {
	aCommand("SP")
	cCommand("M", "M+1", "")
}

func decrementSP() {
	aCommand("SP")
	cCommand("M", "M-1", "")
}

func loadSP() {
	aCommand("SP")
	cCommand("A", "M", "")
}

// to stack ops
func valToStack(val string) {
	aCommand(val)
	cCommand("D", "A", "")
	compToStack("D")
}

func staticToStack(val string) {
	aCommand(val)
	cCommand("D", "M", "")
	compToStack("D")
}

func compToStack(val string) {
	loadSP()
	cCommand("M", val, "")
}

func aCommand(address string) {
	output = append(output, "@"+address)
}
func cCommand(dest, comp, jump string) {
	command := ""
	if dest != "" {
		command += dest + "="
	}
	command += comp
	if jump != "" {
		command += ";" + jump
	}
	output = append(output, jump)
}

var RegisterConsts = map[string]int{
	"R0": 0, "R1": 1, "R2": 2, "R3": 3, "R4": 4, "R5": 5, "R6": 6, "R7": 7, "R8": 8, "R9": 9, "R10": 10, "R11": 11,
	"R12": 12, "R13": 13, "R14": 14, "R15": 15, "SP": 0, "LCL": 1, "ARG": 2, "THIS": 3, "THAT": 4, "TEMP": 5,
	"FRAME": 13, "RET": 14, "COPY": 15}
