package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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

//append nonterm symbol to the tail of nontermlist
func getTerm(rules []string, nontermlist map[string]nonterm) map[string]term {
	termlist := make(map[string]term)

	for i:=0; i<len(rules); i++ {
		nonSymbol := rules[i][0]
		//flag is used to identify whether the terminal symbol should in the firstvt
		flag := false
		//lastTerm is used to save the last terminal symbal
		lastTerm := ""
		for j:=0; j<len(rules[i]); j++ {
			//skip nonterminal symbol
			if unicode.IsUpper(rune(rules[i][j]))  {
				continue
			}
			//skip space
			if rules[i][j] == ' ' {
				continue
			}
			//skip arrow
			if rules[i][j] == '-' && rules[i][j+1] == '>' {
				j += 1
				continue
			}
			//skip |
			if rules[i][j] == '|' {
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
			if rules[i][j] == 'i' && j < len(rules[i])-1 && rules[i][j+1] == 'd' {
				termlist["id"] = term{make([]string, 0), make([]string, 0)}
				//if flag is false, then this terminal symbol should append to firstvt
				if !flag {
					temp := append(nontermlist[string(nonSymbol)].first, "id")
					nontermlist[string(nonSymbol)] = nonterm{temp, nontermlist[string(nonSymbol)].last}
					flag = true
					lastTerm = "id"
				}
			} else {
				termlist[string(rules[i][j])] = term{make([]string, 0), make([]string, 0)}
				if !flag {
					temp := append(nontermlist[string(nonSymbol)].first, string(rules[i][j]))
					nontermlist[string(nonSymbol)] = nonterm{temp, nontermlist[string(nonSymbol)].last}
					flag = true
				}
				lastTerm = string(rules[i][j])
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

func getVT(rules []string, nontermlist map[string]nonterm) {
	//lastflag & firstflag is used to identify whether the firstvt & lastvt should append to the nonsymbol
	var firstflag bool
	var lastflag bool
	var flag bool
	var nonsymbol string
	var lastsymbol string
	for i:=len(rules)-1; i>=0; i-- {
		//fmt.Printf("%d \n", i)
		firstflag = false
		flag = false
		nonsymbol = string(rules[i][0])
		lastsymbol = ""
		lastflag = false
		for j:=0; j<len(rules[i]); j++ {
			//fmt.Printf("j: %d \n", j)
			if !flag && rules[i][j] == '-' && j < len(rules[i])-1 && rules[i][j+1] == '>' {
				j++
				flag = true
				continue
			}
			if flag {
				// if nonterminale symbol
				if unicode.IsUpper(rune(rules[i][j])){
					lastflag = true
					lastsymbol = string(rules[i][j])
					if firstflag {
						continue
					} else {
						firstflag = true
						temp := appendVT(nontermlist[nonsymbol].first, nontermlist[string(rules[i][j])].first)
						//fmt.Println("temp:")
						//fmt.Println(temp)
						nontermlist[nonsymbol] = nonterm{temp, nontermlist[nonsymbol].last}
					}
				} else if rules[i][j] == ' ' {
					//skip space
					continue
				} else if rules[i][j] == '|' {
					// if symbol is | , reset
					firstflag = false
					//append last to nonsymbol
					if lastflag && lastsymbol != "" {
						temp := appendVT(nontermlist[nonsymbol].last, nontermlist[lastsymbol].last)
						//fmt.Println("temp:")
						//fmt.Println(temp)
						nontermlist[nonsymbol] = nonterm{nontermlist[nonsymbol].first, temp}
					}
					lastflag = false
					lastsymbol = ""
				} else {
					//if read a terminal symbol, set the firstflag to true, the firstvt shoud not append to the symbol
					firstflag = true
					lastflag = false
				}
			}
		}
	}
	if lastsymbol != "" {
		temp := appendVT(nontermlist[nonsymbol].last, nontermlist[lastsymbol].last)
		nontermlist[nonsymbol] = nonterm{nontermlist[nonsymbol].first, temp}
	}
}

func appendVT(list1 []string, list2 []string) []string{
	sort.Sort(sort.StringSlice(list1))
	sort.Sort(sort.StringSlice(list2))
	var res []string
	for i:=0; i< len(list2); i++ {
		if _, ok := find(list1, list2[i]); !ok {
			res = append(res, list2[i])
		}
	}
	res = append(res, list1...)
	return res
}

func find(list1 []string, target string) (int, bool) {
	for i:=0; i<len(list1); i++ {
		if list1[i] == target {
			return i, true
		}
		if list1[i] > target {
			break
		}
	}
	return -1, false
}

//params: rules, nontermlist, termlist
func genTable(rules []string, nontermlist map[string]nonterm, termlist map[string]term, strToIdx map[string]int) (map[string][]byte, bool){
	size := len(termlist)

	//build a map , size is the length of termlist
	table := make(map[string][]byte)
	i := 0
	for k, _ := range termlist {
		i++
		table[k] = make([]byte, size)
	}

	//first, we should get the realtion ship between $ and start symbol
	table["$"][strToIdx["$"]] = '='
	nonsymbol := rules[0][0]

	for _, v := range nontermlist[string(nonsymbol)].first {
		table["$"][strToIdx[v]]= '<'
	}

	for _, v := range nontermlist[string(nonsymbol)].last {
		table[v][strToIdx["$"]] = '>'
	}

	var preSymbol []string
	var preNonSymbol []byte
	for i=0; i<len(rules); i++ {
		flag := false
		preNonSymbol = nil
		preSymbol = nil
		for j:=0; j<len(rules[i]); j++ {
			if !flag && rules[i][j] == '-' && j < len(rules[i])-1 && rules[i][j+1] == '>'  {
				flag = true
				j++
				continue
			}
			if flag {
				if rules[i][j] == ' ' {
					continue
				} else if rules[i][j] == '|' {
					preSymbol = nil
					preNonSymbol = nil
				} else if unicode.IsUpper(rune(rules[i][j])) {
					if ok := insertFirst(preSymbol, table, nontermlist[string(rules[i][j])].first, strToIdx, '<'); !ok {
						fmt.Println("文法具有二义性!")
						return table, false
					}
					//preSymbol = nil
					preNonSymbol = append(preNonSymbol, rules[i][j])
				} else {
					if preNonSymbol != nil {
						if ok := insertLast(rules[i][j], table, nontermlist[string(preNonSymbol[0])].last, strToIdx, '>'); !ok {
							fmt.Println("文法具有二义性!")
							return table, false
						}
					}

					if len(preSymbol) != 0 {
						for _, k := range preSymbol {
							table[k][strToIdx[string(rules[i][j])]] = '='
						}
					}
					preSymbol = append(preSymbol, string(rules[i][j]))
					preNonSymbol = nil
				}
			}
		}
	}
	return table, true
}

func insertFirst(preSymbol []string, table map[string][]byte, VT []string, strToIdx map[string]int, symbol byte ) bool {
	for _, s := range preSymbol {
		for _, i := range VT {
			//fmt.Println(string(table[s][strToIdx[i]]))
			if table[s][strToIdx[i]] != byte(0){
				return false
			}
			table[s][strToIdx[i]] = symbol
		}
	}
	return true
}

func insertLast(preSymbol byte, table map[string][]byte, VT []string, strToIdx map[string]int, symbol byte) bool {

	for _, i := range VT {
		if table[i][strToIdx[string(preSymbol)]] != byte(0) {
			return false
		}
		table[i][strToIdx[string(preSymbol)]] = symbol
	}
	return true
}

func output(strToIdx map[string]int, res map[string][]byte, size int) {

	idxToStr := make([]string, size)
	for k,v := range strToIdx {
		idxToStr[v] = k
	}

	fmt.Printf(" \t")
	for _, v := range idxToStr {
		fmt.Printf("%s \t", v)
	}

	for _, v := range idxToStr {
		fmt.Println()
		fmt.Printf("%s \t", v)
		for _, j := range res[v] {
			fmt.Printf("%s \t", string(j))
		}
	}
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
	/*
	for k, _ := range nontermlist {
		fmt.Println(k)
	}
	for k, _ := range termlist {
		fmt.Println(k)
	}

	 */
	getVT(input, nontermlist)

	for k, v := range nontermlist {
		fmt.Println("终结符为：")
		fmt.Printf("%s \n", k)
		fmt.Println("firstvt:")
		for _, i := range v.first {
			fmt.Printf("%s, ", i)
		}
		fmt.Println("\nlastvt:")
		for _, i := range v.last {
			fmt.Printf("%s , ", i)
		}
	}
	fmt.Println()


	strToIdx := make(map[string]int)
	i := 0
	for k, _ := range termlist {
		strToIdx[k] = i
		i++
	}

	var res map[string][]byte
	var ok bool
	if res, ok = genTable(input, nontermlist, termlist, strToIdx); !ok {
		return
	}

 	output(strToIdx, res, i)

}


