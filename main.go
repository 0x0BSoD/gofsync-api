package main

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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

func getDeltatime(start time.Time) string {
	delta := time.Since(start)
	res := fmt.Sprint(delta.String())
	return res
}

func main() {
	_, _, _, webServer := CheckArgs(os.Args)

	s := spinner.New(spinner.CharSets[25], 100*time.Millisecond)
	st := time.Now()
	s.Suffix = " Creating DB..."
	s.Start()
	dbActions()
	s.Stop()
	s.FinalMSG = "Complete! Creating DB worked: " + getDeltatime(st) + "\n"

	if webServer {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})
		http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hi")
		})
		log.Fatal(http.ListenAndServe(":8081", nil))
	} else {
		host, hosts, count, _ := CheckArgs(os.Args)
		if len(hosts) > 1 {
			sHosts := strings.Split(string(hosts), "\n")
			for _, _host := range sHosts {
				if !strings.HasPrefix(_host, "#") {
					// INIT ============
					overallT := time.Now()
					fmt.Println(_host)
					fmt.Println("=============================")

					s.Restart()
					st = time.Now()
					s.Suffix = " Getting Locations..."
					s.Start()
					getLocations(_host)
					s.Stop()
					s.FinalMSG = "Complete! Getting Locations worked: " + getDeltatime(st) + "\n"

					s.Restart()
					st = time.Now()
					s.Suffix = " Getting Host Groups..."
					s.Start()
					getHostGroups(_host, count)
					s.Stop()
					s.FinalMSG = "Complete! Getting Host Groups worked: " + getDeltatime(st) + "\n"

					s.Restart()
					st = time.Now()
					s.Suffix = " Getting Puppet Classes..."
					s.Start()
					getPuppetClasses(_host, count)
					s.Stop()
					s.FinalMSG = "Complete! Getting Puppet Classes worked: " + getDeltatime(st) + "\n"

					s.Restart()
					st = time.Now()
					s.Suffix = " Filling Smart Classes table..."
					s.Start()
					InsertPuppetSmartClasses(_host)
					s.Stop()
					s.FinalMSG = "Complete! Filling Smart Classes table worked: " + getDeltatime(st) + "\n"

					s.Restart()
					st = time.Now()
					s.Suffix = " Filling Smart Classes Base table..."
					s.Start()
					InsertToOverridesBase(_host)
					s.Stop()
					s.FinalMSG = "Complete! Filling Smart Classes Base table worked: " + getDeltatime(st) + "\n"

					s.Restart()
					st = time.Now()
					s.Suffix = " Filling Smart Classes Overrides parameters table..."
					s.Start()
					InsertOverridesParameters(_host)
					s.Stop()
					s.FinalMSG = "Complete! Filling Smart Classes Overrides parameters table worked: " + getDeltatime(st) + "\n"

					fmt.Println()
					sOverall := getDeltatime(overallT)
					fmt.Println(_host)
					fmt.Println("Done by ", sOverall)
					fmt.Println()
				}

				fmt.Println("Actions for all instances")
				s.Restart()
				st = time.Now()
				s.Suffix = " Filling SWE table..."
				s.Start()
				fillTableSWEState()
				s.Stop()
				s.FinalMSG = "Complete! Filling SWE table worked: " + getDeltatime(st) + "\n"

				s.Restart()
				st = time.Now()
				s.Suffix = " Checking SWE..."
				s.Start()
				checkSWEState()
				s.Stop()
				s.FinalMSG = "Complete! Checking SWE worked: " + getDeltatime(st) + "\n"
			}
		} else {
			// INIT ============
			dbActions()
			getLocations(host)
			getHostGroups(host, count)
			getPuppetClasses(host, count)
			fillTableSWEState()
			checkSWEState()
			InsertToOverridesBase(host)
			InsertPuppetSmartClasses(host)
			InsertOverridesParameters(host)
		}
	}

}
