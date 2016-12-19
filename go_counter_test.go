/*
Test for go_counter main package
Copyright (C) Vadim.Karasev 2016

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"log"
	"net/http"
	"os"
	"testing"
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
	if n := <-chanResult; n != wantGoCount {
		t.Fatalf("Must be %d 'Go' in https://golang.org, got %d", wantGoCount, n)
	}
}

func TestControlExecutorFlow(t *testing.T) {
	chanUrls := make(chan string)
	chanFinished := make(chan int)
	simultaneouslyUrlsCount := 5
	go controlExecuteFlow(chanUrls, chanFinished, simultaneouslyUrlsCount)
	chanFinished <- GenerateUrls(chanUrls)
	close(chanUrls)
	waitExecutionFinished(chanFinished)
}

func GenerateUrls(chanUrls chan<- string) int {
	urls := []string{"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org",
		"https://golang.org"}
	for _, url := range urls {
		chanUrls <- url
	}
	return len(urls)
}
