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

// AppendHeaders ...
func AppendHeaders(headers [][]string, newHeader []string) [][]string {
	var updatedHeaders [][]string
	updated := false
	for _, header := range headers {
		if header[0] == newHeader[0] {
			updatedHeaders = append(updatedHeaders, newHeader)
			updated = true
		} else {
			updatedHeaders = append(updatedHeaders, header)
		}
	}
	if !updated {
		updatedHeaders = append(updatedHeaders, newHeader)
	}
	return updatedHeaders
}

// DefaultHeaders ...
func DefaultHeaders() [][]string {
	var headers [][]string
	headerEntries := [][]string{{"Accept", "*/*"}}
	headers = append(headers, headerEntries...)
	return headers
}

// JSONDefaultHeaders ...
func JSONDefaultHeaders() [][]string {
	var headers [][]string
	headerEntries := [][]string{{"Accept", "*/*"}, {"Content-Type", "application/json"}}
	headers = append(headers, headerEntries...)
	return headers
}
