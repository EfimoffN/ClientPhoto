package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"

	"os"

	"mime/multipart"
)

func main() {
	namesFiles()
}

func namesFiles() {

	var wg sync.WaitGroup

	url := "http://localhost:8080/ConvertPNG"

	path := "photos"

	files, err := ioutil.ReadDir("./" + path)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		wg.Add(1)
		go postSendPhoto(f.Name(), path, url, &wg)
	}

	wg.Wait()
}

func postSendPhoto(photoName string, path string, url string, wg *sync.WaitGroup) {

	defer wg.Done()

	pathFile, _ := os.Getwd()
	pathFile += `\` + path + `\` + photoName

	file, err := os.Open(pathFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	paramName := "file"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(pathFile))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(part, file)

	err = writer.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	bod, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bod))
}
