package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
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
func CheckArgs(args []string) (string, []byte, string, bool) {

	var (
		host      string
		hosts     []byte
		count     string
		webServer bool
	)
	CountSet := false
	hostSet := false
	hostsSet := false

	for _, argument := range args[1:] {
		if strings.ContainsAny(argument, "=") { // PARAMETERS ===============
			a := strings.FieldsFunc(argument, SplitArg)
			arg, val := a[0], a[1]
			// Flag checker
			switch arg {
			case "-h", "--host":
				hostSet = true
			case "-f", "--file":
				hostsSet = true
			case "-c", "--cont":
				CountSet = true
			}
			// Flag with val getter
			if CountSet {
				count = val
				continue
			}
			if hostSet {
				host = val
				continue
			}
			if hostsSet {
				f, err := os.Open(val)
				if err != nil {
					log.Fatalf("Not file: %v\n", err)
				}
				hosts, _ = ioutil.ReadAll(f)
				continue
			}
		} else { // FLAGS ===============
			// Flag checker
			switch argument {
			case "-w", "--webserver":
				webServer = true
			}
		}

	}

	// Default values
	if !CountSet {
		count = "10"
	}

	return host, hosts, count, webServer
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
	_, _, _, webServer := CheckArgs(os.Args)
	if webServer {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})
		http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hi")
		})
		log.Fatal(http.ListenAndServe(":8081", nil))
	} else {
		host, hosts, _, _ := CheckArgs(os.Args)
		//host, hosts, count, _ := CheckArgs(os.Args)
		if len(hosts) > 1 {
			sHosts := strings.Split(string(hosts), "\n")
			for _, _host := range sHosts {
				if !strings.HasPrefix(_host, "#") {
					getAllPuppetSmartClasses(_host)
					//dbActions()
					//checkSWEState()
					//fillTableSWEState()
					//getHostGroups(_host, count)
					//getPuppetClasses(_host, count)
					//getLocations(_host)
				}
			}
		} else {
			getAllPuppetSmartClasses(host)
			//dbActions()
			//fillTableSWEState()
			//getHostGroups(host, count)
			//getPuppetClasses(host, count)
			//getLocations(host)
		}
	}

}

// https://spb01-puppet.lab.nordigy.ru/api/v2/smart_class_parameters/173/override_values host specific
