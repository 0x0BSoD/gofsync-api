package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/0x0bsod/foremanGetter/entitys"
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
var (
	host  string
	count string
)

type Auth struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

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

// CheckArgs return parsed parameters
func CheckArgs(args []string) {

	CountSet := false

	if len(args) == 1 {
		getError("Host not specified")
	}

	host = args[1]

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

func worker() {

	spaces := 10
	var result []entitys.SWEs

	fmt.Printf("Getting from %s \n", host)

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.CloseIdleConnections()

	client := &http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", "https://"+host+"/api/hostgroups?format=json&per_page="+count, nil)
	auth := configParser("./config.json")
	req.SetBasicAuth(auth.Username, auth.Pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	bodyText, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal([]byte(bodyText), &result)
	if err != nil {
		log.Printf("%q:\n %s\n", err, bodyText)
		return
	}

	for _, item := range result {
		if item.Hostgroup.Name != "SWE" {
			fmt.Println("Name             :  ", item.Hostgroup.Name)
			fmt.Println("ID               :  ", item.Hostgroup.ID)
			fmt.Println("SubnetID         :  ", item.Hostgroup.SubnetID)
			fmt.Println("OperatingsystemID:  ", item.Hostgroup.OperatingsystemID)
			fmt.Println("DomainID         :  ", item.Hostgroup.DomainID)
			fmt.Println("EnvironmentID    :  ", item.Hostgroup.EnvironmentID)
			fmt.Println("Ancestry         :  ", item.Hostgroup.Ancestry)

			if len(item.Hostgroup.Parameters) > 1 {
				fmt.Println("Parameters       :=>  ")
				for name, item := range item.Hostgroup.Parameters {
					length := len(name)
					strSpaces := giveMeSpaces(spaces - length)
					fmt.Printf("    %s%s:  %s\n", name, strSpaces, item)
				}
			} else {
				fmt.Println("Parameters       :   nil")
			}

			fmt.Println("PuppetclassIds   :  ", item.Hostgroup.PuppetclassIds)
			fmt.Println("==================================================")
		}

	}
}

func main() {
	CheckArgs(os.Args)
	worker()
}
