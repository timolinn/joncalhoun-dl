package main

import "net/http"

// RoundTripperFunc implements the http.RoundTripper interface
type RoundTripperFunc func(r *http.Request) (*http.Response, error)

// RoundTrip function
func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

// TestingHTTPClient returns http client with stubbed transport
// func TestingHTTPClient(fn RoundTripperFunc) *http.Client {
// 	return &http.Client{
// 		Transport: RoundTripperFunc(fn),
// 	}
// }
