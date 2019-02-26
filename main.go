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
type Auth struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

// =====================
//  Helpers
// =====================
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
func CheckArgs(args []string) (string, []byte, string) {

	var (
		host  string
		hosts []byte
		count string
	)
	CountSet := false

	if len(args) == 1 {
		getError("Host not specified")
	}

	host = args[1]

	f, err := os.Open(host)
	if err != nil {
		log.Fatalf("Not file: %v\n", err)
	}
	hosts, _ = ioutil.ReadAll(f)

	for _, argument := range args[2:] {

		a := strings.FieldsFunc(argument, SplitArg)
		arg, val := a[0], a[1]

		switch arg {
		case "-c", "--cont":
			CountSet = true
		}

		// Flag checker
		if CountSet {
			count = val
			continue
		}

	}

	// Default values
	if !CountSet {
		count = "10"
	}

	return host, hosts, count
}

// Return Auth structure with Username and Password for Foreman api
func configParser(path string) Auth {

	var auth Auth

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &auth)

	return auth
}

func main() {

	host, hosts, count := CheckArgs(os.Args)
	if len(hosts) > 1 {
		sHosts := strings.Split(string(hosts), "\n")
		for _, _host := range sHosts {
			if !strings.HasPrefix(_host, "#") {
				dbActions()
				checkSWEState()
				//fillTableSWEState()
				//getHostGroups(_host, count)
				//getPuppetClasses(_host, count)
				//getLocations(_host)
			}
		}
	} else {
		dbActions()
		fillTableSWEState()
		getHostGroups(host, count)
		getPuppetClasses(host, count)
		getLocations(host)
	}
}
