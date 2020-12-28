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
	"testing"
)

// TestFetch
func TestFetch(t *testing.T) {
	t.Run("Fetch", testFetch)
}

// testFetch
func testFetch(t *testing.T) {
	endPoint := "https://fakerapi.it/api/v1/books?_quantity=1"                     // string
	method := "GET"                                                                // string
	headers := [][]string{{"Accept", "*/*"}, {"Content-Type", "application/json"}} // [][]string
	f, err := NewFetch(endPoint, method, headers, nil)
	if err != nil {
		t.Errorf("Error Initializing New Fetch: %v", err)
	}
	f.Execute("")
	f.Resolve()
	if f.Res.StatusCode != 200 {
		t.Errorf("Error Executing and Resolving Test Request: %v", f.Res.Status)
	}
}
