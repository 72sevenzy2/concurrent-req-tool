package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RunTestCommand(parts []string, store *Data, timeout *int) {
	target, ok := HandleTestValdiation(parts, store)
	if !ok {
		return
	}

	utilConf := InitTestUtilConf()
	reqConf := InitTestReqConf()

	method := http.MethodGet
	headers := make(http.Header)
	var body io.Reader
	hasBody := false

	for i := 2; i < len(parts); i++ {
		switch strings.ToUpper(parts[i]) {
		case "-H":
			headerSlice, ok := HandleHeaderAppending(parts, i, "-H")
			if !ok {
				return
			}

			if len(headerSlice) != 2 {
				fmt.Println("invalid header")
				return
			}

			headers.Add(headerSlice[0], headerSlice[1])
			i++

		case "-X":
			newMethod, ok := HandleRequestMethodValidation(parts, i)
			if !ok {
				return
			}

			if newMethod == "" {
				fmt.Println("invalid HTTP method")
				return
			}

			method = strings.ToUpper(newMethod)
			utilConf.reqType = method
			i++

		case "-D":
			if hasBody {
				fmt.Println("specify only one request body")
				return
			}

			cleanedData, ok := HandleJsonDataInput(parts, i-2)
			if !ok {
				// The helper expects the index of "-X".
				// Handle JSON directly if "-D" is parsed independently.
				if i+1 >= len(parts) {
					fmt.Println("missing JSON data after -D")
					return
				}

				cleanedData = strings.Join(parts[i+1:], " ")
				if !IsValidJson([]byte(cleanedData)) {
					fmt.Println("invalid JSON format")
					return
				}
			}

			body = strings.NewReader(cleanedData)
			hasBody = true

			if headers.Get("Content-Type") == "" {
				headers.Set("Content-Type", "application/json")
			}

			// JSON data consumes the remaining command text.
			i = len(parts)

		case "-F":
			if hasBody {
				fmt.Println("specify only one request body")
				return
			}

			if i+1 >= len(parts) {
				fmt.Println("missing form field; use -F key:value")
				return
			}

			field := strings.SplitN(parts[i+1], ":", 2)
			if len(field) != 2 || field[0] == "" {
				fmt.Println("invalid form field; use -F key:value")
				return
			}

			data := url.Values{}
			data.Set(field[0], field[1])
			body = strings.NewReader(data.Encode())
			hasBody = true
			i++

			if headers.Get("Content-Type") == "" {
				headers.Set(
					"Content-Type",
					"application/x-www-form-urlencoded",
				)
			}

		default:
			fmt.Printf("unknown TEST argument: %s\n", parts[i])
			return
		}
	}

	if utilConf.reqType == "" {
		utilConf.reqType = method
	}

	// todo: validate http method before constructing.
	reqConf.cl, reqConf.clErr = http.NewRequest(
		method,
		target,
		body,
	)
	if reqConf.clErr != nil {
		fmt.Println("request creation failed:", reqConf.clErr)
		return
	}

	reqConf.cl.Header = headers.Clone()

	client := &http.Client{
		Timeout: time.Duration(*timeout),
	}

	resp, err := client.Do(reqConf.cl)
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response status:", resp.Status)
	fmt.Println("response headers:")

	responseHeaders := resp.Header.Clone()
	responseHeaders.Del("Authorization")

	for header, values := range responseHeaders {
		for _, value := range values {
			fmt.Printf("%s: %s\n", header, value)
		}
	}
	fmt.Println("response body:")
	if method == http.MethodHead {
		fmt.Println("(HEAD responses do not include a body.)")
		return
	}
	const maxResponseBody = 1 << 20 // 1 MiB
	responseBody, err := io.ReadAll(
		io.LimitReader(resp.Body, maxResponseBody+1),
	)
	if err != nil {
		fmt.Println("reading response body failed:", err)
		return
	}

	if len(responseBody) > maxResponseBody {
		fmt.Printf(
			"%s\n[response truncated at %d bytes]\n",
			string(responseBody[:maxResponseBody]),
			maxResponseBody,
		)
		return
	}

	fmt.Println(string(responseBody))
}
