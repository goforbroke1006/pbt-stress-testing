package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"runtime"
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
	auth              = flag.String("auth", "", "Usernames and password for auth in API, format 'login1:pass1;login2:pass2;login3:pass3'")
	concurrency       = flag.Uint64("concurrency", 1, "Count of parallel streams of tasks")
	attempts          = flag.Uint64("attempts", 100, "Count of tasks in one stream")
	timeout           = flag.Uint64("timeout", 1000, "Timeout in milliseconds between requests")
)

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(int(*concurrency))
}

func main() {

	var tokens []string

	authUPP := strings.Split(*auth, ";")
	for _, cred := range authUPP {
		credArr := strings.Split(cred, ":")
		s, err := getToken(credArr[0], credArr[1])
		if nil != err {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		tokens = append(tokens, s)
	}

	if len(tokens) == 0 {
		fmt.Println("API broken!!!")
		return
	}

	totalCount := (*concurrency) * (*attempts)
	reports := make(chan string, *concurrency)

	startTime := time.Now()

	var i uint64
	for i = 0; i < *concurrency; i++ {
		ti := i % uint64(len(tokens))
		go checkBalanceTask(i, attempts, timeout, &tokens[ti], reports)
	}

	for i = 0; i < totalCount; i++ {
		fmt.Printf("Output [# %d]: %s\n", i, <-reports)
	}

	delta := time.Now().Sub(startTime)
	fmt.Println("Spend time: ", delta.Seconds())
}

func checkBalanceTask(index uint64, requestsCount, timeout *uint64, token *string, reportCh chan string) {
	fmt.Printf("starting %d-th task\n", index)

	var i uint64
	for i = 0; i < *requestsCount; i++ {
		url := fmt.Sprintf("%s/%s", *baseUrlAndContext, BalanceUri)
		sc, respBody, err := bag.Get(url, token)

		if nil != err {
			reportCh <- fmt.Sprintf("[c=%d] unexpected error : %s", index, err.Error())
			continue
		}

		if 200 != sc {
			reportCh <- fmt.Sprintf("[c=%d] unexpected status code %d", index, sc)
			continue
		}

		var restObj interface{}
		json.Unmarshal(respBody, &restObj)
		amountStr := restObj.(map[string]interface{})["response"].(map[string]interface{})["amount"].(float64)

		reportCh <- fmt.Sprintf("[c=%d] %f", index, amountStr)
		time.Sleep(time.Duration(*timeout) * time.Millisecond)

	}
}

func getToken(username, password string) (string, error) {
	url := fmt.Sprintf("%s/%s", *baseUrlAndContext, LoginUri)
	fmt.Println("URL:> ", url)

	bodyStr := fmt.Sprintf(LoginBody, username, password)
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	fmt.Println("BODY:> ", bodyStr)

	_, respBody, err := bag.Post(url, bodyStr, nil)
	if nil != err {
		return "", err
	}

	var restObj interface{}
	err = json.Unmarshal(respBody, &restObj)
	if nil != err {
		return "", err
	}

	authToken := restObj.(map[string]interface{})["response"].(map[string]interface{})["token"].(string)
	fmt.Println("TOKEN:> ", authToken)

	return authToken, nil
}
