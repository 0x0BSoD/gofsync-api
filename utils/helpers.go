package utils

import (
	"encoding/json"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"strconv"
	"strings"
)

// Split []int to [](parts * []int)
func SplitToQueue(item []int, parts int) [][]int {

	var result [][]int
	length := len(item)
	sliceLength := 0

	// Checks ==========
	if length <= parts {
		return append(result, item)
	}

	if length%parts == 0 {
		sliceLength = length / parts
	} else {
		if length/parts == 1 {
			return append(result, item)
		}
		sliceLength = (length / parts) + 1
	}
	if sliceLength == 1 {
		return append(result, item)
	}

	start := 0
	stop := sliceLength

	for i := 0; i < parts; i++ {
		if stop < length {
			result = append(result, item[start:stop])
		} else {
			result = append(result, item[start:])
			break
		}
		start = stop
		stop = start + sliceLength
	}
	return result
}

func StringInSlice(a string, list []string) bool {
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
	return totalPages/perPage + 1
}

func PrintJsonStep(step models.Step) string {
	str, err := json.Marshal(step)
	if err != nil {
		Error.Println("Error on printing step")
	}
	return string(str)
}
