package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"net/http"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

var file *os.File
var output *os.File
var scanner *bufio.Scanner

// wg is used to wait for goroutines to finish.
var wg sync.WaitGroup
var client http.Client

func init(){

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 2 {
		log.Fatal("Filenames must be provided as argument.")
	}

	file, err := os.Open(argsWithoutProg[0])
    if err != nil {
        log.Fatal(err)
	}

    scanner = bufio.NewScanner(file)
    if err := scanner.Err(); err != nil {
        log.Fatal(err)
	}

	setupClient()

}

func main() {

	argsWithoutProg := os.Args[1:]

	defer file.Close()

	output, err := os.OpenFile(argsWithoutProg[1], os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666 )
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		wg.Add(1)
		go get(scanner, output, client)
	}

	wg.Wait()
}


func get(scanner *bufio.Scanner, output *os.File, client http.Client){
	defer wg.Done()
	url := formatURL(scanner)

	// Issue the search
	resp, err := client.Get(url)

	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	defer resp.Body.Close()

	// Decode the JSON response into our struct type.
	var gr plpResponse
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err == nil {
		//write to file
		output.WriteString(fmt.Sprintln(url,":", gr.PagingOptions.Total))
	}
}

func setupClient(){
	timeout := time.Duration(5 * time.Minute)
	client = http.Client{
    	Timeout: timeout,
	}
}

func formatURL(scanner *bufio.Scanner) string {
	var b strings.Builder
	url:= scanner.Text()

	//arrange the url to search
	queryStrSign := "?"
	b.WriteString(url)
	if !strings.Contains(url, queryStrSign) {
		b.WriteString("?format=json")
	}else {
		b.WriteString("&format=json")
	}
	return b.String()
}

// Model json response
type plpResponse struct {
	PagingOptions struct {
		Total int `json:"Total"`
	} `json:"PagingOptions"`
}

