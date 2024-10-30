package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// checkPathes - проверяет пути на корректность,
// если папки dirPath не существует, то создает ее
func checkPathes(srcPath string, dirPath string) error {
	_, err := os.Stat(srcPath)
	if err != nil {
		return errors.New("Файл не существует")
	}
	if _, err = os.Stat(dirPath); err != nil {
		if err = os.Mkdir(dirPath, os.ModePerm); err != nil {
			return errors.New("Папка не может быть создана")
		}
		return errors.New("Папка не существует")
	}
	return nil
}

// getPageData - делает GET-запрос на сайт address,
// если запрос был успешен, сохраняет ответ в папку dirPath
func getPageData(address string, dirPath string) error {
	resp, err := http.Get("http://" + address)
	if err != nil {
		return errors.New("Адрес " + address + " не корректен")
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(dirPath, address+".html"))
	if err != nil {
		return errors.New("Файл для сайта " + address + "не может быть создан")
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.New("Файл для сайта " + address + "не можен быть заполнен")
	}
	return nil
}

func main() {
	srcFlag := flag.String("src", "", "src path ")
	dirFlag := flag.String("dir", "", "dir path")
	flag.Parse()
	startExec := time.Now()
	defer func() { fmt.Println(time.Since(start)) }()

	if *srcFlag == "" {
		fmt.Println("Используйте --src, чтобы указать путь до файла со ссылками")
		return
	}
	if *dirFlag == "" {
		fmt.Println("Используйте --dir, чтобы указать путь для папки с файлами,")
		fmt.Println("если такого пути не существует, то он будет создан автоматически")
		return
	}
	if err := checkPathes(*srcFlag, *dirFlag); err != nil {
		fmt.Println(err.Error())
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
			err := getPageData(address, *dirFlag)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Выполненно для " + address)
			}
			wg.Done()
		}(scanner.Text())
	}

	wg.Wait()
	fmt.Println("Завершено")
}
