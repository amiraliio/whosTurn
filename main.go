package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/slack", processHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func processHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	check(err)
	file, err := os.OpenFile("./data.txt", os.O_RDWR, 0755)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	content := []string{}
	for scanner.Scan() {
		content = append(content, scanner.Text())
	}
	switch req.FormValue("command") {
	case "/whosturn------>[randomly]":
		randomTurn(content, file, res)
	case "/whosturn------>[sequentially]":
		sequenceTurn(content, file, res)
	default:
		log.Println("commands are wrong")
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func randomTurn(content []string, file *os.File, res http.ResponseWriter) {
	rand.Seed(time.Now().Unix())
	name := rand.Int() % len(content)
	log.Fatal(res.Write([]byte(content[name])))
}

func sequenceTurn(content []string, file *os.File, res http.ResponseWriter) {
	firstElement := content[0]
	err := removeFileFirstLine(file)
	check(err)
	_, err = file.Seek(0, io.SeekEnd)
	check(err)
	writer := bufio.NewWriter(file)
	writer.WriteString("\n" + firstElement)
	writer.Flush()
	log.Fatal(res.Write([]byte(firstElement)))
}

func removeFileFirstLine(file *os.File) error {
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
	nw, err := io.Copy(file, buf)
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
