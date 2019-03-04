package main

import (
	"fmt"
	"time"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getDeltaTime(start time.Time) string {
	delta := time.Since(start)
	res := fmt.Sprint(delta.String())
	return res
}
