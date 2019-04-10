package utils

import (
	"encoding/json"
	"fmt"
	"github.com/goforbroke1006/pbt-stress-testing/pkg/pbtHttp/bag"
	"strings"
)

const (
	LoginUri  = "user/rest/security/login"
	LoginBody = `{"username": "%s","password": "%s"}`
)

func LoadTokens(baseUrlAndContext, authDataArg string) []string {
	if len(authDataArg) == 0 {
		return nil
	}

	var tokens []string
	authUPP := strings.Split(authDataArg, ";")
	for _, cred := range authUPP {
		credArr := strings.Split(cred, ":")
		s, err := getToken(baseUrlAndContext, credArr[0], credArr[1])
		if nil != err {
			fmt.Printf("Error: %s\n", err.Error())
			continue
		}
		tokens = append(tokens, s)
	}

	return tokens
}

func getToken(baseUrlAndContext, username, password string) (string, error) {
	url := fmt.Sprintf("%s/%s", baseUrlAndContext, LoginUri)
	fmt.Println("URL:> ", url)

	bodyStr := fmt.Sprintf(LoginBody, username, password)
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	fmt.Println("BODY:> ", bodyStr)

	_, respBody, err := bag.Post(url, bodyStr, nil)
	if nil != err {
		return "", err
	}

	fmt.Println("RESPONSE:> ", string(respBody))
	var restObj interface{}
	err = json.Unmarshal(respBody, &restObj)
	if nil != err {
		return "", err
	}

	authToken := restObj.(map[string]interface{})["response"].(map[string]interface{})["token"].(string)
	fmt.Println("TOKEN:> ", authToken)

	return authToken, nil
}
