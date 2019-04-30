package main

import (
	"bufio"
	"bytes"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"fmt"
	"os"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/slack", processHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func processHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	check(err)
	file, err := os.OpenFile("./data.txt", os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	content := []string{}
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	firstElement := content[0]
	// file.WriteString(firstElement)
	err = removeFileFirstLine(file)
	check(err)
	log.Fatal(res.Write([]byte(firstElement)))
}

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func removeFileFirstLine(file *os.File)  error {
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(make([]byte, 0, fileInfo.Size()))
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = io.Copy(buf, file)
	if err != nil {
		return err
	}
	line, err := buf.ReadString('\n')
	if err != nil && err != io.EOF {
		return err
	}
	fmt.Println(line)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	nw, err := io.Copy(file,buf)
	if err != nil {
		return err
	}
	err = file.Truncate(nw)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}
