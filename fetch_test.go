/*
Author: John Connor Sanders
License: Apache Version 2.0
Version: 0.0.2
Released: 01/29/2021
Copyright (c) 2021 John Connor Sanders

-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
----------------FETCH--------------------
-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-
*/

package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

// jsonData
type jsonData struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// fileData
type fileData struct {
	Id      string
	Name    string
	Type    string
	Content []byte
}

// testDataCache
type testDataCache struct {
	JSONData []*jsonData
	FileData []*fileData
}

// getJSONData - returns all jsonData in the testDataCache
func (tc *testDataCache) getJSONData() []*jsonData {
	return tc.JSONData
}

// getJSONDatum - returns a specific jsonData entry from the testDataCache
func (tc *testDataCache) getJSONDatum(id string) (*jsonData, error) {
	for _, td := range tc.JSONData {
		if td.Id == id {
			return td, nil
		}
	}
	return &jsonData{}, errors.New(id + " Not found!")
}

// postJSONData - adds a new jsonData entry to the testDataCache
func (tc *testDataCache) postJSONData(td *jsonData) {
	td.Id = strconv.Itoa(len(tc.JSONData))
	tc.JSONData = append(tc.JSONData, td)
}

// getFileData - returns a fileData entry from the testDataCache
func (tc *testDataCache) getFileData(id string) (*fileData, error) {
	for _, fd := range tc.FileData {
		if fd.Id == id {
			return fd, nil
		}
	}
	return &fileData{}, errors.New(id + " Not found!")
}

// postFileData - adds a new fileData entry to the testDataCache
func (tc *testDataCache) postFileData(fd *fileData) {
	fd.Id = strconv.Itoa(len(tc.FileData))
	tc.FileData = append(tc.FileData, fd)
}

// newTestDataCache returns a pointer to a new testDataCache struct
func newTestDataCache(td []*jsonData, fd []*fileData) *testDataCache {
	return &testDataCache{
		td,
		fd,
	}
}

// testErr
type testErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// testRes
type testRes struct {
	Code int    `json:"code"`
	Id   string `json:"id"`
}

// TestFetch...
func TestFetch(t *testing.T) {
	t.Run("JSONGet", testJSONGet)
	t.Run("JSONPost", testJSONPost)
	t.Run("FileGet", testFileGet)
	t.Run("FilePost", testFilePost)
}

// testJSONGet
func testJSONGet(t *testing.T) {
	fmt.Println("\n<-----------BEGINNING 1 of 4: testJSONGet...")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			fmt.Println("Expected 'GET' received: ", r.Method)
			t.Errorf("Expected 'GET' received: '%s'", r.Method)
		}
		if r.URL.EscapedPath() != "/data/1" {
			fmt.Println("Incorrect url endpoint: ", r.URL.EscapedPath())
			t.Errorf("Incorrect url endpoint: '%s'", r.URL.EscapedPath())
		}
		dataCache := newTestDataCache([]*jsonData{{Id: "1", Name: "test", Type: "GET"}}, []*fileData{})
		id := strings.TrimPrefix(r.URL.Path, "/data/")
		td, err := dataCache.getJSONDatum(id)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			tErr := testErr{Code: http.StatusNotFound, Text: err.Error()}
			if err = json.NewEncoder(w).Encode(tErr); err != nil {
				fmt.Println("Error encoding testErr struct into JSON for response: ", err.Error())
				t.Errorf("Error encoding testErr struct into JSON for response: '%s'", err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(*td); err != nil {
			fmt.Println("Error encoding jsonData struct into JSON for response: ", err.Error())
			t.Errorf("Error encoding jsonData struct into JSON for response: '%s'", err.Error())
		}
		return
	}))
	defer ts.Close()
	api := ts.URL
	endPoint := fmt.Sprintf("%s/data/%s", api, "1")
	fmt.Printf("TEST Url: GET %s\n", endPoint)
	method := "GET"
	f, err := NewFetch(endPoint, method, JSONDefaultHeaders(), nil)
	if err != nil {
		fmt.Println("Error Initializing New Fetch: ", err.Error())
		t.Errorf("Error Initializing New Fetch: %v", err.Error())
	}
	err = f.Execute("")
	if err != nil {
		fmt.Println("Error Executing Fetch Request: ", err.Error())
		t.Errorf("Error Executing Fetch Request: %v", err.Error())
	}
	f.Resolve()
	if f.Res.StatusCode != 200 {
		fmt.Println("Error Executing and Resolving Test Request: ", f.Res.Status)
		t.Errorf("Error Executing and Resolving Test Request: %v", f.Res.Status)
	}
	fmt.Println("<-----------COMPLETE 1 of 4: testJSONGet")
}

// testJSONPost
func testJSONPost(t *testing.T) {
	fmt.Println("\n<-----------BEGINNING 2 of 4: testJSONPost...")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Println("Expected 'POST' received: ", r.Method)
			t.Errorf("Expected 'POST' received: '%s'", r.Method)
		}
		if r.URL.EscapedPath() != "/data" {
			fmt.Println("Incorrect url endpoint: ", r.URL.EscapedPath())
			t.Errorf("Incorrect url endpoint: '%s'", r.URL.EscapedPath())
		}
		var td jsonData
		dataCache := newTestDataCache([]*jsonData{}, []*fileData{})
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			fmt.Println("Error encoding json post request body into jsonData struct: ", err.Error())
			t.Errorf("Error encoding json post request body into jsonData struct: '%s'", err.Error())
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err = json.Unmarshal(body, &td); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			tErr := testErr{Code: http.StatusUnprocessableEntity, Text: err.Error()}
			if err = json.NewEncoder(w).Encode(tErr); err != nil {
				fmt.Println("Error encoding testErr struct into JSON for response: ", err.Error())
				t.Errorf("Error encoding testErr struct into JSON for response: '%s'", err.Error())
			}
			return
		}
		dataCache.postJSONData(&td)
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(td); err != nil {
			fmt.Println("Error encoding error into JSON for response: ", err.Error())
			t.Errorf("Error encoding error into JSON for response: '%s'", err.Error())
		}
		return
	}))
	defer ts.Close()
	api := ts.URL
	endPoint := fmt.Sprintf("%s/data", api)
	fmt.Printf("TEST Url: POST %s\n", endPoint)
	method := "POST"
	body := []byte(`{"name":"test","type":"POST"}`)
	f, err := NewFetch(endPoint, method, JSONDefaultHeaders(), bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error Initializing New Fetch: ", err.Error())
		t.Errorf("Error Initializing New Fetch: %v", err.Error())
	}
	err = f.Execute("")
	if err != nil {
		fmt.Println("Error Executing Fetch Request: ", err.Error())
		t.Errorf("Error Executing Fetch Request: %v", err.Error())
	}
	f.Resolve()
	if f.Res.StatusCode != 201 {
		fmt.Println("Error Executing and Resolving Test Request: ", f.Res.Status)
		t.Errorf("Error Executing and Resolving Test Request: %v", f.Res.Status)
	}
	fmt.Println("<-----------COMPLETE 2 of 4: testJSONPost")
}

// testFileGet
func testFileGet(t *testing.T) {
	fmt.Println("\n<-----------BEGINNING 3 of 4: testFileGet...")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			fmt.Println("Expected 'GET' received: ", r.Method)
			t.Errorf("Expected 'GET' received: '%s'", r.Method)
		}
		if r.URL.EscapedPath() != "/files/1" {
			fmt.Println("Incorrect url endpoint: ", r.URL.EscapedPath())
			t.Errorf("Incorrect url endpoint: '%s'", r.URL.EscapedPath())
		}
		fc := []byte("This is the GET test file's contents!")
		dataCache := newTestDataCache([]*jsonData{}, []*fileData{{Id: "1", Name: "test.txt", Type: "txt", Content: fc}})
		id := strings.TrimPrefix(r.URL.Path, "/files/")
		fd, err := dataCache.getFileData(id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			tErr := testErr{Code: http.StatusNotFound, Text: err.Error()}
			if err = json.NewEncoder(w).Encode(tErr); err != nil {
				fmt.Println("Error encoding testErr struct into JSON for response: ", err.Error())
				t.Errorf("Error encoding testErr struct into JSON for response: '%s'", err.Error())
			}
			return
		}
		modTime := time.Now()
		cd := mime.FormatMediaType("attachment", map[string]string{"filename": fd.Name})
		w.Header().Set("Content-Disposition", cd)
		w.Header().Set("Content-Type", "application/octet-stream")
		bReader := bytes.NewReader(fd.Content)
		http.ServeContent(w, r, fd.Name, modTime, bReader)
	}))
	defer ts.Close()
	api := ts.URL
	endPoint := fmt.Sprintf("%s/files/%s", api, "1")
	fmt.Printf("TEST Url: GET %s\n", endPoint)
	method := "GET"
	f, err := NewFetch(endPoint, method, JSONDefaultHeaders(), nil)
	if err != nil {
		fmt.Println("Error Initializing New Fetch: ", err.Error())
		t.Errorf("Error Initializing New Fetch: %v", err.Error())
	}
	err = f.Execute("")
	if err != nil {
		fmt.Println("Error Executing Fetch Request: ", err.Error())
		t.Errorf("Error Executing Fetch Request: %v", err.Error())
	}
	f.Resolve()
	if f.Res.StatusCode != 200 {
		fmt.Println("Error Executing and Resolving Test Request: ", f.Res.Status)
		t.Errorf("Error Executing and Resolving Test Request: %v", f.Res.Status)
	}
	d, err := ioutil.ReadAll(f.Res.Body)
	if err != nil {
		fmt.Println("Error reading response file body: ", err.Error())
		t.Errorf("Error reading response file body: %v", err.Error())
	}
	if string(d) != "This is the GET test file's contents!" {
		fmt.Println("Expected file contents and returned do not match: ", string(d))
		t.Errorf("Expected file contents and returned do not match: %v", string(d))
	}
	fmt.Println("<-----------COMPLETE 3 of 4: testFileGet")
}

// testFilePost
func testFilePost(t *testing.T) {
	fmt.Println("\n<-----------BEGINNING 4 of 4: testFilePost...")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Println("Expected 'POST' received: ", r.Method)
			t.Errorf("Expected 'POST' received: '%s'", r.Method)
		}
		if r.URL.EscapedPath() != "/files" {
			fmt.Println("Incorrect url endpoint: ", r.URL.EscapedPath())
			t.Errorf("Incorrect url endpoint: '%s'", r.URL.EscapedPath())
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Println("Error parsing multipart form", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			tErr := testErr{Code: http.StatusBadRequest, Text: err.Error()}
			if err = json.NewEncoder(w).Encode(tErr); err != nil {
				fmt.Println("Error encoding testErr struct into JSON for response: ", err.Error())
				t.Errorf("Error encoding testErr struct into JSON for response: '%s'", err.Error())
			}
			return
		}
		dataCache := newTestDataCache([]*jsonData{}, []*fileData{})
		var tf fileData
		file, h, err := r.FormFile("file")
		if err != nil {
			fmt.Println("Error getting form file", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			tErr := testErr{Code: http.StatusBadRequest, Text: err.Error()}
			if err = json.NewEncoder(w).Encode(tErr); err != nil {
				fmt.Println("Error encoding testErr struct into JSON for response: ", err.Error())
				t.Errorf("Error encoding testErr struct into JSON for response: '%s'", err.Error())
			}
			return
		}
		defer file.Close()
		buf := bytes.NewBuffer(nil)
		if _, err = io.Copy(buf, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			tErr := testErr{Code: http.StatusInternalServerError, Text: err.Error()}
			if err = json.NewEncoder(w).Encode(tErr); err != nil {
				fmt.Println("Error encoding testErr struct into JSON for response: ", err.Error())
				t.Errorf("Error encoding testErr struct into JSON for response: '%s'", err.Error())
			}
			return
		}
		tf.Name = h.Filename
		tf.Content = buf.Bytes()
		tf.Type = "txt"
		if strings.Contains(h.Filename, ".") {
			tf.Type = strings.Split(h.Filename, ".")[1]
		}
		dataCache.postFileData(&tf)
		w.WriteHeader(http.StatusOK)
		res := testRes{http.StatusOK, tf.Id}
		if err = json.NewEncoder(w).Encode(res); err != nil {
			fmt.Println("Error encoding testRes struct into JSON for response: ", err.Error())
			t.Errorf("Error encoding testRes struct into JSON for response: '%s'", err.Error())
		}
		return
	}))
	defer ts.Close()
	api := ts.URL
	endPoint := fmt.Sprintf("%s/files", api)
	fmt.Printf("TEST Url: POST %s\n", endPoint)
	method := "POST"
	fContent := []byte("This is the POST test file's contents!")
	f, err := NewFileFetch("test.txt", endPoint, method, DefaultHeaders(), bytes.NewBuffer(fContent))
	if err != nil {
		fmt.Println("Error Initializing New File Fetch: ", err.Error())
		t.Errorf("Error Initializing New File Fetch: %v", err.Error())
	}
	err = f.Execute("")
	if err != nil {
		fmt.Println("Error Executing File Fetch Request: ", err.Error())
		t.Errorf("Error Executing File Fetch Request: %v", err.Error())
	}
	f.Resolve()
	if f.Res.StatusCode != 200 {
		fmt.Println("Error Executing and Resolving Test Request: ", f.Res.Status)
		t.Errorf("Error Executing and Resolving Test Request: %v", f.Res.Status)
	}
	fmt.Println("<-----------COMPLETE 4 of 4: testFilePost")
}
