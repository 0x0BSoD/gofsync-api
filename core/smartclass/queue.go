package smartclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"git.ringcentral.com/alexander.simonov/goFsync/utils"
	logger "git.ringcentral.com/alexander.simonov/goFsync/utils"
	"log"
	"sync"
)

var writeLock sync.Mutex
var WorkerChannel = make(chan chan Work)

type Collector struct {
	Work chan Work
	End  chan bool
}

type Work struct {
	ID        int
	ForemanID int
	Host      string
	Cfg       *models.Config
	Results   *[]models.SCParameter
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Work
	Channel       chan Work
	End           chan bool
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				work(w.ID, job.ForemanID, job.Host, job.Results, job.Cfg)
			case <-w.End:
				return
			}
		}
	}()
}
func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}

func StartDispatcher(workerCount int) Collector {
	var i int
	var workers []Worker
	input := make(chan Work)
	end := make(chan bool)
	collector := Collector{Work: input, End: end}

	for i < workerCount {
		i++
		fmt.Println("starting worker: ", i)
		worker := Worker{
			ID:            i,
			Channel:       make(chan Work),
			WorkerChannel: WorkerChannel,
			End:           make(chan bool)}
		worker.Start()
		workers = append(workers, worker)
	}

	// start collector
	go func() {
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop()
				}
				return
			case work := <-input:
				worker := <-WorkerChannel
				worker <- work
			}
		}
	}()

	return collector
}

func CreateJobs(foremanIDS []int, host string, res *[]models.SCParameter, cfg *models.Config) []Work {
	var jobs []Work

	for i, fID := range foremanIDS {
		jobs = append(jobs, Work{
			ID:        i,
			ForemanID: fID,
			Host:      host,
			Results:   res,
			Cfg:       cfg,
		})
	}
	return jobs
}

//======================================================================================================================
func work(wrkID int, i int, host string, summary *[]models.SCParameter, cfg *models.Config) {
	fmt.Printf("Worker %d got task: { foremanID:%d }\n", wrkID, i)

	var r models.SCParameter

	fmt.Printf("W: got task, scId: %d, HOST: %s\n", i, host)

	uri := fmt.Sprintf("smart_class_parameters/%d", i)
	response, _ := utils.ForemanAPI("GET", host, uri, "", cfg)
	if response.StatusCode != 200 {
		fmt.Println("SC Parameters, ID:", i, response.StatusCode, host)
	}

	err := json.Unmarshal(response.Body, &r)
	if err != nil {
		logger.Error.Printf("Error on getting override: %q \n%s\n", err, uri)
	}

	writeLock.Lock()
	*summary = append(*summary, r)
	writeLock.Unlock()
	fmt.Printf("Worker %d finish task: { foremanID:%d, data: ", wrkID, i)
	fmt.Println(r, " }")
}
