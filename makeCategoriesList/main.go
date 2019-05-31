package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	categExc := readList("../list/categExc.list")
	categAll := readList("../list/categAll.list")

	// Create list with only categories to use.
	b := bytes.Buffer{}
	exclude := false
	for _, l := range categAll {
		for _, lExc := range categExc {
			if strings.HasPrefix(l, lExc) {
				exclude = true
				continue
			}
		}
		if !exclude {
			b.WriteString(l + "\n")
		}
		exclude = false
		// fmt.Printf("%d-%s\n", i, l)
	}
	// Write to file.
	err := ioutil.WriteFile("../list/categUse", b.Bytes(), 0664)
	if err != nil {
		log.Fatal(err)
	}
}

// readlist uppercase, remove spaces and create a list of lines.
func readList(fileName string) []string {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Replace(string(b), " ", "", -1)
	s = strings.ToUpper(s)
	return strings.Split(s, "\n")
}

func printLines(lines []string) {
	for i, l := range lines {
		fmt.Printf("%d-%s\n", i, l)
	}
}
