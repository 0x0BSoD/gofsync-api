package utils

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/archops/goFsync/models"
	"github.com/streadway/amqp"
	"runtime"
)

// =====================================================================================================================
// RABBITmQ
// =====================================================================================================================
func InitializeAMQP(cfg *models.Config) {
	connString := fmt.Sprintf("amqp://%s:%s@:%s:%d/",
		cfg.AMQP.Username,
		cfg.AMQP.Password,
		cfg.AMQP.Host,
		cfg.AMQP.Port)
	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	cfg.AMQP.Channel = ch
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"jobs",
		false,
		false,
		false,
		false,
		nil,
	)
	cfg.AMQP.Queue = &q
	failOnError(err, "Failed to declare a queue")
}

func SendToQueue(cfg *models.Config, msg models.Job) error {
	body, err := json.Marshal(msg)
	failOnError(err, "Error on marshal job")
	err = cfg.AMQP.Channel.Publish(
		"",
		cfg.AMQP.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		Error.Printf("%s: %s", msg, err)
	}
}

// =====================================================================================================================
// Goroutine workers
// =====================================================================================================================

// =====================================================================================================================
// WorkQueue is a channel type that you can send Work on.
type WorkQueue chan Work

// New creates a WorkQueue with runtime.NumCPU() workers.
func New() WorkQueue {
	return NewN(runtime.NumCPU())
}

// NewN creates and returns a new WorkQueue that has the specified number
// of workers.
func NewN(numWorkers int) WorkQueue {
	queue := make(WorkQueue)
	d := make(dispatcher, numWorkers)
	go d.dispatch(queue)
	return queue
}

// Work is a task to perform that can be sent over a WorkQueue.
type Work func()

type dispatcher chan chan Work

func (d dispatcher) dispatch(queue WorkQueue) {
	// Create and start all of our workers.
	for i := 0; i < cap(d); i++ {
		w := make(worker)
		go w.work(d)
	}

	// Start the main loop in a goroutine.
	go func() {
		for work := range queue {
			go func(work Work) {
				worker := <-d
				worker <- work
			}(work)
		}

		// If we get here, the work queue has been closed, and we should
		// stop all of the workers.
		for i := 0; i < cap(d); i++ {
			w := <-d
			close(w)
		}
	}()
}

type worker chan Work

func (w worker) work(d dispatcher) {
	// Add ourselves to the dispatcher.
	d <- w

	// Start the main loop.
	go w.wait(d)
}

func (w worker) wait(d dispatcher) {
	for work := range w {
		// Do the work.
		if work == nil {
			panic("nil work received")
		}

		work()

		// Re-add ourselves to the dispatcher.
		d <- w
	}
}
