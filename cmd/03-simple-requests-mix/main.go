package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/goforbroke1006/pbt-stress-testing/pkg/pbtHttp/bag"
	"github.com/goforbroke1006/pbt-stress-testing/pkg/utils"
)

const (
	balanceUri        = "user/rest/profile/balance"
	loadProfileUri    = "user/rest/profile/load"
	paymentMethodsUri = "user/rest/system/getFranchiserPaymentMethods"
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
	if nil == baseUrlAndContext {
	}

	tokens := utils.LoadTokens(*baseUrlAndContext, *auth)
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
		go doSomeRandRequest(i, attempts, timeout, &tokens[ti], reports)
	}

	for i = 0; i < totalCount; i++ {
		fmt.Printf("Output [# %d]: %s\n", i, <-reports)
	}

	delta := time.Now().Sub(startTime)
	fmt.Println("Spend time: ", delta.Seconds())
}

func doSomeRandRequest(index uint64, requestsCount, timeout *uint64, token *string, reportCh chan string) {
	fmt.Printf("starting %d-th task\n", index)

	var sc int
	//var respBody []byte
	var err error
	var reqName string

	var i uint64
	for i = 0; i < *requestsCount; i++ {

		requestType := i % 3

		go func() {
			switch requestType {
			case 0:
				reqName = "balance"
				sc, _, err = getBalanceTask(token)
			case 1:
				reqName = "profile"
				sc, _, err = loadProfileTask(token)
			case 2:
				reqName = "payment"
				sc, _, err = loadPaymentMethodsTask(token)
			}

			if nil != err {
				reportCh <- fmt.Sprintf("[c=%d %s] unexpected error : %s", index, reqName, err.Error())
				return
			}

			if 200 != sc {
				reportCh <- fmt.Sprintf("[c=%d %s] unexpected status code %d", index, reqName, sc)
				return
			}

			reportCh <- fmt.Sprintf("[c=%d %s] status code %d", index, reqName, sc)
		}()

		time.Sleep(time.Duration(*timeout) * time.Millisecond)

	}
}

func getBalanceTask(token *string) (int, []byte, error) {
	url := fmt.Sprintf("%s/%s", *baseUrlAndContext, balanceUri)
	return bag.Get(url, token)
}

func loadProfileTask(token *string) (int, []byte, error) {
	url := fmt.Sprintf("%s/%s", *baseUrlAndContext, loadProfileUri)
	return bag.Get(url, token)
}

func loadPaymentMethodsTask(token *string) (int, []byte, error) {
	url := fmt.Sprintf("%s/%s", *baseUrlAndContext, paymentMethodsUri)
	return bag.Get(url, token)
}
