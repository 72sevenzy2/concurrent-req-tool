package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// return config
type RequestResult struct {
	ReqResponse *http.Response
	Reqtime     time.Duration
	Reqbody     string
	Reqlen      int
	Reqerr      error
}

type LoggerConf struct {
	Results []*RequestResult
}

// for params
type LoggerParamConf struct {
	Reqclient    *http.Client
	Req          *http.Request
	AllowReqBody *bool
	ReqBodySize  *int
	ReqNo        *int
}

func Log(conf *LoggerParamConf) *LoggerConf { // move time.Now() (start) and end to inside gorountines.
	resultChan := make(chan *RequestResult, *conf.ReqNo)
	// ^ channels for number of workers per request
	for range *conf.ReqNo {
		go func(conf *LoggerParamConf, result chan<- *RequestResult) {
			start := time.Now()
			resp, resperr := conf.Reqclient.Do(conf.Req)
			if resperr != nil {
				result <- &RequestResult{
					Reqerr: resperr,
				}
				return
			}
			end := time.Since(start)

			result <- &RequestResult{
				ReqResponse: resp,
				Reqtime:     end,
				Reqbody:     "",
				Reqlen:      0,
				Reqerr:      nil,
			}

		}(conf, resultChan)
	}

	var requests = make([]*RequestResult, 0, *conf.ReqNo)

	for range *conf.ReqNo {
		requests = append(requests, <-resultChan)
	}

	fmt.Printf("visited to %s, with method %s", conf.Req.URL.Path, conf.Req.Method)
	fmt.Println(" full url:", conf.Req.URL)
	if conf.Req.URL.RawQuery != "" {
		fmt.Println("with query:", conf.Req.URL.RawQuery)
	}

	fmt.Println("request header details:")
	newHeaders := RedactSensitiveHeaders(conf.Req.Header.Clone())
	for header, value := range newHeaders {
		log.Printf("\nheader: %s | value: %s", header, value)
	}

	// request body printing
	for v := range len(requests) {
		var bodyprev string // body preview (string) var
		var max *int
		var bodybytes []byte
		var bodyerr error

		if *conf.AllowReqBody {
			max = conf.ReqBodySize
			if *max == 0 {
				max = IntPtr(1024) // default to 1 KiB logging limit.
			}
			if requests[v].Reqerr != nil {
				continue
			}

			bodyReader := io.LimitReader(requests[v].ReqResponse.Body, int64(*max))
			bodybytes, bodyerr = io.ReadAll(bodyReader)
			if bodyerr != nil {
				fmt.Println(bodyerr.Error())
			}

			requests[v].ReqResponse.Body = io.NopCloser(bytes.NewBuffer(bodybytes))
			fmt.Println("after truncating:", len(bodybytes))
			bodyprev = string(bodybytes)

		}

		final := len(bodybytes)

		requests[v].Reqlen = final
		requests[v].Reqbody = bodyprev
	}
	return &LoggerConf{
		Results: requests,
	}
}
