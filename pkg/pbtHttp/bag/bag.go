package bag

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var client http.Client

func Get(url string, token *string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	if nil != token {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func Post(url string, body string, token *string) ([]byte, error) {
	body = strings.ReplaceAll(body, "\n", "")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	if nil != token {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)

}
