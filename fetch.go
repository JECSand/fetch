/*
Author: John Connor Sanders
License: Apache Version 2.0
Version: 0.0.3
Released: 06/17/2021
Copyright (c) 2021 John Connor Sanders

-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
----------------FETCH--------------------
-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
*/

package fetch

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

// A Request ...
type Request struct {
	Headers [][]string
	Body    io.Reader
	Type    string
}

// defaultHeaders
func (re *Request) defaultHeaders() {
	var headers [][]string
	headerEntries := [][]string{{"Accept", "*/*"}}
	headers = append(headers, headerEntries...)
	re.Headers = headers
}

// Fetch ...
type Fetch struct {
	URL     string
	Req     *Request
	Res     *http.Response
	Promise *Promise
	Error   error
}

// NewFetch ...
func NewFetch(url string, method string, headers [][]string, body io.Reader) (*Fetch, error) {
	if url == "" {
		return &Fetch{}, errors.New("error: URL String Required")
	} else if method == "" {
		return &Fetch{}, errors.New("error: Method String Required")
	}
	d := Fetch{URL: url, Error: nil}
	d.Req = &Request{Headers: headers, Type: method, Body: body}
	if len(headers) == 0 {
		d.Req.defaultHeaders()
	}
	return &d, nil
}

// NewFileFetch ...
func NewFileFetch(fileName string, url string, method string, headers [][]string, body io.Reader) (*Fetch, error) {
	if fileName == "" {
		return &Fetch{}, errors.New("error: fileName String Required")
	}
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return &Fetch{}, err
	}
	_, err = io.Copy(fw, body)
	if err != nil {
		return &Fetch{}, err
	}
	_ = writer.Close()
	headerEntry := []string{"Content-Type", writer.FormDataContentType()}
	headers = AppendHeaders(headers, headerEntry)
	return NewFetch(url, method, headers, bytes.NewReader(buf.Bytes()))
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
		r, err = http.NewRequest(d.Req.Type, urlStr, d.Req.Body) // Body
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

// recursiveResolve
func (d *Fetch) recursiveResolve(attempt int) {
	if attempt < 4 {
		err := <-d.Promise.Error
		if err != nil {
			time.Sleep(5 * time.Second)
			d.Execute("")
			time.Sleep(10 * time.Second)
			attempt = attempt + 1
			d.recursiveResolve(attempt)
		} else {
			resp := <-d.Promise.Channel
			if resp == nil {
				time.Sleep(5 * time.Second)
				d.Execute("")
				time.Sleep(10 * time.Second)
				attempt = attempt + 1
				d.recursiveResolve(attempt)
			} else {
				d.Res = resp
			}
		}
	} else {
		d.Res = nil
	}
}

// Resolve Request
func (d *Fetch) Resolve() {
	d.recursiveResolve(0)
}
