package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {

	fmt.Printf("\r%s", strings.Repeat(" ", 50))
	fmt.Printf("\rDownloading... %v MB complete", (float32(wc.Total/1024) / 1024))
}

func FileNameFinder(url string) string {
	file_name := ""
	for i := len(url) - 1; i >= 0; i-- {
		if string(url[i]) == "/" {
			//find the / and get the string after / for use of filename
			for j := i + 1; j < len(url); j++ {
				file_name = file_name + string(url[j])
			}
			return file_name
		}
	}
	return file_name
}

func Downloader(url string, file_name string) {
	file, err := os.Create(file_name)
	if err != nil {
		fmt.Print("Unable to create File!")
		os.Exit(0)
	}

	defer file.Close()

	result, err := http.Get(url)
	if err != nil {
		fmt.Print("[+] Connection Error!Unable to get the file..")
		os.Exit(0)
	}

	defer result.Body.Close()

	//Write to the file
	counter := &WriteCounter{}
	_, err = io.Copy(file, io.TeeReader(result.Body, counter))
	if err != nil {
		fmt.Print("Write Failed!")
		os.Exit(0)
	}

}

func main() {
	url_ := flag.String("u", "", "URL of the File")
	directory_ := flag.String("d", "", "Download Directory")
	flag.Parse() //Parsing the command line Flags

	if *url_ == "" && *directory_ == "" {
		fmt.Printf("\n Welcome! \n\n NOTE : Use Direct File Links \n\n USE -> ZDownloader <Options> \n \n OPTIONS  \n -u => Mention URL \n -d => Mention Directory \n\n EXAMPLE :ZDownloader.exe -u <url> -d <Directory> ")
		os.Exit(0)
	}

	// url->package url_-> URL of file
	file_name := FileNameFinder(*url_)
	final, _ := url.QueryUnescape(file_name)
	//Getting user home directory to apend
	user_default_directory, _ := os.UserHomeDir()
	final_directory := user_default_directory + `Downloads`

	if *directory_ != "" {
		final_directory = *directory_
	}

	//use ` `  to mention string without escaping
	file_path := final_directory + final

	fmt.Print("File Name :", final)
	fmt.Print("\nLocation:", file_path)
	fmt.Println()

	Downloader(*url_, file_path)

}
