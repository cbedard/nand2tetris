package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var symbolTable = map[string]int{
	"R0": 0, "R1": 1, "R2": 2, "R3": 3, "R4": 4, "R5": 5, "R6": 6, "R7": 7, "R8": 8, "R9": 9, "R10": 10, "R11": 11,
	"R12": 12, "R13": 13, "R14": 14, "R15": 15, "SCREEN": 16384, "KBD": 24576, "SP": 0, "LCL": 1, "ARG": 2, "THIS": 3,
	"THAT": 4}
var nextSymbolAddress = 16

var destCodes = map[string]string{"": "000", "M": "001", "D": "010", "MD": "011", "A": "100", "AM": "101", "AD": "110",
	"AMD": "111"}
var compCodes0 = map[string]string{"0": "0101010", "1": "0111111", "-1": "0111010", "D": "0001100",
	"A": "0110000", "!D": "0001101", "!A": "0110001", "-D": "0001111", "-A": "0110011", "D+1": "0011111",
	"A+1": "0110111", "D-1": "0001110", "A-1": "0110010", "D+A": "0000010", "D-A": "0010011", "A-D": "0000111",
	"D&A": "0000000", "D|A": "0010101"}
var compCodes1 = map[string]string{"M": "1110000", "!M": "1110001", "-M": "1110011", "M+1": "1110111",
	"M-1": "1110010", "D+M": "1000010", "D-M": "1010011", "M-D": "1000111", "D&M": "1000000", "D|M": "1010101"}
var jumpCodes = map[string]string{"": "000", "JGT": "001", "JEQ": "010", "JGE": "011", "JLT": "100", "JNE": "101",
	"JLE": "110", "JMP": "111"}

func main() {
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

	//first pass to fill symbol table
	removedLookupLabels := 0
	currLine := 0
	for currLine < len(input) {
		if instructionType(input[currLine]) == "L" {
			symbol := parseL(input[currLine])
			if _, ok := symbolTable[symbol]; !ok {
				symbolTable[symbol] = currLine - removedLookupLabels
				removedLookupLabels++
			}
		}
		currLine++
	}

	//output file
	f, _ := os.Create(os.Args[1][:len(os.Args[1])-4] + ".hack")

	//create and write binary commands
	currLine = 0
	for ; currLine < len(input); currLine++ {
		outputString := "0000000000000000"
		if instructionType(input[currLine]) == "A" {
			symbol := parseA(input[currLine])
			//A Symbol -> Binary 0vvv... where vvv is 15-bit binary val of symbol
			num, parseError := strconv.Atoi(symbol)
			if parseError != nil { //is a label
				address, ok := symbolTable[symbol]
				if !ok { //not in symbol table
					symbolTable[symbol] = nextSymbolAddress
					num = nextSymbolAddress
					nextSymbolAddress++
				} else {
					num = address
				}
			}

			outputString += strconv.FormatInt(int64(num), 2)
			outputString = outputString[len(outputString)-16:] //trim output string to len 16
		} else if instructionType(input[currLine]) == "L" {
			continue
		} else { // Comp statement
			dest, comp, jump := parseC(input[currLine])
			//Binary 111accccccdddjjj, a is comp set 0 or 1, c is comp, d dest, j jump
			c := compCodes0[comp]
			if x, ok := compCodes1[comp]; ok {
				c = x
			}
			outputString = "111" + c + destCodes[dest] + jumpCodes[jump]
		}
		//write output string
		f.WriteString(outputString + "\n")
	}
	f.Close()
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

func instructionType(line string) string {
	if strings.Contains(line, "@") {
		return "A"
	} else if strings.Contains(line, "(") {
		return "L"
	} else {
		return "C"
	}
}

func parseA(instr string) string {
	dest := strings.Split(instr, "@")[1]
	return dest
}

func parseL(instr string) string {
	dest := strings.Trim(instr, "()")
	return dest
}

func parseC(instr string) (string, string, string) {
	line := strings.Split(instr, ";")
	destComp := strings.Split(line[0], "=")
	dest, comp, jump := "", "", ""
	if len(line) > 1 {
		jump = line[1]
	}

	if len(destComp) > 1 {
		dest = destComp[0]
		comp = destComp[1]
	} else {
		comp = destComp[0]
	}

	return dest, comp, jump
}
