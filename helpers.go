package main

import (
	"database/sql"
	"fmt"
	"github.com/briandowns/spinner"
	"log"
	"time"
)

type toWaiter func(host string, count string)

func getDBConn() *sql.DB {

	db, err := sql.Open("sqlite3", globConf.DBFile)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

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
	s.FinalMSG = fmt.Sprintf("Complete! %s worked: %s\n", msg, getDeltaTime(st))

	s.Start()
	f(host, count)
	s.Stop()

}

func String(n int64) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
