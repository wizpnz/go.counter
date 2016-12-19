package main

import (
	"testing"
	"net/http"
	"os"
	"log"
)

func TestMain(m *testing.M) {
	log.SetOutput(os.Stdout)
	ret := m.Run()
	os.Exit(ret)
}

func TestGetAndProcess(t *testing.T) {
	chanResponse := make(chan http.Response)
	chanResult := make(chan int)
	chanErrors := make(chan bool)
	go getHTMLBody("https://golang.org", chanResponse, chanErrors)
	go processResponse(chanResponse, chanResult, chanErrors)
	wantGoCount := 9
	if n := <-chanResult ; n != wantGoCount {
		t.Fatalf("Must be %d 'Go' in https://golang.org, got %d", wantGoCount, n)
	}
}

func TestControlExecutorFlow (t *testing.T) {
	chanUrls := make(chan string)
	chanFinished := make(chan int)
	simultaneouslyUrlsCount := 5
	go controlExecuteFlow(chanUrls, chanFinished, simultaneouslyUrlsCount)
	chanFinished <- GenerateUrls(chanUrls)
	close(chanUrls)
 	waitExecutionFinished(chanFinished)
}

func GenerateUrls (chanUrls chan<- string) int  {
	urls := []string{	"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",
				"https://golang.org",}
	for _,url := range(urls) {
		chanUrls <- url
	}
	return len(urls)
}