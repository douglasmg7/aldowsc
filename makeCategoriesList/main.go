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
	// printLines(categExc)

	for _, l := range categAll {
		// No blank line;
		if strings.TrimSpace(l) == "" {
			continue
		}
		for _, lExc := range categExc {
			// No blank line;
			if strings.TrimSpace(lExc) == "" {
				continue
			}
			if strings.HasPrefix(l, lExc) {
				// fmt.Printf("Prefix : %s\n", lExc)
				// fmt.Printf("Exclude: %s\n\n", l)
				exclude = true
				break
			}
		}
		if !exclude {
			b.WriteString(l + "\n")
		}
		exclude = false
		// fmt.Printf("%d-%s\n", i, l)
	}
	// Remove last new line.
	sb := bytes.TrimRight(b.Bytes(), "\n")
	// Write to file.
	err := ioutil.WriteFile("../list/categUse.list", sb, 0664)
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
