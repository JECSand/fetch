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
	"net/http"
)

// Promise struct to store executing thread ...
type Promise struct {
	Channel     chan *http.Response
	Error       chan error
	httpRequest *http.Request
}

// worker
func (p *Promise) worker(done chan *http.Response, doneErr chan error) {
	client := &http.Client{}
	resp, err := client.Do(p.httpRequest)
	doneErr <- err
	done <- resp
	<-done
	<-doneErr
}

// execute
func (p *Promise) execute() {
	done := make(chan *http.Response)
	doneErr := make(chan error)
	go p.worker(done, doneErr)
	p.Channel = done
	p.Error = doneErr
}

// dispatch
func dispatch(r *http.Request) *Promise {
	promise := Promise{httpRequest: r}
	promise.execute()
	return &(promise)
}
