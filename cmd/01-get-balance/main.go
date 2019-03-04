package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/goforbroke1006/pbt-stress-testing/pkg/pbtHttp/bag"
)

const (
	LoginUri  = "user/rest/security/login"
	LoginBody = `
{
  "username": "%s",
  "password": "%s"
}`
	BalanceUri = "user/rest/profile/balance"
)

var (
	baseUrlAndContext = flag.String("base-url", "", "Base URL of API")
	username          = flag.String("username", "", "Username for auth in API")
	password          = flag.String("password", "", "Password for auth in API")
	concurrency       = flag.Uint64("concurrency", 1, "Count of parallel streams of tasks")
	attempts          = flag.Uint64("attempts", 100, "Count of tasks in one stream")
	timeout           = flag.Uint64("timeout", 1000, "Timeout in milliseconds between requests")
)

func init() {
	flag.Parse()
}

func checkBalanceTask(requestsCount, timeout *uint64, token *string, reportCh chan string) {
	var i uint64
	for i = 0; i < *requestsCount; i++ {
		url := fmt.Sprintf("%s/%s", *baseUrlAndContext, BalanceUri)
		respBody, _ := bag.Get(url, token)
		var restObj interface{}
		json.Unmarshal(respBody, &restObj)
		amountStr := restObj.(map[string]interface{})["response"].(map[string]interface{})["amount"].(float64)

		reportCh <- fmt.Sprintf("%f", amountStr)
		time.Sleep(time.Duration(*timeout) * time.Millisecond)

	}
}

func main() {
	url := fmt.Sprintf("%s/%s", *baseUrlAndContext, LoginUri)
	fmt.Println("URL:> ", url)

	bodyStr := fmt.Sprintf(LoginBody, *username, *password)
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	fmt.Println("BODY:> ", bodyStr)

	respBody, _ := bag.Post(url, bodyStr, nil)

	var restObj interface{}
	json.Unmarshal(respBody, &restObj)

	authToken := restObj.(map[string]interface{})["response"].(map[string]interface{})["token"].(string)
	fmt.Println("TOKEN:> ", authToken)

	totalCount := (*concurrency) * (*attempts)
	reports := make(chan string, totalCount)

	startTime := time.Now()

	var i uint64
	for i = 0; i < *concurrency; i++ {
		go checkBalanceTask(attempts, timeout, &authToken, reports)
	}

	for i = 0; i < totalCount; i++ {
		fmt.Printf("Output [# %d]: %s\n", i, <-reports)
	}

	delta := time.Now().Sub(startTime)
	fmt.Println("Spend time: ", delta.Seconds())
}
