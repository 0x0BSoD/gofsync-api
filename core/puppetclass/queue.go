package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
)

var ResultQueue = make(chan models.PCSCParameters, 100)
var WorkQueue = make(chan int, 100)

type Worker struct {
	ID          int
	Work        chan int
	WorkerQueue chan chan int
	ResultQueue chan models.PCSCParameters
	Host        string
	Cfg         *models.Config
	QuitChan    chan bool
}

//func NewWorker(id int, host string, cfg *models.Config) Worker {
//	// Create, and return the worker.
//	worker := Worker{
//		ID:          id,
//		Work:        make(chan int),
//		ResultQueue: ResultQueue,
//		WorkerQueue: make(chan WorkQueue),
//		Host:        host,
//		Cfg:         cfg,
//		QuitChan:    make(chan bool),
//	}
//
//	return worker
//}

func (w *Worker) start() {
	var r models.PCSCParameters
	w.WorkerQueue <- w.Work
	select {
	case i := <-w.Work:
		uri := fmt.Sprintf("puppetclasses/%d", i)
		response, _ := utils.ForemanAPI("GET", w.Host, uri, "", w.Cfg)
		if response.StatusCode != 200 {
			fmt.Println("PuppetClasses updates, ID:", i, response.StatusCode, w.Host)
		}

		err := json.Unmarshal(response.Body, &r)
		if err != nil {
			logger.Error.Printf("%q:\n %q\n", err, response)
		}
		w.ResultQueue <- r

	case <-w.QuitChan:
		// We have been asked to stop.
		fmt.Printf("worker%d stopping\n", w.ID)
		return
	}

}

func worker(i int,
	host string,
	cfg *models.Config) models.PCSCParameters {
	var r models.PCSCParameters
	uri := fmt.Sprintf("puppetclasses/%d", i)
	response, _ := logger.ForemanAPI("GET", host, uri, "", cfg)
	if response.StatusCode != 200 {
		fmt.Println("PuppetClasses updates, ID:", i, response.StatusCode, host)
	}

	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("%q:\n %q\n", err, response)
	}

	return r
}
