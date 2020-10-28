package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const MY_WEBSITE string = "https://cf-assgn.gramadur.workers.dev/links"
const ACCEPT_HEADER string = "Accept: text/html,application/json,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\n"
const CONNECTION_HEADER string = "Connection: close\r\n"
const HELP_STRING string = "\nThis profiler takes in a url and request count and makes those many requests to the url \n" +
	"options \n" +
	"--help Print this message" +
	"--url <string> The url you want the profiler to hit \n" +
	"--profile <int> the number of times you want to hit the url. This also additionally prints summary statistics of the execution\n" +
	""

/**
build the HTTP RequestBody and return it as byte array
@params: the url to send the request to
returns: the byte array []byte requestString
*/
func buildRequestByteArray(url2 url.URL) []byte {

	requestBody := fmt.Sprintf("GET %v HTTP/1.1\r\n", url2.String())
	requestBody += fmt.Sprintf("Host: %v\r\n", url2.Host)
	requestBody += fmt.Sprintf(ACCEPT_HEADER)
	requestBody += fmt.Sprintf(CONNECTION_HEADER)
	requestBody += fmt.Sprintf("\r\n")

	return []byte(requestBody)
}

//check if err is nil and print error message and stop the program
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

/**
Send a request to the specified URL and
@params : takes in the url object to which we need to send the request
1. Print the response to terminal,
2. returns: the status code int, responseLength int, requestDuration float64
*/
func sendRequest(url url.URL) (int, int, float64) {
	conn, err := net.Dial("tcp", url.Host+":80")
	//connection error - failed to establish socket connection
	checkError(err)

	//lets build the HTTP request body based on the url requested
	var requestObj = buildRequestByteArray(url)

	start := time.Now()
	_, err = conn.Write(requestObj)
	//could not write to conn fd
	checkError(err)

	result, err := ioutil.ReadAll(conn)
	elapsed := time.Since(start)

	//could not read from conn fd
	checkError(err)
	//all this is post processing so we wont include it in calculating metrics

	responseString := string(result)
	responseArray := strings.Split(responseString, " ")
	//split returned error
	checkError(err)
	statusCode, err := strconv.Atoi(responseArray[1])
	//atoi error
	checkError(err)
	//fmt.Println(responseString)
	//close the connection
	conn.Close()
	return statusCode, len(result), elapsed.Seconds()
}

/**
Given a list returns mean, median, smallest, largest values
*/
func calcMeanMedianSmallestLargest(list []float64) (float64, float64, float64, float64) {
	sort.Float64s(list)
	total := 0.0
	mean := 0.0
	median := 0.0

	for _, v := range list {
		total += v
	}
	mean = total / float64(len(list))
	middle := len(list) / 2
	if middle%2 == 1 {
		median = list[middle]
	} else {
		median = (list[middle-1] + list[middle]) / 2
	}
	return mean, median, list[len(list)-1], list[0]

}

func PrintMetrics(statusCodes []int, responseLengths []int, executionDurations []float64, errorCodes []int) {
	var requestCount int = len(statusCodes)
	fmt.Println("---------------Metrics---------------")
	fmt.Println("Number of Requests = ", requestCount)
	mean, median, slowest, fastest := calcMeanMedianSmallestLargest(executionDurations)
	fmt.Println("Fastest Time =", fastest*1000, "ms")
	fmt.Println("Slowest Time =", slowest*1000, "ms")
	fmt.Println("Mean Time =", mean*1000, "ms")
	fmt.Println("Median Time =", median*1000, "ms")
	numErrorCodes := len(errorCodes)
	fmt.Println("Percentage of Requests Succeeded =", (len(statusCodes)-numErrorCodes)/len(statusCodes)*100, "%")
	fmt.Println("Error Codes Returned = ", errorCodes)
	sort.Ints(responseLengths)
	fmt.Println("Size of the smallest response =", responseLengths[0], " bytes")
	fmt.Println("Size of the largest response =", responseLengths[len(responseLengths)-1], " bytes")
}

func main() {
	//fmt.Println(os.Args)
	if len(os.Args) == 1 {
		fmt.Println(HELP_STRING)
		return
	}

	helpArg := flag.String("help", "", "URL to make the request to")
	if *helpArg != "" {
		fmt.Println(HELP_STRING)
		return
	}
	urlArg := flag.String("url", MY_WEBSITE, "URL to make the request to")
	profilerArg := flag.Int("profile", 0, "URL to make the request to")
	flag.Parse()
	url, parseError := url.Parse(*urlArg)
	//if profiler is 0 run without metrics calculation
	if *profilerArg == 0 {
		_, _, _ = sendRequest(*url)
		return
	}

	if parseError != nil {
		log.Fatal("Error")
	}

	statusCodes := make([]int, *profilerArg)
	responseLengths := make([]int, *profilerArg)
	runTimes := make([]float64, *profilerArg)
	errorCount := make([]int, 0)

	for i := 0; i < *profilerArg; i++ {
		//fmt.Println("Sending request #", i+1)
		statusCode, respLength, duration := sendRequest(*url)
		//fmt.Println("duration=", duration)
		if statusCode >= 400 {
			errorCount = append(errorCount, statusCode)
		}
		statusCodes[i] = statusCode
		responseLengths[i] = respLength
		runTimes[i] = duration
	}
	PrintMetrics(statusCodes, responseLengths, runTimes, errorCount)

}
