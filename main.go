package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

// terminal symbols
type term struct {
	first []string
	last []string
}

// nonterminal symbols
type nonterm struct {
	first []string
	last []string
}

// read rules from selected file
func readFile(path string) []string {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Open file failed.")
		panic(err)
	}
	defer file.Close()

	var input []string
	i := 0

	buf := bufio.NewReader(file)
	for {
		a, _, c := buf.ReadLine()
		if c == io.EOF {
			break
		}
		input = append(input, string(a))
		i += 1
	}
	return input
}

func getNoTerm(input []string) map[string]nonterm {
	//the list of nonterm
	termlist := make(map[string]nonterm)

	for i:=0; i<len(input); i++ {
		termlist[string(input[i][0])] = nonterm{make([]string, 0), make([]string, 0)}
	}
	return termlist
}

func getTerm(input []string, nontermlist map[string]nonterm) map[string]term {
	termlist := make(map[string]term)

	for i:=0; i<len(input); i++ {
		nonSymbol := input[i][0]
		//flag is used to identify whether the terminal symbol should in the firstvt
		flag := false
		//lastTerm is used to save the last terminal symbal
		lastTerm := ""

		for j:=0; j<len(input[i]); j++ {
			//skip nonterminal symbol
			if unicode.IsUpper(rune(input[i][j]))  {
				continue
			}
			//skip space
			if input[i][j] == ' ' {
				continue
			}
			//skip arrow
			if input[i][j] == '-' && input[i][j+1] == '>' {
				j += 1
				continue
			}
			//skip |
			if input[i][j] == '|' {
				//when scan a | , should reset flag and lastTerm, also append lastTerm to nontermlist
				flag = false
				if lastTerm != "" {
					temp := append(nontermlist[string(nonSymbol)].last, lastTerm)
					nontermlist[string(nonSymbol)] = nonterm{nontermlist[string(nonSymbol)].first, temp}
				}
				lastTerm = ""
				continue
			}
			//append id and other terminal symbols to termlist
			if input[i][j] == 'i' && j < len(input[i])-1 && input[i][j+1] == 'd' {
				termlist["id"] = term{make([]string, 0), make([]string, 0)}
				//if flag is false, then this terminal symbol should append to firstvt
				if !flag {
					temp := append(nontermlist[string(nonSymbol)].first, "id")
					nontermlist[string(nonSymbol)] = nonterm{temp, nontermlist[string(nonSymbol)].last}
					flag = true
					lastTerm = "id"
				}
			} else {
				termlist[string(input[i][j])] = term{make([]string, 0), make([]string, 0)}
				if !flag {
					temp := append(nontermlist[string(nonSymbol)].first, string(input[i][j]))
					nontermlist[string(nonSymbol)] = nonterm{temp, nontermlist[string(nonSymbol)].last}
					flag = true
					lastTerm = string(input[i][j])
				}
			}
		}

		if lastTerm != "" {
			temp := append(nontermlist[string(nonSymbol)].last, lastTerm)
			nontermlist[string(nonSymbol)] = nonterm{nontermlist[string(nonSymbol)].first,temp}
		}
	}
	termlist["$"] = term{make([]string, 0), make([]string, 0)}
	return termlist
}

func (r nonterm) getFirst(symbol string) {

}
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error, usage main.exe filepath.")
	}
	path := os.Args[1]
	var input []string
	input = readFile(path)
	//for i:=0; i< len(input); i++ {
	//	fmt.Println(input[i])
	//}
	nontermlist := make(map[string]nonterm)
	termlist := make(map[string]term)
	nontermlist = getNoTerm(input)
	termlist = getTerm(input, nontermlist)
	//now, we get the list of terminal symbol & the list of nonterminal symbol
	for k, _ := range nontermlist {
		fmt.Println(k)
	}
	for k, _ := range termlist {
		fmt.Println(k)
	}
}


