package main

import (
	"fmt"
	"net/http"
)

// for redacting sensitive request headers before displaying.
func RedactSensitiveHeaders(headers http.Header) http.Header {
	// headers are excluded case-insensitively.
	headers.Del("authorization")
	headers.Del("proxy-authorization")
	headers.Del("cookie")
	headers.Del("set-cookie")
	headers.Del("x-api-key")
	headers.Del("x-auth-token")
	headers.Del("x-access-token")
	headers.Del("x-refresh-token")
	headers.Del("x-csrf-token")
	headers.Del("x-amz-security-token")
	headers.Del("x-goog-api-key")
	headers.Del("authentication-info")
	headers.Del("proxy-authenticate")

	return headers
}

// utils displaying LoggerConf.Results

func (b LoggerConf) DisplayReqStatus() {
	for v := range len(b.Results) {
		fmt.Println(b.Results[v].ReqResponse.Status)
	}
}

func (b LoggerConf) DisplayReqLatency() {
	for v := range len(b.Results) {
		fmt.Println("request latency:", b.Results[v].Reqtime)
	}
}

func (b LoggerConf) DisplayReqBodyLen() {
	for v := range len(b.Results) {
		fmt.Println("request body length:", b.Results[v].Reqlen)
	}
}

func (b LoggerConf) DisplayReqBody() {
	for v := range len(b.Results) {
		fmt.Println(b.Results[v].Reqbody)
	}
}

func (b LoggerConf) Err() {
	for _, v := range b.Results {
		if v.Reqerr != nil {
			fmt.Println(v.Reqerr)
		}
	}
}
