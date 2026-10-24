package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// fetch cmd
func HandleFetch(parts []string, store *Data) bool {
	if len(parts) > 1 {
		fmt.Println("invalid fetch format: only run <FETCH> to retrieve all set variables.")
		return false
	}

	vals, ok := store.GetAll()
	if !ok {
		fmt.Println("no variables set.")
		return false
	}
	for v, i := range vals {
		fmt.Println(i, v)
	}

	return true
}

// var cmd
func HandleVar(parts []string, store *Data) bool {
	if len(parts) < 3 || len(parts) > 3 { // validate arguments
		fmt.Println("please consider the correct format: var <KeyName> <KeyValue>")
		return false
	}

	err := store.Set(parts[1], parts[2])
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("successful.")
	return true
}

// get cmd
func HandleGet(parts []string, store *Data) bool {
	if len(parts) < 2 || len(parts) > 2 {
		fmt.Println("please use the correct format: get <KeyName>")
		return false
	}

	val, err := store.Get(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println(val)
	return true
}

// delete cmd
func HandleDel(parts []string, store *Data) bool {
	if len(parts) < 2 || len(parts) > 2 { // validate if parts[1] (key) only is included in input
		fmt.Println("please use the correct format: del <KeyName>")
		return false
	}

	msg, ok := store.Del(parts[1])
	if !ok {
		fmt.Println(msg)
		return false
	}

	fmt.Println("deleted key.")
	return true
}

// modules for test cmd:

// config structs for test cmd dependent data.

// utility fields
type TestUtilConf struct {
	pass     bool
	reqType  string
	jsonData string
}

func InitTestUtilConf() *TestUtilConf {
	return &TestUtilConf{
		pass:     true, // default
		reqType:  "",
		jsonData: "",
	}
}

// request dependent fields
type TestReqConf struct {
	cl    *http.Request
	clErr error

	reqHeaders []string
	formBody   io.Reader
}

func InitTestReqConf() *TestReqConf {
	return &TestReqConf{
		cl:    nil,
		clErr: nil,

		reqHeaders: nil,
		formBody:   nil,
	}
}

func HandleTestValdiation(parts []string, store *Data) (string, bool) {
	if len(parts) < 2 {
		fmt.Println("variable as second argument does not exit, consider setting a var.")
		return "", false
	}

	// check if variable stored exists
	val, err := store.Get(parts[1])
	if err != nil {
		fmt.Println(err.Error())
		return "", false
	}

	return val, true
}

func HandleHeaderAppending(parts []string, partsIndex int, uppercasedH string) ([]string, bool) {
	// validate if values exist
	if partsIndex+1 >= len(parts) {
		fmt.Println("missing header values.")
		return nil, false
	}

	var returnSlice []string
	returnSlice = strings.SplitN(parts[partsIndex+1], ":", 2) // would create 2 substrings on : seperation

	if len(returnSlice) < 2 || len(returnSlice) > 2 {
		fmt.Println("please include both header name and value.")
		return nil, false
	}

	return returnSlice, true
}

func HandleRequestMethodValidation(parts []string, partsIndex int) (string, bool) {
	if partsIndex+1 >= len(parts) { // check if arguments after partsIndex + 1 exist
		fmt.Println("missing argument values, consider either POST or GET.")
		return "", false
	}

	return strings.ToUpper(parts[partsIndex+1]), true
}

func HandlePostValidation(parts []string, partsIndex int) (string, bool) {
	if partsIndex+2 >= len(parts) {
		fmt.Println("missing request body flag, consider: -d [data]")
		return "", false
	}

	return strings.ToUpper(parts[partsIndex+2]), true
}

// small util func for validating json inputs
func IsValidJson(str []byte) bool {
	if v := json.Valid(str); !v {
		return false
	}
	return true
}

func HandleJsonDataInput(parts []string, partsIndex int) (string, bool) {
	if partsIndex+3 >= len(parts) {
		fmt.Println("please include actual data in json format.")
		return "", false
	}

	cleaned := strings.Join(parts[partsIndex+3:], " ") // join json inputs

	valid := IsValidJson([]byte(cleaned))
	if !valid {
		fmt.Println("invalid json format.")
		return "", false
	}

	return cleaned, true
}

func HandleFormDataInput(parts []string, partsIndex int) ([]string, bool) {
	if partsIndex+3 >= len(parts) {
		fmt.Println("please include the necessary form data with format: 'title:value'")
		return nil, false
	}

	formParts := strings.SplitN(parts[partsIndex+3], ":", 2) // splits form data seperated by ":", into substrings

	// validate if formParts is of correct length before accessing slice (or will panic)
	if len(formParts) < 2 || len(formParts) > 2 {
		fmt.Println("please use the correct format.")
		return nil, false
	}

	return formParts, true
}

// spam cmd utils:

// worker response format
type Result struct {
	ReqHeaders    http.Header // http.Header is map[string][]string
	Body          string
	Error         error
	ReqPath       string
	ReqStatusCode int
	ReqMethod     string
}

// util funcs

func (b Result) Err() { // display err
	fmt.Println(b.Error)
}

func (b Result) DisplayBody() { // for body
	fmt.Println(b.Body)
}

func (b Result) DisplayReqHeaders() {
	for k, v := range b.ReqHeaders {
		fmt.Println(k, v)
	}
}

func (b Result) DisplayReqPath() {
	fmt.Println(b.ReqPath)
}

func (b Result) DisplayReqStatusCode() {
	fmt.Println(b.ReqStatusCode)
}

func (b Result) DisplayReqMethod() {
	fmt.Println(b.ReqMethod)
}

// config struct for NewWorker() parameters
type WorkerConfig struct {
	Req     *http.Request
	Client  *http.Client
	Results chan<- Result
	Id      int
}

func NewWorker(v *WorkerConfig) {
	resp, err := v.Client.Do(v.Req)
	var clonedHeaders http.Header

	if err != nil {
		v.Results <- Result{
			ReqHeaders:    clonedHeaders,
			Body:          "",
			Error:         err,
			ReqPath:       v.Req.URL.Path,
			ReqStatusCode: 0,
			ReqMethod:     v.Req.Method,
		}
		return
	}

	// shadowing sensitive header values.
	clonedHeaders = resp.Header.Clone()
	clonedHeaders.Del("authorization")
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil { // no body
		v.Results <- Result{
			ReqHeaders:    clonedHeaders,
			Body:          "",
			Error:         err,
			ReqPath:       resp.Request.URL.Path,
			ReqStatusCode: resp.StatusCode,
			ReqMethod:     resp.Request.Method,
		}
		return
	}
	v.Results <- Result{
		ReqHeaders:    clonedHeaders,
		Body:          string(body),
		Error:         nil,
		ReqPath:       resp.Request.URL.Path,
		ReqStatusCode: resp.StatusCode,
		ReqMethod:     resp.Request.Method,
	}
}

type SpamConf struct {
	client http.Client

	reqBody   io.Reader
	reqUrl    string
	reqMethod string
	reqCap    int
}

func InitSpamConfig() *SpamConf {
	return &SpamConf{
		reqBody:   nil,
		reqUrl:    "",
		reqMethod: "",
		reqCap:    0,
	}
}

// spam cmd modularization:

func ExtractSpamUrl(parts []string, store *Data, conf *SpamConf) bool {
	if len(parts) <= 1 {
		fmt.Println("please include a url.")
		return false
	}

	varVal, err := store.Get(parts[1])
	conf.reqUrl = varVal
	if err != nil {
		conf.reqUrl = parts[1] // if not real variable
	}

	return true
}

func HandleMethodsWithoutBody(remainingArgs []string, conf *SpamConf, store *Data) bool {

	if len(remainingArgs) < 2 {
		fmt.Println("not enough arguments, please consider: -C")
		return false
	}

	for i := 0; i < len(remainingArgs); i++ {
		switch strings.ToUpper(remainingArgs[i]) {
		case "-C":
			if len(remainingArgs) <= i+1 {
				fmt.Println("please include a value for -C.")
				continue
			}

			var temp string

			val, err := store.Get(remainingArgs[i+1])
			temp = val

			if err != nil { // if it isnt a variable
				temp = remainingArgs[i+1]
			}

			s := strings.Trim(temp, "\"")
			d, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("invalid number.")
				continue
			}

			conf.reqCap = d

			i++ // skip value of remainingArgs[i+1] so default case doesnt pick it up.

		default:
			fmt.Println("not a valid command.")
			continue
		}
	}
	return true
}

func HandleMethodsWithBody(remainingArgs []string, conf *SpamConf, store *Data) bool {
	if len(remainingArgs) < 4 {
		fmt.Println("not enough arguments, please consider: -c ['int'] and -d ['json'].")
		return false
	}

	for i := range len(remainingArgs) {
		switch strings.ToUpper(remainingArgs[i]) {
		case "-C":
			if len(remainingArgs) <= i+1 {
				fmt.Println("please include a value for -C.")
				return false
			}
			var temp string

			val, err := store.Get(remainingArgs[i+1])
			temp = val

			if err != nil {
				temp = remainingArgs[i+1]
			}

			s := strings.Trim(temp, "\"")
			d, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("not a valid numebr.")
				return false
			}

			conf.reqCap = d
			i++

		case "-D":
			if len(remainingArgs) <= i+1 {
				fmt.Println("please include a value for -D.")
				return false
			}

			var temp string

			val, err := store.Get(remainingArgs[i+1])
			temp = val

			if err != nil {
				temp = remainingArgs[i+1]
			}

			cleanedJson := strings.Trim(temp, "\"") // remove newline slashes
			if ok := IsValidJson([]byte(cleanedJson)); !ok {
				fmt.Println("invalid json format.")
				return false
			}

			conf.reqBody = strings.NewReader(cleanedJson)
			i++

		default:
			fmt.Println("not a valid command.")
			continue
		}
	}
	return true
}

// help cmd func

func DisplayCmds() {
	fmt.Println("usage:")
	fmt.Println("var <VarName> <Value>")
	fmt.Println("to retrieve values:")
	fmt.Println("get <VarName>")
	fmt.Println("to test API's:")
	fmt.Println("test <VarName> -H content-type:application -X post -D { 'name': 'name' } or -X post -F <formTitle:formValue>")
	fmt.Println("though flags shown above are optional, but make sure <VarName> is a valid url.")
}
