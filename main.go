package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	// "net/url"
	"os"

	"mime/multipart"
)

func main() {
	// postConvertJPG()
	postConvertPNG()
}

func postConvertJPG() {

	url := "http://localhost:8080/ConvertJPG"

	data := `{"User": "postConvertJPG", "Key": "red"}`

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func postConvertPNG() {

	url := "http://localhost:8080/ConvertPNG"

	file, err := os.Open("test.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	var sendBufBody bytes.Buffer

	multiPArWriter := multipart.NewWriter(&sendBufBody)

	fileWriter, err := multiPArWriter.CreateFormFile("file_field", "test.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fieldWriter, err := multiPArWriter.CreateFormField("normal_field")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		fmt.Println(err)
		return
	}

	multiPArWriter.Close()

	req, err := http.NewRequest("POST", url, &sendBufBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	ttt := multiPArWriter.FormDataContentType()
	fmt.Println(ttt)

	req.Header.Set("Content-Type", multiPArWriter.FormDataContentType())
	// req.Header.Set("Content-Type", "image/jpeg")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)

	log.Println(result)
}


