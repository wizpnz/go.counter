package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"flag"
)

func init() {
    enableLog := flag.Bool("log", false, "Enable or disable logging. It's useful to check workers count")
    flag.Parse()
    if !(*enableLog) {
	    log.SetOutput(ioutil.Discard)
    }
}

func main() {
	chanUrls := make(chan string)
	chanFinished := make(chan int)
	simultaneouslyUrlsCount := 5

	go controlExecuteFlow(chanUrls, chanFinished, simultaneouslyUrlsCount)
	chanFinished <- readData(chanUrls)
	close(chanUrls)
 	waitExecutionFinished(chanFinished)
}

func readData(chanUrls chan<- string) (urlsCount int) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		chanUrls <- s.Text()
		urlsCount++
	}
	return  urlsCount
}

func waitExecutionFinished(chanFinished <-chan int) {
	for {
		_, ok := <-chanFinished
		if !ok {
			break
		}
	}
}


func controlExecuteFlow(chanUrls <-chan string, chanFinished chan int, workers int) {

	chanResponse := make(chan http.Response)
	chanResult := make(chan int)
	chanErrors := make(chan bool)
	chanWorkers := make(chan bool, workers)
	keys := make(chan os.Signal, 1)

	signal.Notify(keys, os.Interrupt)
	var total int
	var errors int
	var processedUrls int
	inUrls := -1

	for {
		select {
		case url, ok := <-chanUrls:
			if ok {
				go func(url string) {
					chanWorkers <- true
					go getHTMLBody(url, chanResponse, chanErrors)
					go processResponse(chanResponse, chanResult, chanErrors)
				}(url)
			}
		case result := <-chanResult:
			go func() {
				<-chanWorkers
				processedUrls++
				total += result
			}()
		case <-chanErrors:
			<-chanWorkers
			errors++
		case count := <-chanFinished:
			inUrls = count
		case <-keys:
			break

		}

		if inUrls != -1 && inUrls == processedUrls+errors {
			break
		}
	}

	fmt.Printf("Total: %d\n", total)
	if errors > 0 {
		fmt.Printf("\n%d Error", errors)
		if errors > 1 {
			fmt.Print("s")
		}
		fmt.Println(" occured, while processing. Please, check log for additional information\n")
	}
	close(chanFinished)
}

func getHTMLBody(url string, chanRes chan<- http.Response, chanErrors chan<- bool) {
	log.Println("getHTMLBody Started")
	res, err := http.DefaultClient.Get(url)
	if err == nil {
		chanRes <- *res
	} else {
		log.Printf("Failed to get data from %s err: %s", url, err)
		chanErrors <- true
	}
	log.Println("getHTMLBody finished")
	return
}

func processResponse(chanRes <-chan http.Response, chanResult chan<- int, chanErrors chan<- bool) {
	log.Println("processResponse Started")
	res := <-chanRes
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to get Body from %s err: %s", res, err)
		chanErrors <- true
	} else {
		re := regexp.MustCompile("Go")
		//matches := re.FindAll(body, -1)
		matches := re.FindAllString(string(body), -1)
		fmt.Printf("Count for %s: %d\n", res.Request.URL, len(matches))
		chanResult <- len(matches)
	}
	log.Println("processResponse finished")
	return
}

