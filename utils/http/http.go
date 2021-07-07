package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response is an HTTP response
type Response struct {
	StatusCode int
	Body       []byte
	JSON       interface{}
	Headers    http.Header
}

// Authentication contains the keys needed to build an authorization URL
type Authentication struct {
	Scheme   string
	Username string
	Password string
	Token    string
}

// Request is an HTTP request struct
type Request struct {
	URL            string `default:"http://localhost/"`
	Method         string `default:"GET"`
	Headers        map[string]string
	Query          map[string]string
	ContentType    string `default:"application/json"`
	Authentication Authentication
	Client         *http.Client
	Body           []byte
}

// Post isn an HTTP POST Method
func (request *Request) Post() (*Response, error) {
	request.Method = "POST"
	return request.http()
}

// Get is an HTTP GET Method
func (request *Request) Get() (*Response, error) {
	request.Method = "GET"
	return request.http()
}

// Put is an HTTP PUT Method
func (request *Request) Put() (*Response, error) {
	request.Method = "PUT"
	return request.http()
}

// Patch is an HTTP PATCH Method
func (request *Request) Patch() (*Response, error) {
	request.Method = "PATCH"
	return request.http()
}

// Delete is an HTTP DELETE Method
func (request *Request) Delete() (*Response, error) {
	request.Method = "DELETE"
	return request.http()
}

func (request *Request) http() (*Response, error) {

	var client *http.Client
	if request.Client != nil {
		client = request.Client
	} else {
		client = &http.Client{}
	}

	req, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(request.Body))
	if err != nil {
		return nil, err
	}

	if request.ContentType == "" {
		req.Header.Set("Content-type", "application/json")
	} else {
		req.Header.Set("Content-type", request.ContentType)
	}

	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}

	switch strings.ToLower(request.Authentication.Scheme) {
	case "basic":
		req.Header.Set("Authorization", request.Authentication.Scheme+" "+base64.StdEncoding.EncodeToString([]byte(request.Authentication.Username+":"+request.Authentication.Password)))
	case "bearer":
		req.Header.Set("Authorization", request.Authentication.Scheme+" "+request.Authentication.Token)
	}

	q := req.URL.Query()
	for k, v := range request.Query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: " + err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{resp.StatusCode, body, "", resp.Header}, fmt.Errorf("failed to read the response body: " + err.Error())
	}

	var decoded interface{}
	json.Unmarshal(body, &decoded)

	if resp.StatusCode > 399 {
		return &Response{resp.StatusCode, body, decoded, resp.Header}, fmt.Errorf("non-OK response: %s", body)
	}

	return &Response{resp.StatusCode, body, decoded, resp.Header}, nil
}
