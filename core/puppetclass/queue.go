package puppetclass

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	logger "git.ringcentral.com/archops/goFsync/utils"
	"log"
	"sync"
)

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
	Results   *[]models.PCSCParameters
	Lock      *sync.Mutex
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
				work(w.ID, job.ForemanID, job.Host, job.Results, job.Lock, job.Cfg)
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

func CreateJobs(foremanIDS []int, host string, res *[]models.PCSCParameters, cfg *models.Config) []Work {
	var jobs []Work
	fmt.Println("OBJ_3:", len(foremanIDS))
	for i, fID := range foremanIDS {
		jobs = append(jobs, Work{
			ID:        i,
			ForemanID: fID,
			Host:      host,
			Results:   res,
			Cfg:       cfg,
		})
	}
	fmt.Println("OBJ_4:", len(jobs))
	return jobs
}

// =====================================================================================================================
func work(wrkID int, i int, host string, summary *[]models.PCSCParameters, lock *sync.Mutex, cfg *models.Config) {

	//fmt.Printf("Worker %d got task: { foremanID:%d }\n", wrkID, i)

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
	lock.Lock()
	*summary = append(*summary, r)
	lock.Unlock()
	//fmt.Printf("Worker %d finish task: { foremanID:%d, data: ", wrkID, i)
	//fmt.Println(r, " }")
}
