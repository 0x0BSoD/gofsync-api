package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// =====================
//  Structures and vars
// =====================
var (
	host  string
	hosts []byte
	count string
)

type Auth struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

// =====================
//  Helpers
// =====================
// For pretty format output
func giveMeSpaces(num int) string {
	spaces := " "
	for i := 0; i < num; i++ {
		spaces += " "
	}
	return spaces
}

// ReturnHelp return help obviously
func ReturnHelp() {
	msg := `Usage:`

	fmt.Println(msg)
}

func getError(err string) {
	ReturnHelp()
	fmt.Println("Error:", err)
	os.Exit(1)
}

// SplitArg for splitting by '=' symbol
func SplitArg(r rune) bool {
	return r == '='
}

// =====================
//  Functions
// =====================
// CheckArgs return parsed parameters
func CheckArgs(args []string) {

	CountSet := false

	if len(args) == 1 {
		getError("Host not specified")
	}

	host = args[1]

	f, err := os.Open(host)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Not file: %v\n", err)
	}
	hosts, _ = ioutil.ReadAll(f)

	for _, argument := range args[2:] {

		a := strings.FieldsFunc(argument, SplitArg)
		arg, val := a[0], a[1]

		switch arg {
		case "-c", "--cont":
			CountSet = true
		}

		if CountSet {
			count = val
			continue
		}
	}

	if !CountSet {
		count = "10"
	}

}

// Return Auth structure with Username and Password for Foreman api
func configParser(path string) Auth {
	var auth Auth
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &auth)
	defer jsonFile.Close()
	return auth
}

func main() {
	CheckArgs(os.Args)

	//dbActions()

	if len(hosts) > 1 {
		sHosts := strings.Split(string(hosts), "\n")
		for _, host := range sHosts {
			if !strings.HasPrefix(host, "#") {
				log.Println(host)
				getHostgroups(host)
			}
		}
	} else {
		getHostgroups(host)
	}
}
