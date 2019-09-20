package utils

import (
	"database/sql"
	"encoding/json"
	"git.ringcentral.com/archops/goFsync/models"
	"net/http"
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

// IntegerInSlice  replacement
func Search(data []int, s int) bool {
	if len(data) > 1 {
		first := 0
		last := len(data) - 1

		for first <= last {
			middle := (first + last) / 2

			if data[middle] < s {
				first = middle + 1
			} else {
				last = middle - 1
			}
		}

		if last == len(data) || data[first] != s {
			return false
		} else {
			return true
		}
	} else {
		return false
	}

}

// Fast conv int to string
func String(n int) string {
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

func DeferCloseStmt(conn *sql.Stmt) {
	if err := conn.Close(); err != nil {
		Error.Println("Error on closing DB connection")
	}
}

func SendResponse(w http.ResponseWriter, msg string, data interface{}) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		Error.Printf(msg, err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
	}
}
