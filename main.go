package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getPath(scanner *bufio.Scanner) string {
	var path string
	fmt.Print("The full path to the folder: ")
	scanner.Scan()
	path = scanner.Text()
	// fmt.Print(path + "\n")
	return path
}

func getSplitString(scanner *bufio.Scanner, sep string) []string {
	fmt.Print("Extensions for files (sep. by a space): ")
	scanner.Scan()
	line := scanner.Text()
	return strings.Split(line, sep)
}

func getRegexp(extensions []string) *regexp.Regexp {
	regexCond := ""
	if len(extensions) == 1 {
		regexCond = extensions[0]
	} else if len(extensions) > 1 {
		for i, item := range extensions {
			regexCond += item
			if i+1 < len(extensions) {
				regexCond += "|"
			}
		}
	}
	fmt.Println(regexCond)
	return regexp.MustCompile(regexCond)
}

func lineCounter(path string) (int, error) {
	r, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func getAllLines(path string, reg *regexp.Regexp) (int, error) {
	count := 0
	err := filepath.Walk(path,
		func(inPath string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if reg.MatchString(info.Name()) {
				var tmpCount int
				if info.IsDir() {
					fmt.Printf("Dir name: %s\n", info.Name())
					tmpCount, err = getAllLines(inPath+"\\"+info.Name(), reg)
				} else {
					fmt.Printf("File name: %s\n", info.Name())
					tmpCount, err = lineCounter(inPath)
				}
				if err != nil {
					fmt.Println(err)
					return err
				}
				count += tmpCount
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}
	return count, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	path := getPath(scanner)

	extensions := getSplitString(scanner, " ")

	reg := getRegexp(extensions)

	count, err := getAllLines(path, reg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nTotal result: %d", count)
}
