package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	requests := make(chan *request, requestsNumber)
	requests <- newRequests()
	var wg sync.WaitGroup
	for i := 0; i < requestsNumber; i++ {
		wg.Add(1)
		go sendRequest(requests, &wg)
		wg.Wait()
	}
	elapsedTime := time.Since(start)
	printResult(<-requests, elapsedTime)
}

func sendRequest(reqChannel chan *request, wg *sync.WaitGroup) {
	defer wg.Done()
	requestStart := time.Now()
	client := http.Client{Timeout: time.Duration(timeoutMilliseconds)}
	_, err := client.Get(address)
	tempRequests := <-reqChannel
	if err, ok := err.(net.Error); ok && err.Timeout() {
		incRejectedNumber(tempRequests)
		reqChannel <- tempRequests
	} else if err != nil {
		panic(err)
	} else {
		addRequestTime(tempRequests, requestStart)
		reqChannel <- tempRequests
	}
}

var (
	address             string
	requestsNumber      int
	timeoutMilliseconds float64
)

const (
	defaultAddress             = "https://www.google.com/"
	defaultRequestsNumber      = 10
	defaultTimeoutMilliseconds = 100
)

func init() {
	address = *flag.String("address", defaultAddress, "address")
	requestsNumberFlag := flag.String("requestsNumber", strconv.Itoa(defaultRequestsNumber), "requestsNumber")
	timeoutMillisecondsFlag := flag.String("timeoutMilliseconds", strconv.Itoa(defaultTimeoutMilliseconds), "timeoutMilliseconds")
	flag.Parse()
	var err error
	if address == "" {
		address = defaultAddress
		fmt.Printf("address default value is:%s , because address has incorrect value\n", address)
	}
	requestsNumber, err = strconv.Atoi(*requestsNumberFlag)
	if err != nil || requestsNumber <= 0 {
		requestsNumber = defaultRequestsNumber
		fmt.Printf("requestsNumber default value is:%d , because requestsNumber has incorrect value\n", requestsNumber)
	}
	timeoutMilliseconds, err = strconv.ParseFloat(*timeoutMillisecondsFlag, 64)
	timeoutMilliseconds *= 1000000
	if err != nil || requestsNumber <= 0 {
		timeoutMilliseconds = defaultTimeoutMilliseconds
		fmt.Printf("timeoutMilliseconds default value is: %f, because timeoutMilliseconds has incorrect value\n", timeoutMilliseconds)
	}
}

type request struct {
	requestTimes   []time.Duration
	numberRejected int64
}

func newRequests() *request {
	return &request{requestTimes: []time.Duration{}}
}

func maxTime(requestsTime []time.Duration) (maxTime time.Duration) {
	if len(requestsTime) <= 0 {
		return
	}
	maxTime = requestsTime[0]
	for i := 0; i < len(requestsTime); i++ {
		if requestsTime[i] > maxTime {
			maxTime = requestsTime[i]
		}
	}
	return
}
func minTime(requestsTime []time.Duration) (minTime time.Duration) {
	if len(requestsTime) <= 0 {
		return
	}
	minTime = requestsTime[0]
	for i := 0; i < len(requestsTime); i++ {
		if requestsTime[i] < minTime {
			minTime = requestsTime[i]
		}
	}
	return
}

func requestsAverageTime(requestTimes []time.Duration) time.Duration {
	if len(requestTimes) != 0 {
		return sum(requestTimes) / time.Duration(len(requestTimes))
	}
	return 0
}

func sum(times []time.Duration) (sum time.Duration) {
	for i := 0; i < len(times); i++ {
		sum += times[i]
	}
	return
}

func incRejectedNumber(requests *request) {
	requests.numberRejected++
}

func addRequestTime(requests *request, requestStart time.Time) {
	requests.requestTimes = append(requests.requestTimes, time.Since(requestStart))
}

func printResult(requests *request, elapsedTime time.Duration) {
	fmt.Println("End time of requests:", elapsedTime)
	fmt.Println("Average request time:", requestsAverageTime(requests.requestTimes))
	fmt.Println("Longest request time:", maxTime(requests.requestTimes))
	fmt.Println("Faster request time:", minTime(requests.requestTimes))
	fmt.Println("Responds number that didn't wait:", requests.numberRejected)
}
