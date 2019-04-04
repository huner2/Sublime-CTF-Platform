package main

import (
	"math/rand"
	"strings"
)

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func contains(s []string, e string) (bool, string) {
	for _, a := range s {
		if strings.ToLower(a) == strings.ToLower(e) {
			return true, a
		}
	}
	return false, ""
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
