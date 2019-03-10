package main

import (
	"fmt"
	"github.com/briandowns/spinner"
	"time"
)

type toWaiter func(host string, count string)

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

func Spinner(msg string, host string, count string, f toWaiter) {
	s := spinner.New(spinner.CharSets[21], 100*time.Millisecond)
	st := time.Now()

	s.Suffix = fmt.Sprintf(" %s...", msg)
	s.FinalMSG = fmt.Sprintf("Complete! %s worked: %s\n", msg,  getDeltaTime(st))

	s.Start()
	f(host, count)
	s.Stop()

}
