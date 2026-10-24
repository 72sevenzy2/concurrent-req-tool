package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

// start a interactive session
func StartSession(b *bufio.Scanner, store *Data) {
	fmt.Println("session started.")
	for {
		fmt.Print(">")
		isOk := b.Scan() // take input
		if !isOk {
			fmt.Println(b.Err()) // possible scanner error case
		}

		input := b.Text()
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue // skip current iteration if no input
		}

		upperInput := strings.ToUpper(parts[0])

		switch upperInput {
		// declare variables (can be to store headers, urls, etc)
		case "VAR":
			ok := HandleVar(parts, store)
			if !ok {
				continue
			}
		case "GET":
			ok := HandleGet(parts, store)
			if !ok {
				continue
			}

		case "DEL":
			ok := HandleDel(parts, store)
			if !ok {
				continue
			}

			// actual api testing logic (GET only for now)
		case "TEST":
			RunTestCommand(parts, store)
		case "SPAM":
			conf := InitSpamConfig()

			ok := ExtractSpamUrl(parts, store, conf) // conf.reqUrl would hold url
			if !ok {
				continue
			}

			// needed defaults
			req, reqErr := http.NewRequest(http.MethodGet, conf.reqUrl, nil)
			conf.reqCap = 5 // N amount of requests that are to be sent

			if len(parts) <= 2 {
				fmt.Println("please include a request method (-x [method])")
				continue
			}
			uc := strings.ToUpper(parts[2])
			if uc == "-X" {

				conf.reqMethod = parts[3]
				remainingArgs := parts[4:]

				switch conf.reqMethod {
				case http.MethodGet, http.MethodHead, http.MethodDelete, http.MethodOptions, http.MethodTrace, http.MethodConnect:
					if ok := HandleMethodsWithoutBody(remainingArgs, conf, store); !ok {
						continue
					}

				case http.MethodPost, http.MethodPut, http.MethodPatch:
					if ok := HandleMethodsWithBody(remainingArgs, conf, store); !ok {
						continue
					}
				}

				// initialising final req
				req, reqErr = http.NewRequest(conf.reqMethod, conf.reqUrl, conf.reqBody)
				if reqErr != nil {
					fmt.Println(reqErr.Error())
					continue
				}

				res := make(chan Result) // 1 channel for all workers

				// NewWorker() param conf
				WorkerConf := &WorkerConfig{
					Req:     req,
					Client:  &conf.client,
					Results: res,
				}

				for i := range conf.reqCap {
					go func(id int) {
						NewWorker(WorkerConf)
					}(i)
				}

				for i := range conf.reqCap {
					r := <-res // receiving results

					if r.Error != nil {
						r.Err()
						continue
					}

					fmt.Printf("result number %d, successfuly retrieved.", i)
					fmt.Println("\n request body:")
					r.DisplayBody()

					fmt.Println("\n request headers:")
					r.DisplayReqHeaders()

					fmt.Println("\n request path:")
					r.DisplayReqPath()

					fmt.Println("\n request method:")
					r.DisplayReqMethod()

					fmt.Println("\n request status code:")
					r.DisplayReqStatusCode()
				}
			}

		case "FETCH":
			ok := HandleFetch(parts, store)
			if !ok { //
				continue
			}

		case "HELP":
			DisplayCmds()
			continue
		// for exiting
		case "EXIT":
			fmt.Println("exiting will remove saved variables.")
			return

		default:
			fmt.Println("not a valid command.")
			continue
		}
	}
}
