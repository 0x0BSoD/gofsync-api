package utils

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"strconv"
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IntegerInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Fast conv int to string
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

// String wit comma separator to []int
func Integers(s string) []int {
	var tmpInt []int
	ls := strings.Split(s, ",")
	for _, i := range ls {
		Int, _ := strconv.Atoi(i)
		tmpInt = append(tmpInt, Int)
	}
	return tmpInt
}

// ¯\＿(ツ)＿/¯
func Pager(totalPages int, perPage int) int {
	if perPage == 0 {
		perPage = 100
	}
	return totalPages/perPage + 1
}

func PrintJsonStep(step models.Step) string {
	str, err := json.Marshal(step)
	if err != nil {
		Error.Println("Error on printing step")
	}
	return string(str)
}
