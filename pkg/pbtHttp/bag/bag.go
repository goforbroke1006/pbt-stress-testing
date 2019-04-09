package bag

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func Get(url string, token *string) (int, []byte, error) {
	req, _ := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	if nil != token {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, respBodyData, err
}

func Post(url string, body string, token *string) (int, []byte, error) {
	body = strings.ReplaceAll(body, "\n", "")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	if nil != token {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBodyData, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, respBodyData, err

}
