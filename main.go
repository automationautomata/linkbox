package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func CheckPathes(srcPath string, dirPath string) bool {
	_, err := os.Stat(srcPath)
	if err != nil {
		fmt.Println("File doesn't exist")
		return false
	}

	_, err = os.Stat(dirPath)
	if err != nil {
		fmt.Println("Directory doesn't exist")
		if err = os.Mkdir(dirPath, os.ModePerm); err != nil {
			fmt.Println("Directory can't be created")
			return false
		}
	}
	return true
}

func GetPageData(address string, dirPath string) {
	resp, err := http.Get("http://" + address)
	if err != nil {
		fmt.Println(address, err.Error())
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(dirPath, address+".html"))
	if err != nil {
		fmt.Println("File for", address, "can't be created")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("File for", address, "can't be filled", err.Error())
	}
}

func main() {
	srcFlag := flag.String("src", "", "src path ")
	dirFlag := flag.String("dir", "", "dir path")
	flag.Parse()

	if !CheckPathes(*srcFlag, *dirFlag) {
		return
	}

	file, err := os.Open(*srcFlag)
	if err != nil {
		fmt.Println("File", *srcFlag, "can't be read")
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go func(address string) {
			GetPageData(address, *dirFlag)
			wg.Done()
		}(scanner.Text())
	}

	wg.Wait()
	fmt.Println(*srcFlag, *dirFlag)
}
