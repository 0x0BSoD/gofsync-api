package utils

import (
	"encoding/json"
	"fmt"
	"git.ringcentral.com/alexander.simonov/goFsync/models"
	"github.com/streadway/amqp"
)

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
		"jobs", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
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
