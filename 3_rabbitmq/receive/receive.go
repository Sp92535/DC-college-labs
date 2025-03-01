package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672/")
	fail(err, "Error Connecting")
	defer conn.Close()

	ch, err := conn.Channel()
	fail(err, "Error opening channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"YOO",
		false,
		false,
		false,
		false,
		nil,
	)
	fail(err, "Error declaring queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	fail(err, "Erron setting up Consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Recieved Message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func fail(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", err, msg)
	}
}
