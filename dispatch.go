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
	"net/http"
)

// Promise struct to store executing thread ...
type Promise struct {
	Channel     chan *http.Response
	httpRequest *http.Request
}

// worker
func (p *Promise) worker(done chan *http.Response) {
	client := &http.Client{}
	resp, _ := client.Do(p.httpRequest)
	done <- resp
	<-done
}

// execute
func (p *Promise) execute() {
	done := make(chan *http.Response)
	go p.worker(done)
	p.Channel = done
}

// dispatch
func dispatch(r *http.Request) *Promise {
	promise := Promise{httpRequest: r}
	promise.execute()
	return &(promise)
}
