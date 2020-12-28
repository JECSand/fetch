/*
Author: John Connor Sanders
License: Apache Version 2.0
Version: 0.0.1
Released: 12/28/2020
Copyright 2020 John Connor Sanders

-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
----------------FETCH--------------------
-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
*/

package fetch

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
)

// A Request ...
type Request struct {
	Headers [][]string
	Body    []byte
	Type    string
}

// defaultHeaders
func (re *Request) defaultHeaders() {
	var headers [][]string
	headerEntry := []string{"Content-Type", "application/json"}
	headers = append(headers, headerEntry)
	re.Headers = headers
}

// A Fetch ...
type Fetch struct {
	URL     string
	Req     Request
	Res     *http.Response
	Promise *Promise
}

// NewFetch ...
func NewFetch(url string, method string, headers [][]string, body []byte) (*Fetch, error) {
	if url == "" {
		return &Fetch{}, errors.New("Error: URL String Required")
	} else if method == "" {
		return &Fetch{}, errors.New("Error: Method String Required")
	}
	d := Fetch{URL: url}
	d.Req = Request{Headers: headers, Type: method, Body: body}
	if len(headers) == 0 {
		d.Req.defaultHeaders()
	}
	return &d, nil
}

// headers
func (d *Fetch) headers(r *http.Request) *http.Request {
	for _, headerStr := range d.Req.Headers {
		r.Header.Set(headerStr[0], headerStr[1])
	}
	return r
}

// Execute Request
func (d *Fetch) Execute(resType string) error {
	u, err := url.ParseRequestURI(d.URL)
	if err != nil {
		return err
	}
	urlStr := u.String()
	var r *http.Request
	if d.Req.Body == nil {
		r, err = http.NewRequest(d.Req.Type, urlStr, nil) // No Body
	} else {
		r, err = http.NewRequest(d.Req.Type, urlStr, bytes.NewBuffer(d.Req.Body)) // Body
	}
	if err != nil {
		return err
	}
	r = d.headers(r)
	reqPromise := dispatch(r)
	if resType == "discard" {
		<-reqPromise.Channel
	} else {
		d.Promise = reqPromise
	}
	return nil
}

// Resolve Request
func (d *Fetch) Resolve() {
	resp := <-d.Promise.Channel
	d.Res = resp
}
