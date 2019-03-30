package main

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

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
		log.Fatal("")
	}
}

func configParser() {
	viper.SetConfigName(conf)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
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
