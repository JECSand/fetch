# fetch

A pre-made Golang module for making easy async REST calls.

[![Build Status](https://travis-ci.org/JECSand/fetch.svg?branch=main)](https://travis-ci.org/JECSand/fetch)
[![Go Report Card](https://goreportcard.com/badge/github.com/JECSand/fetch)](https://goreportcard.com/report/github.com/JECSand/fetch)

* Author(s): John Connor Sanders
* License: Apache Version 2.0
* Version Release Date: 06/17/2021
* Current Version: 0.0.3

## License
* Copyright 2021 John Connor Sanders

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
	"bytes"
	"github.com/JECSand/fetch"
	"log"
)


func main() {
	
	// 1: Setup Variable Parameters for Requests.
	endPoint := "https://fakerapi.it/api/v1/books?_quantity=1" // string
	method := "GET" // string
	endPoint2 := "https://fakerapi.it/api/v1/companies?_quantity=5" // string
	method2 := "GET" // string
	//**In cases where the Body parameter is not nil, it can be setup as follows:
	// body := bytes.NewBuffer([]byte(`{"username":"test","password":"password123"}`))
	
	// 2: Initialize a new Fetch Structure with parameters.
	f, err := fetch.NewFetch(endPoint, method, fetch.JSONDefaultHeaders(), nil)
	if err != nil {
		log.Fatalf("Failed to initialize new Fetch Struct 1: %v", err)
	}
	log.Printf("Successfully initialized new Fetch Struct 1: %v", f)
	f2, err2 := fetch.NewFetch(endPoint2, method2, fetch.JSONDefaultHeaders(), nil)
	if err2 != nil {
		log.Fatalf("Failed to initialize new Fetch Struct 2: %v", err2)
	}
	log.Printf("Successfully initialized new Fetch Struct 2: %v", f2)
	
	// 3: Initialize a new multipart File Fetch 
	method = "POST"
	fContent := []byte("This is the POST file's contents!")
	f3, err := fetch.NewFileFetch("test.txt", endPoint, method, fetch.DefaultHeaders(), bytes.NewBuffer(fContent))
	if err != nil {
		log.Fatalf("Failed to initialize new Fetch Struct 3: %v", err)
	}
	log.Printf("Successfully initialized new Fetch Struct 3: %v", f3)
	
	// 4: Execute Fetch structs' Async Requests and store structs in a Slice of Fetch.
	f.Execute("") // **Optionally you can use "discard" instead of "" to throw the http response away.
	f2.Execute("") // **Ditto from above.
	f3.Execute("") // **Ditto from above.
	fProcesses := []*fetch.Fetch{f, f2, f3}
	
	// 5: Resolve Fetch structs' as needed.
	fProcesses[0].Resolve()
	fProcesses[1].Resolve()
	fProcesses[2].Resolve()
	
	// 6: Access *http.Response in Fetch structs.
	log.Printf("Successfully resolved Fetch Struct 1: %v", f.Res)
	log.Printf("Successfully resolved Fetch Struct 2: %v", f2.Res)
	log.Printf("Successfully resolved Fetch Struct 3: %v", f3.Res)
}
```
