package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"sync"
)

func addResult(i models.SCParameter, r *SmartClasses, mtx *sync.Mutex) {
	mtx.Lock()
	if i.OverrideValuesCount > 0 {
		fmt.Println(i.ID, i.OverrideValues)
		fmt.Println("====")
	}
	r.SmartClasses = append(r.SmartClasses, i)
	mtx.Unlock()
}

func worker(i int, host string, cfg *models.Config) models.SCParameter {
	var d models.SCParameter
	fmt.Printf("W: got task, scId: %d, HOST: %s\n", i, host)

	uri := fmt.Sprintf("smart_class_parameters/%d", i)
	response, _ := utils.ForemanAPI("GET", host, uri, "", cfg)
	if response.StatusCode != 200 {
		fmt.Println("SC Parameters, ID:", i, response.StatusCode, host)
	}

	err := json.Unmarshal(response.Body, &d)
	if err != nil {
		logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
	}
	return d
}

func asyncWorker(wrkID int,
	in <-chan int,
	out chan<- models.SCParameter,
	host string,
	wg *sync.WaitGroup,
	cfg *models.Config) {
	defer wg.Done()
	var d models.SCParameter
	for {
		i := <-in

		//fmt.Printf("W: %d got task, scId: %d, HOST: %s\n", wrkID, i, host)

		uri := fmt.Sprintf("smart_class_parameters/%d", i)
		response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
		if response.StatusCode != 200 {
			fmt.Println("SC Parameters, ID:", i, response.StatusCode, host)
		}

		err := json.Unmarshal(response.Body, &d)
		if err != nil {
			logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
		} else {
			if d.OverrideValuesCount > 0 {
				fmt.Println("WRK:", wrkID, i, d.OverrideValues)
				fmt.Println("====")
			}
			out <- d
			//addResult(d, r, mtx)
		}
	}
}
