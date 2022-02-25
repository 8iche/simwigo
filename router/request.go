package router

import (
	"fmt"
	"io"
	"net/http"
)

type httpRequest struct {
	method  string
	path    string
	params  map[string][]string
	scheme  string
	host    string
	headers map[string]string
	body    string
}

func strHttpRequest(r httpRequest) string {
	colorBlue := "\033[34m"
	colorRed := "\033[31m"
	//colorGreen := "\033[32m"
	colorReset := "\033[0m"

	strReq := fmt.Sprintf("%s %s %s\n", r.method, r.path, r.scheme)

	strReq += fmt.Sprintf("Host: %s\n", r.host)

	for k, v := range r.headers {
		if k == "Cookie" {
			strReq += fmt.Sprintf("%s: %s%s%s\n", k, colorRed, v, colorReset)
			continue
		}
		strReq += fmt.Sprintf("%s: %s\n", k, v)
	}

	if len(r.body) > 0 {
		strReq += "\n" + colorBlue + r.body
	}

	strReq += colorReset

	return strReq
}

func formatHttpRequest(r *http.Request) httpRequest {
	req := httpRequest{
		method:  r.Method,
		path:    r.URL.String(),
		params:  r.URL.Query(),
		scheme:  r.Proto,
		host:    r.Host,
		headers: make(map[string]string),
	}

	for i, v := range r.Header {
		for _, w := range v {
			req.headers[i] += w
		}
	}

	body, _ := io.ReadAll(r.Body)
	r.Body.Close()
	req.body = string(body)

	return req

}
