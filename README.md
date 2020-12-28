# fetch

A pre-made Golang module for making easy async REST calls.

[![Build Status](https://travis-ci.org/JECSand/fetch.svg?branch=master)](https://travis-ci.org/JECSand/fetch)
[![Go Report Card](https://goreportcard.com/badge/github.com/JECSand/fetch)](https://goreportcard.com/report/github.com/JECSand/fetch)

* Author(s): John Connor Sanders
* License: Apache Version 2.0
* Version Release Date: 12/28/2020
* Current Version: 0.0.1

## License
* Copyright 2020 John Connor Sanders

This source code of this package is released under the Apache Version 2.0 license. Please see
the [LICENSE](https://github.com/JECSand/fetch/blob/main/LICENSE) for the full
content of the license.

## Installation
```bash
$ go get github.com/JECSand/fetch
```

## Usage
```go
package main

import (
	"github.com/JECSand/fetch"
	"log"
)


func main() {
	
	// 1: Setup Variable Parameters for Requests.
	endPoint := "https://fakerapi.it/api/v1/books?_quantity=1" // string
	method := "GET" // string
	headers := [][]string{[]string{"Accept", "*/*"}, []string{"Content-Type", "application/json"}} // [][]string
	endPoint2 := "https://fakerapi.it/api/v1/companies?_quantity=5" // string
	method2 := "GET" // string
	headers2 := [][]string{[]string{"Accept", "*/*"}, []string{"Content-Type", "application/json"}} // [][]string
	//*In cases where the Body parameter is not nil, it can be setup as follows:
	//body := []byte(`{"username":"test","password":"password123"}`)
	
	// 2: Initialize a new Fetch Structure with parameters.
	f, err := fetch.NewFetch(endPoint, method, headers, nil)
	if err != nil {
		log.Fatalf("Failed to initialize new Fetch Struct 1: %v", err)
    }
	log.Printf("Successfully initialized new Fetch Struct 1: %v", f)
	f2, err2 := fetch.NewFetch(endPoint2, method2, headers2, nil)
	if err2 != nil {
		log.Fatalf("Failed to initialize new Fetch Struct 2: %v", err2)
	}
	log.Printf("Successfully initialized new Fetch Struct 2: ", f2)
	
	// 3: Execute Fetch structs' Async Requests and store structs in a Slice of Fetch.
	f.Execute("") // *Optionally you can use "discard" instead of "" to throw the http response away
	f2.Execute("") // *Ditto
	fProcesses := []*fetch.Fetch{f, f2}
	
	// 4: Resolve Fetch structs' as needed
	fProcesses[0].Resolve()
	fProcesses[1].Resolve()
	
	// 5: Access *http.Response in Fetch structs'
	log.Printf("Successfully resolved Fetch Struct 1: ", f.Res)
	log.Printf("Successfully resolved Fetch Struct 2: ", f2.Res)
}
```
