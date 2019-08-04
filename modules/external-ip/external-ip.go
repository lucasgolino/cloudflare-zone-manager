package main

import (
	"io/ioutil"
	"net/http"
)

type externalip string

func (e externalip) Resolve(args interface{}) (string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://ifconfig.me/ip", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")

	res, _ := client.Do(req)
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	return string(bodyBytes)
}

var Plugin externalip