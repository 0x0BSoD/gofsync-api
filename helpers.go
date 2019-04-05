package main

import (
	"crypto/tls"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/logger"
	"github.com/spf13/viper"
	"gopkg.in/ldap.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func ldapTest(username string, password string) (string, error) {
	// The username and password we want to check
	bindUsername := "srv_foreman"
	bindPassword := "2R2gUres"

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "ams01-c01-adc02.ringcentral.com", 389))
	if err != nil {
		return "", err
	}
	defer l.Close()

	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return "", err
	}

	// First bind with a read only user
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		return "", err
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		"ou=RingCentral Admins,dc=ringcentral,dc=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", username),
		[]string{"*"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", err
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("user does not exist or too many entries returned")
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		return "", err
	}

	// Rebind as the read only user for any further queries
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		return "", err
	}

	return userdn, nil
}

func getHosts(file string) {
	if len(file) > 0 {
		// Get hosts from file
		var hosts []byte
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Not file: %v\n", err)
		}
		hosts, _ = ioutil.ReadAll(f)
		tmpHosts := strings.Split(string(hosts), "\n")
		var sHosts []string
		for _, i := range tmpHosts {
			if !strings.HasPrefix(i, "#") && len(i) > 0 {
				sHosts = append(sHosts, i)
			}
		}
		globConf.Hosts = sHosts
	} else {
		logger.Error.Println("Hosts file not found...")
		os.Exit(2)
	}
}

func configParser() {
	viper.SetConfigName(conf)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error.Println("Config file not found...")
		os.Exit(2)
	} else {
		globConf = Config{
			Username: viper.GetString("API.username"),
			Pass:     viper.GetString("API.password"),
			DBFile:   viper.GetString("DB.db_file"),
			Actions:  viper.GetStringSlice("RUNNING.actions"),
			RTPro:    viper.GetString("RT.pro"),
			RTStage:  viper.GetString("RT.stage"),
			PerPage:  viper.GetInt("RUNNING.per_page_def"),
			DbInit:   viper.GetString("DB.init_file"),
		}
		globConf.Initialize(viper.GetString("DB.db_user"),
			viper.GetString("DB.db_password"),
			viper.GetString("DB.db_schema"))
	}
}

// Split []int to [](parts * []int)
func splitToQueue(item []int, parts int) [][]int {

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

func stringInSlice(a string, list []string) bool {
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
func Pager(totalPages int) int {
	pagesRange := totalPages/globConf.PerPage + 1
	return pagesRange
}
