package main

import (
	"github.com/logrusorgru/aurora"
	"github.com/quan-to/slog"
	"io"
	"net/http"
)

type externalip string

func (e externalip) Resolve(args interface{}) string {
	client := &http.Client{}
	slog := slog.Scope("m:External-IP")
	req, err := http.NewRequest("GET", "https://ifconfig.me/ip", nil)
	if err != nil {
		slog.Warn(err)
		panic(err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:61.0) Gecko/20100101 Firefox/61.0")

	res, _ := client.Do(req)
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	resp := string(bodyBytes)
	slog.Log(`	fetch (%s%s)`, aurora.Red(resp), aurora.Cyan(""))

	return resp
}

var Plugin externalip
